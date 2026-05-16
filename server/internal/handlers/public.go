package handlers

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tofex/backend/internal/fsx"
	"github.com/tofex/backend/internal/httpx"
	"github.com/tofex/backend/internal/models"
)

// PublicListProducts lists active products for the storefront.
//
//	@Summary		List active products (public)
//	@Tags			public
//	@Produce		json
//	@Success		200	{object}	apidocs.PublicProductsListResponse
//	@Failure		500	{object}	apidocs.ErrorResponse
//	@Router			/api/v1/public/products [get]
func (s *Server) PublicListProducts(c *gin.Context) {
	var list []models.Product
	if err := s.DB.Where("active = ?", true).Preload("Images").Order("id asc").Find(&list).Error; err != nil {
		httpx.JSONError(c, 500, "server_error", "could not list products")
		return
	}
	c.JSON(http.StatusOK, gin.H{"products": list})
}

// PublicGetProduct returns one active product by numeric id or slug.
//
//	@Summary		Get product (public)
//	@Tags			public
//	@Produce		json
//	@Param			id	path		string	true	"Product ID or slug"
//	@Success		200	{object}	models.Product
//	@Failure		404	{object}	apidocs.ErrorResponse
//	@Router			/api/v1/public/products/{id} [get]
func (s *Server) PublicGetProduct(c *gin.Context) {
	idOrSlug := c.Param("id")
	var p models.Product
	q := s.DB.Preload("Images").Where("active = ?", true)
	if _, err := strconv.ParseUint(idOrSlug, 10, 64); err == nil {
		q = q.Where("id = ?", idOrSlug)
	} else {
		q = q.Where("slug = ?", idOrSlug)
	}
	if err := q.First(&p).Error; err != nil {
		httpx.JSONError(c, 404, "not_found", "product not found")
		return
	}
	c.JSON(http.StatusOK, p)
}

// PublicListStoreBranches lists active pickup locations for customers.
func (s *Server) PublicListStoreBranches(c *gin.Context) {
	var list []models.StoreBranch
	if err := s.DB.Where("active = ?", true).Order("sort_order asc, id asc").Find(&list).Error; err != nil {
		httpx.JSONError(c, 500, "server_error", "could not list branches")
		return
	}
	c.JSON(http.StatusOK, gin.H{"branches": list})
}

// PublicListPaymentAccounts lists active bank accounts customers can pay into.
//
//	@Summary		List payment accounts (public)
//	@Tags			public
//	@Produce		json
//	@Success		200	{object}	apidocs.PaymentAccountsListResponse
//	@Failure		500	{object}	apidocs.ErrorResponse
//	@Router			/api/v1/public/payment-accounts [get]
func (s *Server) PublicListPaymentAccounts(c *gin.Context) {
	var list []models.PaymentAccount
	if err := s.DB.Where("active = ?", true).Order("id asc").Find(&list).Error; err != nil {
		httpx.JSONError(c, 500, "server_error", "could not list payment accounts")
		return
	}
	c.JSON(http.StatusOK, gin.H{"payment_accounts": list})
}

// PublicProductImageFile serves a product image binary.
//
//	@Summary		Download product image (public)
//	@Tags			public
//	@Produce		application/octet-stream
//	@Param			id	path	int	true	"Product image ID"
//	@Success		200	{file}	binary
//	@Failure		400	{object}	apidocs.ErrorResponse
//	@Failure		404	{object}	apidocs.ErrorResponse
//	@Router			/api/v1/public/product-images/{id} [get]
func (s *Server) PublicProductImageFile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httpx.JSONError(c, 400, "validation_error", "invalid image id")
		return
	}
	var img models.ProductImage
	if err := s.DB.First(&img, id).Error; err != nil {
		httpx.JSONError(c, 404, "not_found", "image not found")
		return
	}
	var p models.Product
	if err := s.DB.First(&p, img.ProductID).Error; err != nil || !p.Active {
		httpx.JSONError(c, 404, "not_found", "image not found")
		return
	}
	if img.Path == "" {
		httpx.JSONError(c, 404, "not_found", "no file for image")
		return
	}
	full, err := fsx.SafeJoinUploadDir(s.Cfg.UploadDir, img.Path)
	if err != nil {
		httpx.JSONError(c, 400, "bad_path", "invalid stored path")
		return
	}
	if _, err := os.Stat(full); err != nil {
		httpx.JSONError(c, 404, "not_found", "file missing")
		return
	}
	c.File(full)
}
