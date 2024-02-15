package server

import (
	"database/sql"
	"log/slog"
	"net/http"
	"taskmanager/pkg/db"

	"github.com/gin-gonic/gin"
)

func Roots(conn *sql.DB, router *gin.Engine) {
	api := router.Group("/api", userIdentity)
	{
		api.GET("/home", func(c *gin.Context) {
			c.HTML(http.StatusOK, "home.html", gin.H{
				"Title": "Home",
			})
			// rslt, err := db.GetTasks(conn, 3, 1)
		})
	}

	auth := router.Group("/auth")
	{
		auth.GET("/login", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", gin.H{
				"Title": "Login",
			})
		})
		auth.GET("/registration", func(c *gin.Context) {
			c.HTML(http.StatusOK, "registration.html", gin.H{
				"Title": "Registration",
			})
		})
		auth.POST("/log_in", func(c *gin.Context) {
			user := db.User{}
			user.Username = c.Request.FormValue("inputUsername")
			user.Password = c.Request.FormValue("inputPassword")
			token, err := db.AutorizationUser(conn, user)
			slog.Info(token)
			if err != nil {
				slog.Warn(err.Error())
				c.Redirect(http.StatusFound, "/auth/login")
			}
			c.JSON(http.StatusOK, map[string]interface{}{
				"token": token,
			})
			c.Redirect(http.StatusFound, "/home")
		})
		auth.POST("/register", func(c *gin.Context) {
			user := db.User{}
			user.Username = c.Request.FormValue("inputUsername")
			user.Password = c.Request.FormValue("inputPassword")
			err := db.CreateUser(conn, user)
			if err != nil {
				slog.Warn(err.Error())
			}
			c.Redirect(http.StatusFound, "/auth/login")
		})
	}
}
