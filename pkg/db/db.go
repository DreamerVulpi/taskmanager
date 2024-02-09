package db

import (
	"database/sql"
	"log/slog"

	_ "github.com/lib/pq"
)

type User struct {
	Username      string
	Password_hash string
}

// TODO: CRUD TASKS

func ConnectToDB(params string, driver string) *sql.DB {
	slog.Info(params)
	conn, err := sql.Open(driver, params)
	if err != nil {
		slog.Warn("Connect to DB is failded :C\n")
		panic(err)
	}
	return conn
}

func AutorizationUser(conn *sql.DB, user User) bool {
	var status bool
	slog.Info(user.Username)
	slog.Info(user.Password_hash)
	conn.QueryRow("SELECT EXISTS(SELECT username = $1 FROM users WHERE password_hash = $2)", user.Username, user.Password_hash).Scan(&status)
	if status {
		slog.Info("Autorization user is successful")
		return true
	} else {
		slog.Warn("Autorization user is failded :C")
		return false
	}
}
