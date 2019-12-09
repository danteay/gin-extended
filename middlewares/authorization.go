package middlewares

import (
	"encoding/base64"
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/danteay/ginrest"
	"github.com/gin-gonic/gin"
)

const (
	// BearerAuth bearer authentication strategy
	BearerAuth = "bearer"

	// BasicAuth basic authentication strategy
	BasicAuth = "basic"
)

// AuthValidator callback validation for authentication
type AuthValidator func(key string) bool

// AuthorizationConfig Configuration for Api key authentication
type AuthorizationConfig struct {
	// Type Authorization types to use in validation
	// - bearer
	// - basic
	Type string

	// APIKey static validation for bearer authentication
	APIKey string

	// AuthCredentials static credentials for basic authentication
	// format []string{"user", "password"}
	AuthCredentials []string

	// Validator custom header validation
	Validator AuthValidator

	// Header custom header for authentication
	Header string

	// SkipPaths static paths that will not apply the middleware
	SkipPaths []string

	// RegexpSkipPaths regular expresions for http routes that will skip the middleware
	RegexpSkipPaths []string

	// ObjectResponse value for response payload key "object"
	ObjectResponse string
}

var (
	errInvalidAuth = errors.New("invalid authentication")

	errInvalidAuthMethod = errors.New("invalid authentication method")

	objectAuthentication = "middlewares.authentication"

	// DefaultAuthorization Default configuration for Authorization
	DefaultAuthorization = &AuthorizationConfig{
		Type:            "bearer",
		APIKey:          "",
		AuthCredentials: nil,
		Validator:       nil,
		Header:          "Authorization",
		SkipPaths:       []string{"/ping"},
		ObjectResponse:  objectAuthentication,
	}
)

// Authorization Return an Authentication middleware with default config
func Authorization() gin.HandlerFunc {
	return AuthorizationWithConfig(DefaultAuthorization)
}

// AuthorizationWithConfig Authenticate request by apikey header
func AuthorizationWithConfig(conf *AuthorizationConfig) gin.HandlerFunc {
	if conf.Header == "" {
		conf.Header = "Authorization"
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

		var err error

		if _, ok := skip[path]; !ok && !conf.skipRegexpPath(path) {
			switch conf.Type {
			case "bearer":
				err = conf.bearerAuth(c)
				break
			case "basic":
				err = conf.basicAuth(c)
				break
			default:
				err = errInvalidAuthMethod
				break
			}
		}

		if err != nil {
			u := c.Request.RequestURI
			r := ginrest.New(u, "").SetGin(c)

			if conf.ObjectResponse == "" {
				conf.ObjectResponse = objectAuthentication
			}

			data := ginrest.Payload{"object": conf.ObjectResponse}

			r.Res(http.StatusForbidden, data, "error")
			c.Abort()
			return
		}

		c.Next()
		return
	}
}

func (ac *AuthorizationConfig) bearerAuth(c *gin.Context) error {
	token := c.GetHeader(ac.Header)
	token = strings.ReplaceAll(token, "Bearer ", "")

	if ac.Validator != nil {
		if !ac.Validator(token) {
			return errInvalidAuth
		}
	} else {
		if ac.APIKey != token {
			return errInvalidAuth
		}
	}

	return nil
}

func (ac *AuthorizationConfig) basicAuth(c *gin.Context) error {
	token := c.GetHeader(ac.Header)
	token = strings.ReplaceAll(token, "Basic ", "")

	if ac.Validator != nil {
		if !ac.Validator(token) {
			return errInvalidAuth
		}
	} else {
		if ac.AuthCredentials == nil {
			return nil
		}

		data, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			return errInvalidAuth
		}

		args := strings.Split(string(data), ":")

		if ac.AuthCredentials[0] != args[0] || ac.AuthCredentials[1] != args[1] {
			return errInvalidAuth
		}
	}

	return nil
}

func (ac *AuthorizationConfig) skipRegexpPath(path string) bool {
	for _, reg := range ac.RegexpSkipPaths {
		exp := regexp.MustCompile(reg)

		if exp.MatchString(path) {
			return true
		}
	}

	return false
}
