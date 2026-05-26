package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/user/chat-app/internal/logger"
)

func TestLoggerMiddleware_GeneratesAndPropagatesRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(LoggerMiddleware(nil)) // nil base log is safe for test mode
	r.GET("/test", func(c *gin.Context) {
		// Verify context injection
		reqID := c.Request.Context().Value(logger.RequestIDKey)
		assert.NotNil(t, reqID, "request_id should be injected into context")
		assert.NotEmpty(t, reqID.(string), "request_id should not be empty")

		// Verify header response
		headerID := c.Writer.Header().Get("X-Request-ID")
		assert.Equal(t, reqID.(string), headerID, "X-Request-ID header should match context value")
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLoggerMiddleware_PreservesExistingRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(LoggerMiddleware(nil))
	r.GET("/test", func(c *gin.Context) {
		reqID := c.Request.Context().Value(logger.RequestIDKey)
		assert.Equal(t, "pre-existing-id-123", reqID.(string))
		assert.Equal(t, "pre-existing-id-123", c.Writer.Header().Get("X-Request-ID"))
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Request-ID", "pre-existing-id-123")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
