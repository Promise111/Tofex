package handlers

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tofex/backend/internal/audit"
	"github.com/tofex/backend/internal/fsx"
	"github.com/tofex/backend/internal/httpx"
	"github.com/tofex/backend/internal/middleware"
	"github.com/tofex/backend/internal/models"
	"github.com/tofex/backend/internal/slug"
)

func (s *Server) ensureUniqueProductSlug(base string, ignoreID uint) (string, error) {
	slugStr := slug.FromName(base)
	candidate := slugStr
	for i := 0; i < 50; i++ {
		var count int64
		q := s.DB.Model(&models.Product{}).Where("slug = ?", candidate)
		if ignoreID > 0 {
			q = q.Where("id <> ?", ignoreID)
		}
		if err := q.Count(&count).Error; err != nil {
			return "", err
		}
		if count == 0 {
			return candidate, nil
		}
		candidate = slugStr + "-" + strconv.Itoa(i+2)
	}
	return "", errors.New("could not allocate slug")
}

func (s *Server) AdminListProducts(c *gin.Context) {
	var list []models.Product
	if err := s.DB.Preload("Images").Order("id desc").Find(&list).Error; err != nil {
		httpx.JSONError(c, 500, "server_error", "could not list products")
		return
	}
	c.JSON(http.StatusOK, gin.H{"products": list})
}

type createProductBody struct {
	Name        string `json:"name" binding:"required,max=255"`
	Slug        string `json:"slug" binding:"max=255"`
	Description string `json:"description"`
	PriceCents  int64  `json:"price_cents" binding:"required,min=1"`
	Active      *bool  `json:"active"`
}

func (s *Server) AdminCreateProduct(c *gin.Context) {
	cur := middleware.CurrentUser(c)
	var body createProductBody
	if err := c.ShouldBindJSON(&body); err != nil {
		httpx.JSONError(c, 400, "validation_error", err.Error())
		return
	}
	slugStr := strings.TrimSpace(body.Slug)
	if slugStr == "" {
		var err error
		slugStr, err = s.ensureUniqueProductSlug(body.Name, 0)
		if err != nil {
			httpx.JSONError(c, 500, "server_error", "could not allocate slug")
			return
		}
	} else {
		var count int64
		s.DB.Model(&models.Product{}).Where("slug = ?", slugStr).Count(&count)
		if count > 0 {
			httpx.JSONError(c, 400, "validation_error", "slug already in use")
			return
		}
	}
	active := true
	if body.Active != nil {
		active = *body.Active
	}
	p := models.Product{
		Name:        body.Name,
		Slug:        slugStr,
		Description: body.Description,
		PriceCents:  body.PriceCents,
		Active:      active,
	}
	if err := s.DB.Create(&p).Error; err != nil {
		httpx.JSONError(c, 400, "create_failed", "could not create product")
		return
	}
	audit.Log(s.DB, &cur.ID, "product.create", "product", itoa(p.ID), map[string]any{"slug": p.Slug}, httpx.ClientIP(c), c.Request.UserAgent())
	c.JSON(http.StatusCreated, p)
}

type patchProductBody struct {
	Name        *string `json:"name"`
	Slug        *string `json:"slug"`
	Description *string `json:"description"`
	PriceCents  *int64  `json:"price_cents"`
	Active      *bool   `json:"active"`
}

