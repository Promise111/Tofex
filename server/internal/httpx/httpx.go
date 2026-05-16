package httpx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func JSONError(c *gin.Context, status int, code, message string) {
	c.JSON(status, gin.H{"error": gin.H{"code": code, "message": message}})
}

func JSONOK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, data)
}

func ClientIP(c *gin.Context) string {
	if x := c.GetHeader("X-Forwarded-For"); x != "" {
		return x
	}
	return c.ClientIP()
}
