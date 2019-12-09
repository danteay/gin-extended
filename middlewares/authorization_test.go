package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/assert.v1"
)

func newTestRouter(mw gin.HandlerFunc) *gin.Engine {
	router := gin.New()
	router.Use(mw)
	return router
}

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func performRequestWithHeaders(r http.Handler, method, path string, headers map[string]string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

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
