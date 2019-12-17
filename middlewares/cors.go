package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// CorsConfig base configuration for cors request
type CorsConfig struct {
	// SkipPaths list of static paths that will not apply the middleware
	SkipPaths []string

	// RegexpSkipPaths regular expresions for http routes that will skip the middleware
	RegexpSkipPaths []string

	// AllowOrigin origin that are allowed to request the service
	AllowOrigin string

	// AllowMethods methods that are allowed to be requested
	AllowMethods []string

	// AllowHeaders allowed headers
	AllowHeaders []string

	// AllowCredentials validate allowed credentials
	AllowCredentials string
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
		AllowOrigin:      "*",
		AllowCredentials: "true",
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     DefaultCorsHeaders,
		SkipPaths:        []string{},
		RegexpSkipPaths:  []string{},
	}
)

// Cors retrive a default cors middleware
func Cors() gin.HandlerFunc {
	return CorsWithConfig(DefaultCorsConfig)
}

// CorsWithConfig retrive a cors middleware with custom configuration
func CorsWithConfig(conf *CorsConfig) gin.HandlerFunc {
	if conf.AllowOrigin == "" {
		conf.AllowOrigin = DefaultCorsConfig.AllowOrigin
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

		if _, ok := skip[path]; !ok && !skipRegexpPath(conf.RegexpSkipPaths, path) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", conf.AllowOrigin)

			c.Writer.Header().Set("Access-Control-Allow-Credentials", conf.AllowCredentials)

			methods := strings.Join(conf.AllowMethods, ",")
			c.Writer.Header().Set("Access-Control-Allow-Headers", methods)

			headers := strings.Join(conf.AllowHeaders, ",")
			c.Writer.Header().Set("Access-Control-Allow-Headers", headers)
		}

		c.Next()
	}
}
