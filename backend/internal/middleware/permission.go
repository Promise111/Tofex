package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/tofex/backend/internal/httpx"
	"github.com/tofex/backend/internal/rbac"
)

func RequirePermission(perm string) gin.HandlerFunc {
	return func(c *gin.Context) {
		u := CurrentUser(c)
		if u == nil || !rbac.UserHasPermission(u, perm) {
			httpx.JSONError(c, 403, "forbidden", "insufficient permissions")
			c.Abort()
			return
		}
		c.Next()
	}
}
