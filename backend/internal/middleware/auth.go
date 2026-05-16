package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tofex/backend/internal/config"
	"github.com/tofex/backend/internal/httpx"
	"github.com/tofex/backend/internal/models"
	jwtutil "github.com/tofex/backend/internal/security/jwt"
	"gorm.io/gorm"
)

const CtxUserKey = "auth_user"

func Auth(cfg *config.Config, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		const pfx = "Bearer "
		if len(h) < len(pfx) || !strings.EqualFold(h[:len(pfx)], pfx) {
			httpx.JSONError(c, 401, "unauthorized", "missing or invalid bearer token")
			c.Abort()
			return
		}
		raw := strings.TrimSpace(h[len(pfx):])
		uid, err := jwtutil.ParseUserID(raw, cfg.JWTSecret)
		if err != nil {
			httpx.JSONError(c, 401, "unauthorized", "invalid or expired token")
			c.Abort()
			return
		}
		var user models.User
		if err := db.Preload("Roles.Permissions").First(&user, uid).Error; err != nil {
			httpx.JSONError(c, 401, "unauthorized", "user not found")
			c.Abort()
			return
		}
		c.Set(CtxUserKey, &user)
		c.Next()
	}
}

func CurrentUser(c *gin.Context) *models.User {
	v, ok := c.Get(CtxUserKey)
	if !ok {
		return nil
	}
	u, _ := v.(*models.User)
	return u
}
