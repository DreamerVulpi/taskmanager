package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"taskmanager/config"
	"taskmanager/db"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	var conf config.AppConfig
	conf.Path = "."
	conf.NameFile = "config"
	conf.TypeFile = "yaml"

	config.LoadConfig(&conf)
	slog.Info(conf.DB.Driver)

	data := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", conf.DB.User, conf.DB.Password, conf.DB.DBname, conf.DB.Sslmode)
	slog.Info(data)
	conn, err := sql.Open(conf.DB.Driver, data)
	if err != nil {
		slog.Warn("Connect is FAILDED :C\n")
		panic(err)
	}
	defer conn.Close()

	router := gin.Default()
	router.Static("/css", "templates/css")
	router.Static("/src", "src")

	files := []string{
		"templates/index.html",
		"templates/registration.html",
		"templates/home.html",
	}
	router.LoadHTMLFiles(files...)

	// FIXME: GARBAGE
	var status bool

	// TODO: CREATE -> POST,	READ -> GET,	UPDATE -> PUT,		DELETE -> DELETE
	// TODO: HOME PAGE WITH TASKS
	// TODO: REFACTORING CODE

	router.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", gin.H{
			"Title":    "Home",
			"Username": "",
		})
		rslt := db.GetTasks(conn, 3, 1)

		for _, element := range rslt {
			temp := element.Title + " " + element.Description
			slog.Info(temp)
		}
	})

	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Title":  "Login",
			"Status": status,
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
			status = true
			// TODO: Get username for authorization
			c.Redirect(http.StatusFound, "/home")
		} else {
			status = false
		}
	})

	router.POST("/sign_out", func(c *gin.Context) {
		user := db.User{}
		user.Username = c.Request.FormValue("inputUsername")
		user.Password_hash = c.Request.FormValue("inputPassword")
		db.CreateUser(conn, user)
		c.Redirect(http.StatusFound, "/login")
	})

	router.Run(":" + conf.S.Port)

}
