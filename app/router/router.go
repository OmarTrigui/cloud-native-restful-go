package router

import (
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"

	"myapp/app/app"
	"myapp/app/router/middleware"
)

func New(a *app.App) *gin.Engine {

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: false,
		},
	)

	router := gin.New()

	router.Use(logger.SetLogger())
	router.Use(gin.Recovery())

	// Routes for healthz
	router.GET("/healthz/liveness", app.HandleLive)
	router.GET("/healthz/readiness", a.HandleReady)

	// Routes for APIs
	v1 := router.Group("/api/v1")
	{
		router.Use(middleware.ContentTypeJson())

		// Routes for books
		v1.GET("/books", a.HandleListBooks)
		v1.POST("/books", a.HandleCreateBook)
		v1.GET("/books/:id", a.HandleReadBook)
		v1.PUT("/books/:id", a.HandleUpdateBook)
		v1.DELETE("/books/:id", a.HandleDeleteBook)
	}

	router.GET("/", a.HandleIndex)

	return router
}
