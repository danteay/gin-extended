package middlewares

import (
	"github.com/danteay/ginrest"
	assert "github.com/danteay/openapi-assert"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SwaggerValidatorConfig struct {
	// SkipPaths list of static paths that will not apply the swagger validator middleware
	SkipPaths []string

	// ObjectResponse value for response payload key "object"
	ObjectResponse string

	// Document path to the swagger assert document
	Document string
}

// DefaultSwaggerValidatorConfig holds the base configuration for swagger validator middleware
var (
	objectSwaggerValidator = "middlewares.swaggerValidator"

	DefaultSwaggerValidatorConfig = &SwaggerValidatorConfig{
		ObjectResponse: objectSwaggerValidator,
		Document:       "spec.yml",
	}
)

// SwaggerValidator create a swagger validator middleware with default configuration
func SwaggerValidator() gin.HandlerFunc {
	return SwaggerValidatorWithConfig(DefaultSwaggerValidatorConfig)
}

// SwaggerValidatorWithConfig execute swagger validator middleware with custom configuration
func SwaggerValidatorWithConfig(conf *SwaggerValidatorConfig) gin.HandlerFunc {
	doc, err := assert.LoadFromURI(conf.Document)
	if err != nil {
		panic(err)
	}

	if doc == nil {
		panic("echo: assert middleware requires an openapi-assert document")
	}

	// Load skip paths that will not apply the middleware
	var skip map[string]struct{}

	if length := len(conf.SkipPaths); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range conf.SkipPaths {
			skip[path] = struct{}{}
		}
	}

	// Load assert swagger document
	assert := assert.New(doc)

	return func(c *gin.Context) {
		// Loading standardized response payload with ginrest
		u := c.Request.RequestURI
		r := ginrest.New(u, "").SetGin(c)

		// Validate if the current path should be skipped
		if _, ok := skip[c.Request.URL.Path]; !ok {
			// Validating request schema
			if err := assert.Request(c.Request); err != nil {
				if conf.ObjectResponse == "" {
					conf.ObjectResponse = objectSwaggerValidator
				}

				r.Res(http.StatusBadRequest, ginrest.Payload{
					"object": conf.ObjectResponse,
				}, err.Error())

				c.Abort()
				return
			}
		}

		c.Next()
		return
	}
}
