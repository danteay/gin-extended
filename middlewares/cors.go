package middlewares

import (
	"github.com/gin-gonic/gin"
	"strings"
)

// CorsConfig base configuration for cors request
type CorsConfig struct {
	// SkipPaths list of static paths that will not apply the middleware
	SkipPaths []string

	// RegexpSkipPaths regular expresions for http routes that will skip the middleware
	RegexpSkipPaths []string

	// AllowOrigins origins that are allowed to request the service
	AllowOrigins []string

	// AllowMethods methods that are allowed to be requested
	AllowMethods []string

	// AllowHeaders allowed headers
	AllowHeaders []string
}

var (
	// DefaultCorsHeaders list of default allowed headers for request
	DefaultCorsHeaders = []string{
		"Authentication",
		"Content-Type",
		"Origin",
		"Accept",
		"Cache-Control",
		"Postman-Token",
		"User-Agent",
		"Cache-Control",
		"Host",
		"Accept-Encoding",
		"Connection",
	}

	// DefaultCorsConfig default configuration for cors middleware
	DefaultCorsConfig = &CorsConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: DefaultCorsHeaders,
	}
)

// Cors retrive a default cors middleware
func Cors() gin.HandlerFunc {
	return CorsWithConfig(DefaultCorsConfig)
}

// CorsWithConfig retrive a cors middleware with custom configuration
func CorsWithConfig(conf *CorsConfig) gin.HandlerFunc {
	if conf.AllowOrigins == nil {
		conf.AllowOrigins = DefaultCorsConfig.AllowOrigins
	}

	if conf.AllowMethods == nil {
		conf.AllowMethods = DefaultCorsConfig.AllowMethods
	}

	if conf.AllowHeaders == nil {
		conf.AllowHeaders = DefaultCorsConfig.AllowHeaders
	}

	var skip map[string]struct{}

	if length := len(conf.SkipPaths); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range conf.SkipPaths {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		c.Next()

		if _, ok := skip[path]; !ok && !skipRegexpPath(conf.RegexpSkipPaths, path) {
			origins := strings.Join(conf.AllowOrigins, ",")
			c.Request.Response.Header.Set("Access-Control-Allow-Origin", origins)

			methods := strings.Join(conf.AllowMethods, ",")
			c.Request.Response.Header.Set("Access-Control-Allow-Headers", methods)

			headers := strings.Join(conf.AllowHeaders, ",")
			c.Request.Response.Header.Set("Access-Control-Allow-Headers", headers)
		}
	}
}
