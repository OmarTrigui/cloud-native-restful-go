package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// HandleLive is an http.HandlerFunc that handles liveness checks by
// immediately responding with an HTTP 200 status.
func HandleLive(c *gin.Context) {
	writeHealthy(c)
}

// HandleReady is an http.HandlerFunc that handles readiness checks by
// responding with an HTTP 200 status if it is healthy, 500 otherwise.
func (app *App) HandleReady(c *gin.Context) {
	if err := app.db.DB().Ping(); err != nil {
		app.Logger().Fatal().Err(err).Msg("")
		writeUnhealthy(c)
		return
	}

	writeHealthy(c)
}

func writeHealthy(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write([]byte("."))
}

func writeUnhealthy(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.WriteHeader(http.StatusInternalServerError)
	c.Writer.Write([]byte("."))
}
