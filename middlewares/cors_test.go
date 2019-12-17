package middlewares

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/assert.v1"
)

func TestCors(t *testing.T) {
	router := newTestRouter(CorsWithConfig(DefaultCorsConfig))
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "get")
	})

	w := performRequestOrigin(router, "GET", "/ping", "http://google.com")
	assert.Equal(t, "get", w.Body.String())
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
}
