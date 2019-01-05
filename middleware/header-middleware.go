package middleware

import (
	"github.com/gin-gonic/gin"
)

func  HeaderMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("request-id")
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}
