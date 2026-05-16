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

func (s *Server) AdminListPaymentAccounts(c *gin.Context) {
	var list []models.PaymentAccount
	if err := s.DB.Order("id asc").Find(&list).Error; err != nil {
		httpx.JSONError(c, 500, "server_error", "could not list payment accounts")
		return
	}
	c.JSON(http.StatusOK, gin.H{"payment_accounts": list})
}

type createPaymentAccountBody struct {
	BankName       string `json:"bank_name" binding:"required,max=255"`
	AccountName    string `json:"account_name" binding:"required,max=255"`
	AccountNumber  string `json:"account_number" binding:"required,max=64"`
	Branch         string `json:"branch" binding:"max=255"`
	Currency       string `json:"currency" binding:"max=8"`
	DisplayLabel   string `json:"display_label" binding:"max=255"`
	Active         *bool  `json:"active"`
}

func (s *Server) AdminCreatePaymentAccount(c *gin.Context) {
	cur := middleware.CurrentUser(c)
	var body createPaymentAccountBody
	if err := c.ShouldBindJSON(&body); err != nil {
		httpx.JSONError(c, 400, "validation_error", err.Error())
		return
	}
	curStr := "NGN"
	if body.Currency != "" {
		curStr = body.Currency
	}
	active := true
	if body.Active != nil {
		active = *body.Active
	}
	a := models.PaymentAccount{
		BankName:      body.BankName,
		AccountName:   body.AccountName,
		AccountNumber: body.AccountNumber,
		Branch:        body.Branch,
		Currency:      curStr,
		DisplayLabel:  body.DisplayLabel,
		Active:        active,
	}
	if err := s.DB.Create(&a).Error; err != nil {
		httpx.JSONError(c, 400, "create_failed", "could not create payment account")
		return
	}
	audit.Log(s.DB, &cur.ID, "payment_account.create", "payment_account", itoa(a.ID), nil, httpx.ClientIP(c), c.Request.UserAgent())
	c.JSON(http.StatusCreated, a)
}

type patchPaymentAccountBody struct {
	BankName       *string `json:"bank_name"`
	AccountName    *string `json:"account_name"`
	AccountNumber  *string `json:"account_number"`
	Branch         *string `json:"branch"`
	Currency       *string `json:"currency"`
	DisplayLabel   *string `json:"display_label"`
	Active         *bool   `json:"active"`
}

func (s *Server) AdminPatchPaymentAccount(c *gin.Context) {
	cur := middleware.CurrentUser(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httpx.JSONError(c, 400, "validation_error", "invalid id")
		return
	}
	var a models.PaymentAccount
	if err := s.DB.First(&a, id).Error; err != nil {
		httpx.JSONError(c, 404, "not_found", "payment account not found")
		return
	}
	var body patchPaymentAccountBody
	if err := c.ShouldBindJSON(&body); err != nil {
		httpx.JSONError(c, 400, "validation_error", err.Error())
		return
	}
	if body.BankName != nil {
		a.BankName = *body.BankName
	}
	if body.AccountName != nil {
		a.AccountName = *body.AccountName
	}
	if body.AccountNumber != nil {
		a.AccountNumber = *body.AccountNumber
	}
	if body.Branch != nil {
		a.Branch = *body.Branch
	}
	if body.Currency != nil {
		a.Currency = *body.Currency
	}
	if body.DisplayLabel != nil {
		a.DisplayLabel = *body.DisplayLabel
	}
	if body.Active != nil {
		a.Active = *body.Active
	}
	if err := s.DB.Save(&a).Error; err != nil {
		httpx.JSONError(c, 500, "server_error", "could not update payment account")
		return
	}
	audit.Log(s.DB, &cur.ID, "payment_account.update", "payment_account", itoa(a.ID), nil, httpx.ClientIP(c), c.Request.UserAgent())
	c.JSON(http.StatusOK, a)
}

func (s *Server) AdminDeletePaymentAccount(c *gin.Context) {
	cur := middleware.CurrentUser(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		httpx.JSONError(c, 400, "validation_error", "invalid id")
		return
	}
	res := s.DB.Delete(&models.PaymentAccount{}, id)
	if res.Error != nil {
		httpx.JSONError(c, 500, "server_error", "delete failed")
		return
	}
	if res.RowsAffected == 0 {
		httpx.JSONError(c, 404, "not_found", "payment account not found")
		return
	}
	audit.Log(s.DB, &cur.ID, "payment_account.delete", "payment_account", strconv.FormatUint(id, 10), nil, httpx.ClientIP(c), c.Request.UserAgent())
	c.Status(http.StatusNoContent)
}
