package middlewares

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
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

func performRequestOrigin(r http.Handler, method, path, origin string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	if len(origin) > 0 {
		req.Header.Set("Origin", origin)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func performRequestWithHeadersOrigin(r http.Handler, method, path, origin string, headers map[string]string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if len(origin) > 0 {
		req.Header.Set("Origin", origin)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
