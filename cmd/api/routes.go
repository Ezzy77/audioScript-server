package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (app *application) routes() *gin.Engine {

	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	router.Use(cors.New(corsConfig))

	router.POST("/v1/api/transcript", app.createTranscriptJobHandler)
	router.POST("/v1/api/upload", app.uploadMedia)

	// user auth
	//router.POST("/v1/api/auth/login", app.loginHandler)
	//router.POST("/v1/api/auth/logout", app.logoutHandler)
	//router.POST("/v1/api/auth/register", app.registerHandler)

	return router
}
