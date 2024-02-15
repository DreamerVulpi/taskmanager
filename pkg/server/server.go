package server

import (
	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	router := gin.Default()
	router.Static("/css", "pkg/templates/css")
	router.Static("/src", "pkg/src")

	html := []string{
		"pkg/templates/index.html",
		"pkg/templates/registration.html",
		"pkg/templates/home.html",
	}
	router.LoadHTMLFiles(html...)
	return router
}

func Run(router *gin.Engine, port string) *gin.Engine {
	router.Run(":" + port)
	return router
}
