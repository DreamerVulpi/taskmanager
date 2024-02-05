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

type inputData struct {
	Objective string
	Deadline  string
}

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
	router.LoadHTMLFiles("templates/index.html")

	// TODO: CREATE -> POST,	READ -> GET,	UPDATE -> PUT,		DELETE -> DELETE
	// TODO: DO SIMPLE UI IN HTML+CSS

	router.GET("", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Title": "Home Page",
			"Rows":  db.READ(conn),
		})

	})

	router.Run(":" + conf.S.Port)

}
