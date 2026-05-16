package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tofex/backend/internal/audit"
	"github.com/tofex/backend/internal/httpx"
	"github.com/tofex/backend/internal/middleware"
	"github.com/tofex/backend/internal/models"
)

// AdminListStoreBranches lists all pickup locations.
func (s *Server) AdminListStoreBranches(c *gin.Context) {
	var list []models.StoreBranch
	if err := s.DB.Order("sort_order asc, id asc").Find(&list).Error; err != nil {
		httpx.JSONError(c, 500, "server_error", "could not list branches")
		return
	}
	c.JSON(http.StatusOK, gin.H{"branches": list})
}

type createStoreBranchBody struct {
	Name      string   `json:"name" binding:"required,max=255"`
	Address   string   `json:"address" binding:"required"`
	City      string   `json:"city" binding:"max=128"`
	Phone     string   `json:"phone" binding:"max=64"`
	Hours     string   `json:"hours" binding:"max=255"`
	MapsURL   string   `json:"maps_url" binding:"max=1024"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
	Active    *bool    `json:"active"`
	SortOrder *int     `json:"sort_order"`
}

// AdminCreateStoreBranch adds a pickup location.
func (s *Server) AdminCreateStoreBranch(c *gin.Context) {
	cur := middleware.CurrentUser(c)
	var body createStoreBranchBody
	if err := c.ShouldBindJSON(&body); err != nil {
		httpx.JSONError(c, 400, "validation_error", err.Error())
		return
	}
	active := true
	if body.Active != nil {
		active = *body.Active
	}
	sortOrder := 0
	if body.SortOrder != nil {
		sortOrder = *body.SortOrder
	}
	b := models.StoreBranch{
		Name:      body.Name,
		Address:   body.Address,
		City:      body.City,
		Phone:     body.Phone,
		Hours:     body.Hours,
		MapsURL:   body.MapsURL,
		Latitude:  body.Latitude,
		Longitude: body.Longitude,
		Active:    active,
		SortOrder: sortOrder,
	}
	if err := s.DB.Create(&b).Error; err != nil {
		httpx.JSONError(c, 400, "create_failed", "could not create branch")
		return
	}
	audit.Log(s.DB, &cur.ID, "store_branch.create", "store_branch", itoa(b.ID), nil, httpx.ClientIP(c), c.Request.UserAgent())
	c.JSON(http.StatusCreated, b)
}

type patchStoreBranchBody struct {
	Name      *string  `json:"name"`
	Address   *string  `json:"address"`
	City      *string  `json:"city"`
	Phone     *string  `json:"phone"`
	Hours     *string  `json:"hours"`
	MapsURL   *string  `json:"maps_url"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
	Active    *bool    `json:"active"`
	SortOrder *int     `json:"sort_order"`
}

// AdminPatchStoreBranch updates a pickup location.
func (s *Server) AdminPatchStoreBranch(c *gin.Context) {
	cur := middleware.CurrentUser(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httpx.JSONError(c, 400, "validation_error", "invalid id")
		return
	}
	var b models.StoreBranch
	if err := s.DB.First(&b, id).Error; err != nil {
		httpx.JSONError(c, 404, "not_found", "branch not found")
		return
	}
	var body patchStoreBranchBody
	if err := c.ShouldBindJSON(&body); err != nil {
		httpx.JSONError(c, 400, "validation_error", err.Error())
		return
	}
	if body.Name != nil {
		b.Name = *body.Name
	}
	if body.Address != nil {
		b.Address = *body.Address
	}
	if body.City != nil {
		b.City = *body.City
	}
	if body.Phone != nil {
		b.Phone = *body.Phone
	}
	if body.Hours != nil {
		b.Hours = *body.Hours
	}
	if body.MapsURL != nil {
		b.MapsURL = *body.MapsURL
	}
	if body.Latitude != nil {
		b.Latitude = body.Latitude
	}
	if body.Longitude != nil {
		b.Longitude = body.Longitude
	}
	if body.Active != nil {
		b.Active = *body.Active
	}
	if body.SortOrder != nil {
		b.SortOrder = *body.SortOrder
	}
	if err := s.DB.Save(&b).Error; err != nil {
		httpx.JSONError(c, 500, "server_error", "could not update branch")
		return
	}
	audit.Log(s.DB, &cur.ID, "store_branch.update", "store_branch", itoa(b.ID), nil, httpx.ClientIP(c), c.Request.UserAgent())
	c.JSON(http.StatusOK, b)
}

// AdminDeleteStoreBranch soft-deletes a pickup location.
func (s *Server) AdminDeleteStoreBranch(c *gin.Context) {
	cur := middleware.CurrentUser(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httpx.JSONError(c, 400, "validation_error", "invalid id")
		return
	}
	res := s.DB.Delete(&models.StoreBranch{}, id)
	if res.Error != nil {
		httpx.JSONError(c, 500, "server_error", "delete failed")
		return
	}
	if res.RowsAffected == 0 {
		httpx.JSONError(c, 404, "not_found", "branch not found")
		return
	}
	audit.Log(s.DB, &cur.ID, "store_branch.delete", "store_branch", strconv.FormatUint(id, 10), nil, httpx.ClientIP(c), c.Request.UserAgent())
	c.Status(http.StatusNoContent)
}
