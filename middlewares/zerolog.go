package middlewares

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	// HeaderXRequestID  Request ID
	HeaderXRequestID = "X-Request-ID"

	// HeaderReferer Return URL from referer if is send
	HeaderReferer = "Referer"
)

// ZerologConfig defines the config for ZeroLog middleware.
type ZerologConfig struct {
	// FieldMap set a list of fields with tags
	//
	// Tags to constructed the logger fields.
	//
	// - @id (Request ID)
	// - @remote_ip
	// - @uri
	// - @host
	// - @method
	// - @path
	// - @referer
	// - @user_agent
	// - @status
	// - @latency (In nanoseconds)
	// - @latency_human (Human readable)
	// - @bytes_in (Bytes received)
	// - @bytes_out (Bytes sent)
	// - @header:<NAME>
	// - @query:<NAME>
	// - @form:<NAME>
	// - @cookie:<NAME>
	FieldMap map[string]string

	// Logger it is a zerolog logger
	Logger zerolog.Logger

	// SkipPaths list of http routes that will not apply the middleware
	SkipPaths []string
}

// DefaultZeroLogConfig is the default ZeroLog middleware config.
var DefaultZeroLogConfig = &ZerologConfig{
	FieldMap: map[string]string{
		"remote_ip": "@remote",
		"uri":       "@uri",
		"host":      "@host",
		"method":    "@method",
		"status":    "@status",
		"latency":   "@latency",
	},
	Logger: log.Output(zerolog.ConsoleWriter{Out: os.Stderr}),
}

// Zerolog Create a default ZeroLog middleware
func Zerolog() gin.HandlerFunc {
	return ZerologWithConfig(DefaultZeroLogConfig)
}

// ZerologWithConfig Create ZeroLog middleware for logging
func ZerologWithConfig(conf *ZerologConfig) gin.HandlerFunc {
	var skip map[string]struct{}

	if length := len(conf.SkipPaths); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range conf.SkipPaths {
			skip[path] = struct{}{}
		}
	}

	log.Logger = conf.Logger

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		if _, ok := skip[path]; !ok {
			stop := time.Now()
			entry := conf.Logger.Info()

			if raw != "" {
				path = path + "?" + raw
			}

			for k, v := range conf.FieldMap {
				if v == "" {
					continue
				}

				switch v {
				case "@id":
					id := c.GetHeader(HeaderXRequestID)

					if id == "" {
						id = c.GetHeader(HeaderXRequestID)
					}

					entry = entry.Str(k, id)
				case "@remote":
					entry = entry.Str(k, c.ClientIP())
				case "@uri":
					entry = entry.Str(k, c.Request.RequestURI)
				case "@host":
					entry = entry.Str(k, c.Request.Host)
				case "@method":
					entry = entry.Str(k, c.Request.Method)
				case "@path":
					p := path

					if p == "" {
						p = "/"
					}

					entry = entry.Str(k, p)
				case "@referer":
					entry = entry.Str(k, c.Request.Referer())
				case "@user_agent":
					entry = entry.Str(k, c.Request.UserAgent())
				case "@status":
					entry = entry.Int(k, c.Writer.Status())
				case "@latency":
					l := stop.Sub(start)
					entry = entry.Str(k, strconv.FormatInt(int64(l), 10))
				case "@latency_human":
					entry = entry.Str(k, stop.Sub(start).String())
				case "@bytes_in":
					entry = entry.Int64(k, c.Request.ContentLength)
				case "@bytes_out":
					entry = entry.Int(k, c.Writer.Size())
				default:
					switch {
					case strings.HasPrefix(v, "@header:"):
						entry = entry.Str(k, c.GetHeader(v[8:]))
					case strings.HasPrefix(v, "@query:"):
						entry = entry.Str(k, c.Query(v[7:]))
					case strings.HasPrefix(v, "@form:"):
						entry = entry.Str(k, c.PostForm(v[6:]))
					case strings.HasPrefix(v, "@cookie:"):
						cookie, err := c.Cookie(v[8:])
						if err == nil {
							entry = entry.Str(k, cookie)
						}
					}
				}
			}

			entry.Msg("request")
		}
	}
}
