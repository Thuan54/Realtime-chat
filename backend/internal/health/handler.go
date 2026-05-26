package health

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Checker holds dependency ping functions injected at startup
type Checker struct {
	DBPing    func(ctx context.Context) error
	RedisPing func(ctx context.Context) error
}

// Health returns 200 OK if the application process is running
func (h *Checker) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Ready returns 200 only if PostgreSQL and Redis respond within timeout
func (h *Checker) Ready(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	if h.DBPing != nil {
		if err := h.DBPing(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "database_unreachable",
				"error":  err.Error(),
			})
			return
		}
	}

	if h.RedisPing != nil {
		if err := h.RedisPing(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "redis_unreachable",
				"error":  err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
