package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tofex/backend/internal/audit"
	"github.com/tofex/backend/internal/config"
	"github.com/tofex/backend/internal/database"
	"github.com/tofex/backend/internal/httpx"
	"github.com/tofex/backend/internal/models"
	"github.com/tofex/backend/internal/security/password"
	jwtutil "github.com/tofex/backend/internal/security/jwt"
	"gorm.io/gorm"
)

type Server struct {
	DB  *gorm.DB
	Cfg *config.Config
}

func NewServer(db *gorm.DB, cfg *config.Config) *Server {
	return &Server{DB: db, Cfg: cfg}
}

type bootstrapBody struct {
	Email       string `json:"email" binding:"required,email"`
	Username    string `json:"username" binding:"required,min=2,max=64"`
	Password    string `json:"password" binding:"required,min=10,max=128"`
	DisplayName string `json:"display_name" binding:"max=255"`
}

func (s *Server) Bootstrap(c *gin.Context) {
	n, err := database.CountUsers(s.DB)
	if err != nil {
		httpx.JSONError(c, 500, "server_error", "could not verify bootstrap state")
		return
	}
	if n > 0 {
		httpx.JSONError(c, 403, "bootstrap_disabled", "bootstrap is only allowed before the first user exists")
		return
	}
	var body bootstrapBody
	if err := c.ShouldBindJSON(&body); err != nil {
		httpx.JSONError(c, 400, "validation_error", err.Error())
		return
	}
	hash, err := password.Hash(body.Password)
	if err != nil {
		httpx.JSONError(c, 500, "server_error", "password hashing failed")
		return
	}
	var role models.Role
	if err := s.DB.Where("name = ?", "super_admin").First(&role).Error; err != nil {
		httpx.JSONError(c, 500, "server_error", "super_admin role missing; restart server to seed RBAC")
		return
	}
	u := models.User{
		Email:        body.Email,
		Username:     body.Username,
		PasswordHash: hash,
		DisplayName:  body.DisplayName,
	}
	if err := s.DB.Create(&u).Error; err != nil {
		httpx.JSONError(c, 400, "create_failed", "could not create user (email/username must be unique)")
		return
	}
	if err := s.DB.Model(&u).Association("Roles").Append(&role); err != nil {
		httpx.JSONError(c, 500, "server_error", "could not assign super_admin role")
		return
	}
	uid := u.ID
	audit.Log(s.DB, &uid, "bootstrap.first_user", "user", itoa(u.ID), map[string]any{"email": u.Email}, httpx.ClientIP(c), c.Request.UserAgent())
	httpx.JSONOK(c, gin.H{"message": "first admin created", "user_id": u.ID})
}

type loginBody struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (s *Server) Login(c *gin.Context) {
	var body loginBody
	if err := c.ShouldBindJSON(&body); err != nil {
		httpx.JSONError(c, 400, "validation_error", err.Error())
		return
	}
	var u models.User
	if err := s.DB.Preload("Roles.Permissions").Where("email = ? OR username = ?", body.Login, body.Login).First(&u).Error; err != nil {
		audit.Log(s.DB, nil, "auth.login_failed", "auth", "", map[string]any{"login": body.Login}, httpx.ClientIP(c), c.Request.UserAgent())
		httpx.JSONError(c, 401, "invalid_credentials", "invalid email/username or password")
		return
	}
	if !password.Verify(u.PasswordHash, body.Password) {
		audit.Log(s.DB, &u.ID, "auth.login_failed", "auth", itoa(u.ID), nil, httpx.ClientIP(c), c.Request.UserAgent())
		httpx.JSONError(c, 401, "invalid_credentials", "invalid email/username or password")
		return
	}
	tok, err := jwtutil.Sign(u.ID, s.Cfg.JWTSecret, s.Cfg.JWTExpiry)
	if err != nil {
		httpx.JSONError(c, 500, "server_error", "could not issue token")
		return
	}
	audit.Log(s.DB, &u.ID, "auth.login", "auth", itoa(u.ID), nil, httpx.ClientIP(c), c.Request.UserAgent())
	c.JSON(http.StatusOK, gin.H{
		"access_token": tok,
		"token_type":   "Bearer",
		"expires_in":   int(s.Cfg.JWTExpiry.Seconds()),
		"user":         u,
	})
}

