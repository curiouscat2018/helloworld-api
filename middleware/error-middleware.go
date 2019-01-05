package middleware

import (
	"github.com/gin-gonic/gin"
)

func  ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// skip if no errors
		if c.Errors.Last() == nil {
			return
		}

		c.JSON(c.Writer.Status(),  gin.H{
			"request-id": c.GetString("request-id"),
			"err": c.Errors.String()})
	}
}
