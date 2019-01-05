package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

func ContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := uuid.NewV4()
		c.Set("request-id", id.String())
		c.Next()
	}
}