type forgotBody struct {
	Email string `json:"email" binding:"required,email"`
}

func (s *Server) ForgotPassword(c *gin.Context) {
	var body forgotBody
	if err := c.ShouldBindJSON(&body); err != nil {
		httpx.JSONError(c, 400, "validation_error", err.Error())
		return
	}
	var u models.User
	if err := s.DB.Where("email = ?", body.Email).First(&u).Error; err != nil {
		// uniform response
		c.JSON(http.StatusOK, gin.H{"message": "if the account exists, reset instructions will follow"})
		return
	}
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		httpx.JSONError(c, 500, "server_error", "could not generate token")
		return
	}
	token := hex.EncodeToString(raw)
	sum := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(sum[:])
	pr := models.PasswordReset{
		Email:     u.Email,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(s.Cfg.PasswordResetExpiry),
	}
	_ = s.DB.Where("email = ?", u.Email).Delete(&models.PasswordReset{}).Error
	if err := s.DB.Create(&pr).Error; err != nil {
		httpx.JSONError(c, 500, "server_error", "could not store reset token")
		return
	}
	audit.Log(s.DB, &u.ID, "auth.forgot_password", "user", itoa(u.ID), nil, httpx.ClientIP(c), c.Request.UserAgent())
	if s.Cfg.LogPasswordResetLink {
		c.JSON(http.StatusOK, gin.H{
			"message":     "password reset token (dev only — disable LOG_PASSWORD_RESET_LINK in production)",
			"reset_token": token,
			"expires_in":  int(s.Cfg.PasswordResetExpiry.Seconds()),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "if the account exists, reset instructions will follow"})
}

type resetBody struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=10,max=128"`
}

func (s *Server) ResetPassword(c *gin.Context) {
	var body resetBody
	if err := c.ShouldBindJSON(&body); err != nil {
		httpx.JSONError(c, 400, "validation_error", err.Error())
		return
	}
	sum := sha256.Sum256([]byte(body.Token))
	tokenHash := hex.EncodeToString(sum[:])
	var pr models.PasswordReset
	if err := s.DB.Where("token_hash = ? AND used_at IS NULL", tokenHash).Order("id desc").First(&pr).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			httpx.JSONError(c, 400, "invalid_token", "reset token is invalid or already used")
			return
		}
		httpx.JSONError(c, 500, "server_error", "lookup failed")
		return
	}
	if time.Now().After(pr.ExpiresAt) {
		httpx.JSONError(c, 400, "expired_token", "reset token has expired")
		return
	}
	var u models.User
	if err := s.DB.Where("email = ?", pr.Email).First(&u).Error; err != nil {
		httpx.JSONError(c, 400, "invalid_token", "user no longer exists")
		return
	}
	hash, err := password.Hash(body.NewPassword)
	if err != nil {
		httpx.JSONError(c, 500, "server_error", "password hashing failed")
		return
	}
	now := time.Now()
	if err := s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&pr).Update("used_at", now).Error; err != nil {
			return err
		}
		return tx.Model(&u).Update("password_hash", hash).Error
	}); err != nil {
		httpx.JSONError(c, 500, "server_error", "could not reset password")
		return
	}
	audit.Log(s.DB, &u.ID, "auth.reset_password", "user", itoa(u.ID), nil, httpx.ClientIP(c), c.Request.UserAgent())
	httpx.JSONOK(c, gin.H{"message": "password updated"})
}

func itoa(u uint) string {
	return strconv.FormatUint(uint64(u), 10)
}
