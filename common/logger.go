package common

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	log.Logger = zerolog.New(os.Stdout)
}

func TraceInfo(c *gin.Context) *zerolog.Event {
	return trace(c).Info()
}

func TraceWarn(c *gin.Context) *zerolog.Event {
	return trace(c).Warn()
}
func TraceError(c *gin.Context) *zerolog.Event {
	return trace(c).Error()
}

func TraceFatal(c *gin.Context) *zerolog.Event {
	return trace(c).Fatal()
}

func trace(c *gin.Context) *zerolog.Logger {
	requestID := ""
	if c != nil {
		requestID = c.GetString("request-id")
	}
	temp := log.With().
		Str("tag", "diagnostic-trace").
		Str("request-id", requestID).
		Logger()
	return &temp
}
