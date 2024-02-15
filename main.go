package main

import (
	"fmt"
	"log/slog"
	"taskmanager/config"
	"taskmanager/pkg/db"
	"taskmanager/pkg/server"

	_ "github.com/lib/pq"
)

func main() {
	var conf config.AppConfig
	conf.Path = "."
	conf.NameFile = "config"
	conf.TypeFile = "yaml"
	config.Load(&conf)
	slog.Info(conf.DB.Driver)

	params := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", conf.DB.User, conf.DB.Password, conf.DB.DBname, conf.DB.Sslmode)
	conn, err := db.ConnectToDB(params, conf.DB.Driver)
	if err != nil {
		slog.Warn(err.Error())
	}
	host := server.Init()
	server.Roots(conn, host)
	server.Run(host, conf.S.Port)
	defer conn.Close()
}
