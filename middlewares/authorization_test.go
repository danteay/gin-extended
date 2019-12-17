package middlewares

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/assert.v1"
)

func TestSkipRegexpPaths(t *testing.T) {
	router := newTestRouter(AuthorizationWithConfig(&AuthorizationConfig{
		Type:   BearerAuth,
		APIKey: "123456789",
		RegexpSkipPaths: []string{
			`^\/ping`,
		},
	}))

	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "/ping")
	})

	router.GET("/ping/dos", func(c *gin.Context) {
		c.String(http.StatusOK, "/ping/dos")
	})

	resp := performRequest(router, "GET", "/ping")
	assert.Equal(t, "/ping", resp.Body.String())

	resp = performRequest(router, "GET", "/ping/dos")
	assert.Equal(t, "/ping/dos", resp.Body.String())
}
