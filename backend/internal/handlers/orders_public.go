package handlers

import (
	"encoding/json"
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
	"gorm.io/gorm"
)

type orderLineIn struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

func (s *Server) PublicCreateOrder(c *gin.Context) {
	customerName := strings.TrimSpace(c.PostForm("customer_name"))
	customerEmail := strings.TrimSpace(c.PostForm("customer_email"))
	customerPhone := strings.TrimSpace(c.PostForm("customer_phone"))
	customerNote := strings.TrimSpace(c.PostForm("customer_note"))
	paymentAccountIDStr := c.PostForm("payment_account_id")
	itemsJSON := c.PostForm("items")

	if customerName == "" || customerEmail == "" {
		httpx.JSONError(c, 400, "validation_error", "customer_name and customer_email are required")
		return
	}
	paymentAccountID, err := parseUint(paymentAccountIDStr)
	if err != nil || paymentAccountID == 0 {
		httpx.JSONError(c, 400, "validation_error", "payment_account_id is required")
		return
	}
	if itemsJSON == "" {
		httpx.JSONError(c, 400, "validation_error", "items JSON is required")
		return
	}
	var lines []orderLineIn
	if err := json.Unmarshal([]byte(itemsJSON), &lines); err != nil {
		httpx.JSONError(c, 400, "validation_error", "items must be a JSON array of {product_id, quantity}")
		return
	}
	if len(lines) == 0 {
		httpx.JSONError(c, 400, "validation_error", "at least one line item is required")
		return
	}
	fh, err := c.FormFile("receipt")
	if err != nil {
		httpx.JSONError(c, 400, "validation_error", "receipt file is required")
		return
	}
	if fh.Size > int64(s.Cfg.MaxUploadMB)*1024*1024 {
		httpx.JSONError(c, 400, "payload_too_large", "receipt exceeds upload limit")
		return
	}

	var acct models.PaymentAccount
	if err := s.DB.Where("id = ? AND active = ?", paymentAccountID, true).First(&acct).Error; err != nil {
		httpx.JSONError(c, 400, "validation_error", "invalid or inactive payment account")
		return
	}
	snap, _ := json.Marshal(acct)

	type lineResolved struct {
		Product models.Product
		Qty     int
	}
	resolved := make([]lineResolved, 0, len(lines))
	var total int64
	for _, ln := range lines {
		if ln.ProductID == 0 || ln.Quantity <= 0 {
			httpx.JSONError(c, 400, "validation_error", "each item needs product_id and positive quantity")
			return
		}
		var p models.Product
		if err := s.DB.Where("id = ? AND active = ?", ln.ProductID, true).First(&p).Error; err != nil {
			httpx.JSONError(c, 400, "validation_error", "unknown or inactive product id")
			return
		}
		total += p.PriceCents * int64(ln.Quantity)
		resolved = append(resolved, lineResolved{Product: p, Qty: ln.Quantity})
	}

	orderID := uuid.New()
	receiptExt := strings.ToLower(filepath.Ext(fh.Filename))
	if receiptExt == "" {
		receiptExt = ".bin"
	}
	storedReceiptName := orderID.String() + receiptExt
	relReceipt := filepath.ToSlash(filepath.Join("receipts", orderID.String(), storedReceiptName))

	opened, err := fh.Open()
	if err != nil {
		httpx.JSONError(c, 400, "validation_error", "could not read receipt")
		return
	}
	defer opened.Close()

	head := make([]byte, 512)
	n, err := opened.Read(head)
	if err != nil && !errors.Is(err, io.EOF) {
		httpx.JSONError(c, 400, "validation_error", "could not read receipt")
		return
	}
	if n == 0 {
		httpx.JSONError(c, 400, "validation_error", "receipt file is empty")
		return
	}
	mime := http.DetectContentType(head[:n])
	if !strings.HasPrefix(mime, "image/") {
		httpx.JSONError(c, 400, "validation_error", "receipt must be an image (jpeg, png, webp, etc.)")
		return
	}
	if _, err := opened.Seek(0, 0); err != nil {
		httpx.JSONError(c, 500, "server_error", "could not rewind upload")
		return
	}

	destAbs, err := fsx.SafeJoinUploadDir(s.Cfg.UploadDir, relReceipt)
	if err != nil {
		httpx.JSONError(c, 500, "server_error", "invalid upload path")
		return
	}
	if err := os.MkdirAll(filepath.Dir(destAbs), 0o755); err != nil {
		httpx.JSONError(c, 500, "server_error", "could not create upload directory")
		return
	}
	out, err := os.Create(destAbs)
	if err != nil {
		httpx.JSONError(c, 500, "server_error", "could not store receipt")
		return
	}
	if _, err := io.Copy(out, opened); err != nil {
		out.Close()
		httpx.JSONError(c, 500, "server_error", "could not store receipt")
		return
	}
	_ = out.Close()

	o := models.Order{
		ID:                  orderID,
		CustomerName:        customerName,
		CustomerEmail:       customerEmail,
		CustomerPhone:       customerPhone,
		PaymentAccountID:    acct.ID,
		PaymentSnapshotJSON: string(snap),
		TotalCents:          total,
		Status:              "pending",
		CustomerNote:        customerNote,
	}
	if err := s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&o).Error; err != nil {
			return err
		}
		for _, r := range resolved {
			it := models.OrderItem{
				OrderID:         orderID,
				ProductID:       r.Product.ID,
				ProductNameSnap: r.Product.Name,
				UnitPriceCents:  r.Product.PriceCents,
				Quantity:        r.Qty,
			}
			if err := tx.Create(&it).Error; err != nil {
				return err
			}
		}
		rec := models.OrderReceipt{
			OrderID:   orderID,
			FilePath:  relReceipt,
			MIME:      mime,
			SizeBytes: fh.Size,
		}
		return tx.Create(&rec).Error
	}); err != nil {
		_ = os.Remove(destAbs)
		httpx.JSONError(c, 500, "server_error", "could not create order")
		return
	}

	u := middleware.CurrentUser(c)
	var uid *uint
	if u != nil {
		uid = &u.ID
	}
	audit.Log(s.DB, uid, "order.created_public", "order", orderID.String(), map[string]any{"total_cents": total}, httpx.ClientIP(c), c.Request.UserAgent())

	c.JSON(http.StatusCreated, gin.H{"order_id": orderID.String(), "total_cents": total, "status": o.Status})
}

func parseUint(s string) (uint, error) {
	n, err := strconv.ParseUint(strings.TrimSpace(s), 10, 64)
	return uint(n), err
}
