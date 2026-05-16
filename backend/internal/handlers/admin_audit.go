package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tofex/backend/internal/httpx"
	"github.com/tofex/backend/internal/models"
)

func (s *Server) AdminListAuditLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	per, _ := strconv.Atoi(c.DefaultQuery("per_page", "50"))
	if page < 1 {
		page = 1
	}
	if per < 1 || per > 200 {
		per = 50
	}
	offset := (page - 1) * per
	q := s.DB.Model(&models.AuditLog{})
	if v := c.Query("user_id"); v != "" {
		if uid, err := strconv.ParseUint(v, 10, 64); err == nil {
			u := uint(uid)
			q = q.Where("user_id = ?", u)
		}
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		httpx.JSONError(c, 500, "server_error", "could not count audit logs")
		return
	}
	var logs []models.AuditLog
	if err := q.Order("id desc").Limit(per).Offset(offset).Find(&logs).Error; err != nil {
		httpx.JSONError(c, 500, "server_error", "could not list audit logs")
		return
	}
	c.JSON(http.StatusOK, gin.H{"audit_logs": logs, "total": total, "page": page, "per_page": per})
}
