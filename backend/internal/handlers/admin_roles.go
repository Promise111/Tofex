package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tofex/backend/internal/audit"
	"github.com/tofex/backend/internal/httpx"
	"github.com/tofex/backend/internal/middleware"
	"github.com/tofex/backend/internal/models"
	"gorm.io/gorm"
)

func (s *Server) AdminListRoles(c *gin.Context) {
	var roles []models.Role
	if err := s.DB.Preload("Permissions").Order("id asc").Find(&roles).Error; err != nil {
		httpx.JSONError(c, 500, "server_error", "could not list roles")
		return
	}
	c.JSON(http.StatusOK, gin.H{"roles": roles})
}

type patchRoleBody struct {
	Permissions []string `json:"permissions" binding:"required"`
}

func (s *Server) AdminPatchRole(c *gin.Context) {
	cur := middleware.CurrentUser(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httpx.JSONError(c, 400, "validation_error", "invalid id")
		return
	}
	var role models.Role
	if err := s.DB.First(&role, id).Error; err != nil {
		httpx.JSONError(c, 404, "not_found", "role not found")
		return
	}
	var body patchRoleBody
	if err := c.ShouldBindJSON(&body); err != nil {
		httpx.JSONError(c, 400, "validation_error", err.Error())
		return
	}
	perms := body.Permissions
	if role.Name == "super_admin" {
		perms = []string{"*"}
	}
	if err := s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", role.ID).Delete(&models.RolePermission{}).Error; err != nil {
			return err
		}
		for _, p := range perms {
			if p == "" {
				continue
			}
			rp := models.RolePermission{RoleID: role.ID, Permission: p}
			if err := tx.Create(&rp).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		httpx.JSONError(c, 500, "server_error", "could not update role permissions")
		return
	}
	audit.Log(s.DB, &cur.ID, "role.update_permissions", "role", itoa(role.ID), map[string]any{"permissions": perms}, httpx.ClientIP(c), c.Request.UserAgent())
	httpx.JSONOK(c, gin.H{"message": "updated"})
}
