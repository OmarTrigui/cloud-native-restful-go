package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (app *App) HandleIndex(c *gin.Context) {
	c.Writer.Header().Set("Content-Length", "12")
	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Writer.Header().Set("X-Content-Type-Options", "nosniff")

	c.Writer.WriteHeader(http.StatusOK)

	c.Writer.Write([]byte("Hello World!"))
}
