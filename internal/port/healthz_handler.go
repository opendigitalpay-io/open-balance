package port

import (
	"github.com/gin-gonic/gin"
)

func (h *HTTPServer) GetHealthStatus() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		h.RespondWithOK(ctx, gin.H{"status": "ok"})
	}
}