func (s *Server) AdminPatchProduct(c *gin.Context) {
	cur := middleware.CurrentUser(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httpx.JSONError(c, 400, "validation_error", "invalid id")
		return
	}
	var p models.Product
	if err := s.DB.First(&p, id).Error; err != nil {
		httpx.JSONError(c, 404, "not_found", "product not found")
		return
	}
	var body patchProductBody
	if err := c.ShouldBindJSON(&body); err != nil {
		httpx.JSONError(c, 400, "validation_error", err.Error())
		return
	}
	if body.Name != nil {
		p.Name = *body.Name
	}
	if body.Description != nil {
		p.Description = *body.Description
	}
	if body.PriceCents != nil {
		if *body.PriceCents < 1 {
			httpx.JSONError(c, 400, "validation_error", "price_cents must be positive")
			return
		}
		p.PriceCents = *body.PriceCents
	}
	if body.Active != nil {
		p.Active = *body.Active
	}
	if body.Slug != nil {
		slugStr := strings.TrimSpace(*body.Slug)
		if slugStr == "" {
			newSlug, err := s.ensureUniqueProductSlug(p.Name, p.ID)
			if err != nil {
				httpx.JSONError(c, 500, "server_error", "could not allocate slug")
				return
			}
			p.Slug = newSlug
		} else {
			var count int64
			s.DB.Model(&models.Product{}).Where("slug = ? AND id <> ?", slugStr, p.ID).Count(&count)
			if count > 0 {
				httpx.JSONError(c, 400, "validation_error", "slug already in use")
				return
			}
			p.Slug = slugStr
		}
	}
	if err := s.DB.Save(&p).Error; err != nil {
		httpx.JSONError(c, 500, "server_error", "could not update product")
		return
	}
	audit.Log(s.DB, &cur.ID, "product.update", "product", itoa(p.ID), nil, httpx.ClientIP(c), c.Request.UserAgent())
	c.JSON(http.StatusOK, p)
}

func (s *Server) AdminDeleteProduct(c *gin.Context) {
	cur := middleware.CurrentUser(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httpx.JSONError(c, 400, "validation_error", "invalid id")
		return
	}
	res := s.DB.Delete(&models.Product{}, id)
	if res.Error != nil {
		httpx.JSONError(c, 500, "server_error", "delete failed")
		return
	}
	if res.RowsAffected == 0 {
		httpx.JSONError(c, 404, "not_found", "product not found")
		return
	}
	audit.Log(s.DB, &cur.ID, "product.delete", "product", strconv.FormatUint(id, 10), nil, httpx.ClientIP(c), c.Request.UserAgent())
	c.Status(http.StatusNoContent)
}

func (s *Server) AdminUploadProductImage(c *gin.Context) {
	cur := middleware.CurrentUser(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httpx.JSONError(c, 400, "validation_error", "invalid id")
		return
	}
	var p models.Product
	if err := s.DB.First(&p, id).Error; err != nil {
		httpx.JSONError(c, 404, "not_found", "product not found")
		return
	}
	fh, err := c.FormFile("file")
	if err != nil {
		httpx.JSONError(c, 400, "validation_error", "file is required")
		return
	}
	if fh.Size > int64(s.Cfg.MaxUploadMB)*1024*1024 {
		httpx.JSONError(c, 400, "payload_too_large", "file exceeds upload limit")
		return
	}
	ext := strings.ToLower(filepath.Ext(fh.Filename))
	if ext == "" {
		ext = ".bin"
	}
	stored := uuid.New().String() + ext
	rel := filepath.ToSlash(filepath.Join("products", strconv.FormatUint(id, 10), stored))
	destAbs, err := fsx.SafeJoinUploadDir(s.Cfg.UploadDir, rel)
	if err != nil {
		httpx.JSONError(c, 500, "server_error", "invalid path")
		return
	}
	if err := os.MkdirAll(filepath.Dir(destAbs), 0o755); err != nil {
		httpx.JSONError(c, 500, "server_error", "could not create directory")
		return
	}
	src, err := fh.Open()
	if err != nil {
		httpx.JSONError(c, 400, "validation_error", "could not read file")
		return
	}
	defer src.Close()
	dst, err := os.Create(destAbs)
	if err != nil {
		httpx.JSONError(c, 500, "server_error", "could not store file")
		return
	}
	if _, err := io.Copy(dst, src); err != nil {
		dst.Close()
		httpx.JSONError(c, 500, "server_error", "could not store file")
		return
	}
	_ = dst.Close()

	img := models.ProductImage{
		ProductID: uint(id),
		Path:      rel,
		SortOrder: 0,
	}
	if err := s.DB.Create(&img).Error; err != nil {
		_ = os.Remove(destAbs)
		httpx.JSONError(c, 500, "server_error", "could not save image metadata")
		return
	}
	audit.Log(s.DB, &cur.ID, "product.image_upload", "product", strconv.FormatUint(id, 10), map[string]any{"image_id": img.ID}, httpx.ClientIP(c), c.Request.UserAgent())
	c.JSON(http.StatusCreated, img)
}
