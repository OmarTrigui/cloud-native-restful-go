package middleware

import (
	"github.com/gin-gonic/gin"
)

func ContentTypeJson() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json;charset=utf8")
		c.Next()
	}
}
