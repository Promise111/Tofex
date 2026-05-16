package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Health returns service liveness.
//
//	@Summary		Health check
//	@Description	Returns whether the API process is up.
//	@Tags			system
//	@Produce		json
//	@Success		200	{object}	apidocs.HealthResponse
//	@Router			/healthz [get]
func (s *Server) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
