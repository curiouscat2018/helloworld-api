package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
	"time"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasSuffix(c.Request.URL.Path, "/ping") {
			return
		}

		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		err := ""
		if len(c.Errors) > 0 {
			err = c.Errors.String()
		}

		subLogger := log.With().
			Str("tag", "per-request-trace").
			Str("request-id", c.GetString("request-id")).
			//Str("hostname", config.Config.HostName()).
			Int("status", c.Writer.Status()).
			Str("method", c.Request.Method).
			Str("protocol", c.Request.Proto).
			Str("path", path).
			Str("ip", c.ClientIP()).
			Dur("latency", latency).
			//Str("user-agent", c.Request.UserAgent()).
			Logger()

		switch {
		case c.Writer.Status() >= http.StatusBadRequest && c.Writer.Status() < http.StatusInternalServerError:
			subLogger.Warn().Msg(err)
		case c.Writer.Status() >= http.StatusInternalServerError:
			subLogger.Error().Msg(err)
		default:
			subLogger.Info().Msg(err)
		}
	}
}
