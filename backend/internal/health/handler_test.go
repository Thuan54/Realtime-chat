package health

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealth_AlwaysReturnsOK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	checker := &Checker{}
	r := gin.New()
	r.GET("/health", checker.Health)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"status":"ok"}`, w.Body.String())
}

func TestReady_Returns200_WhenAllDependenciesHealthy(t *testing.T) {
	gin.SetMode(gin.TestMode)
	checker := &Checker{
		DBPing:    func(ctx context.Context) error { return nil },
		RedisPing: func(ctx context.Context) error { return nil },
	}
	r := gin.New()
	r.GET("/ready", checker.Ready)

	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestReady_Returns503_WhenDatabaseFails(t *testing.T) {
	gin.SetMode(gin.TestMode)
	checker := &Checker{
		DBPing: func(ctx context.Context) error { return assert.AnError },
	}
	r := gin.New()
	r.GET("/ready", checker.Ready)

	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	assert.Contains(t, w.Body.String(), "database_unreachable")
}

func TestReady_Returns503_WhenRedisFails(t *testing.T) {
	gin.SetMode(gin.TestMode)
	checker := &Checker{
		DBPing:    func(ctx context.Context) error { return nil },
		RedisPing: func(ctx context.Context) error { return assert.AnError },
	}
	r := gin.New()
	r.GET("/ready", checker.Ready)

	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	assert.Contains(t, w.Body.String(), "redis_unreachable")
}
