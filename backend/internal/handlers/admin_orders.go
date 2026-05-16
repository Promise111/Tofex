package handlers

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tofex/backend/internal/audit"
	"github.com/tofex/backend/internal/fsx"
	"github.com/tofex/backend/internal/httpx"
	"github.com/tofex/backend/internal/middleware"
	"github.com/tofex/backend/internal/models"
)

var allowedOrderStatuses = map[string]struct{}{
	"pending":   {},
	"confirmed": {},
	"ready":     {},
	"completed": {},
	"cancelled": {},
}

func (s *Server) AdminListOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	per, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	status := strings.TrimSpace(c.Query("status"))
	if page < 1 {
		page = 1
	}
	if per < 1 || per > 100 {
		per = 20
	}
	offset := (page - 1) * per
	q := s.DB.Model(&models.Order{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		httpx.JSONError(c, 500, "server_error", "could not count orders")
		return
	}
	var orders []models.Order
	db := s.DB.Preload("Items").Preload("Receipts").Order("created_at desc").Limit(per).Offset(offset)
	if status != "" {
		db = db.Where("status = ?", status)
	}
	if err := db.Find(&orders).Error; err != nil {
		httpx.JSONError(c, 500, "server_error", "could not list orders")
		return
	}
	c.JSON(http.StatusOK, gin.H{"orders": orders, "total": total, "page": page, "per_page": per})
}

func (s *Server) AdminGetOrder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.JSONError(c, 400, "validation_error", "invalid order id")
		return
	}
	var o models.Order
	if err := s.DB.Preload("Items").Preload("Receipts").First(&o, "id = ?", id).Error; err != nil {
		httpx.JSONError(c, 404, "not_found", "order not found")
		return
	}
	c.JSON(http.StatusOK, o)
}

type patchOrderBody struct {
	Status string `json:"status" binding:"required"`
}

func (s *Server) AdminPatchOrder(c *gin.Context) {
	cur := middleware.CurrentUser(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.JSONError(c, 400, "validation_error", "invalid order id")
		return
	}
	var body patchOrderBody
	if err := c.ShouldBindJSON(&body); err != nil {
		httpx.JSONError(c, 400, "validation_error", err.Error())
		return
	}
	if _, ok := allowedOrderStatuses[strings.TrimSpace(body.Status)]; !ok {
		httpx.JSONError(c, 400, "validation_error", "invalid status")
		return
	}
	var o models.Order
	if err := s.DB.First(&o, "id = ?", id).Error; err != nil {
		httpx.JSONError(c, 404, "not_found", "order not found")
		return
	}
	prev := o.Status
	o.Status = strings.TrimSpace(body.Status)
	if err := s.DB.Save(&o).Error; err != nil {
		httpx.JSONError(c, 500, "server_error", "could not update order")
		return
	}
	audit.Log(s.DB, &cur.ID, "order.update_status", "order", id.String(), map[string]any{"from": prev, "to": o.Status}, httpx.ClientIP(c), c.Request.UserAgent())
	c.JSON(http.StatusOK, o)
}

func (s *Server) AdminDownloadOrderReceipt(c *gin.Context) {
	oid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httpx.JSONError(c, 400, "validation_error", "invalid order id")
		return
	}
	rid, err := strconv.ParseUint(c.Param("rid"), 10, 64)
	if err != nil {
		httpx.JSONError(c, 400, "validation_error", "invalid receipt id")
		return
	}
	var rec models.OrderReceipt
	if err := s.DB.Where("id = ? AND order_id = ?", rid, oid).First(&rec).Error; err != nil {
		httpx.JSONError(c, 404, "not_found", "receipt not found")
		return
	}
	full, err := fsx.SafeJoinUploadDir(s.Cfg.UploadDir, rec.FilePath)
	if err != nil {
		httpx.JSONError(c, 400, "bad_path", "invalid path")
		return
	}
	if _, err := os.Stat(full); err != nil {
		httpx.JSONError(c, 404, "not_found", "file missing")
		return
	}
	c.File(full)
}
