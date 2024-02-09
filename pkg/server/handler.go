package server

import (
	"database/sql"
	"log/slog"
	"net/http"
	"taskmanager/pkg/db"

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

func Roots(conn *sql.DB, router *gin.Engine) {
	// TODO: CREATE -> POST,	READ -> GET,	UPDATE -> PUT,		DELETE -> DELETE
	// TODO: HOME PAGE WITH TASKS

	router.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", gin.H{
			"Title":    "Home",
			"Username": "",
		})
		rslt, err := db.GetTasks(conn, 3, 1)
		if err != nil {
			slog.Warn(err.Error())
		} else {
			for _, element := range rslt {
				temp := element.Title + " " + element.Description
				slog.Info(temp)
			}
		}
	})

	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Title":  "Login",
			"Status": true,
		})
	})

	router.GET("/registration", func(c *gin.Context) {
		c.HTML(http.StatusOK, "registration.html", gin.H{
			"Title": "Registration",
		})
	})

	router.POST("/sign_in", func(c *gin.Context) {
		user := db.User{}
		user.Username = c.Request.FormValue("inputUsername")
		user.Password_hash = c.Request.FormValue("inputPassword")
		if db.AutorizationUser(conn, user) {
			// TODO: Get username for authorization
			c.Redirect(http.StatusFound, "/home")
		} else {
			c.Redirect(http.StatusFound, "/login")
		}
	})

	router.POST("/sign_out", func(c *gin.Context) {
		user := db.User{}
		user.Username = c.Request.FormValue("inputUsername")
		user.Password_hash = c.Request.FormValue("inputPassword")
		db.CreateUser(conn, user)
		c.Redirect(http.StatusFound, "/login")
	})
}
