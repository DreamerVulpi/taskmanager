package db

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"log/slog"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	_ "github.com/lib/pq"
)

const (
	jwtSecretKey = "fdasgdsg532413dsaf"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"id"`
}

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password_hash"`
	Token    string `json:"token"`
}

func ConnectToDB(params string, driver string) (*sql.DB, error) {
	slog.Info(params)
	conn, err := sql.Open(driver, params)
	if err != nil {
		return nil, err
	}
	return conn, err
}

func HashPassword(password string) ([]byte, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(password))
	if err != nil {
		return []byte{}, err
	}
	rslt := hash.Sum(nil)
	return rslt, nil
}

func GetUser(conn *sql.DB, user User) (int, error) {
	password_hash, err := HashPassword(user.Password)
	if err != nil {
		return 0, err
	}
	row, err := conn.Query("SELECT id FROM users WHERE username=$1 AND password_hash=$2", user.Username, password_hash)
	result := 0
	if err != nil {
		return 0, err
	} else {
		row.Scan(&result)
		return result, err
	}
}

func GenerateToken(user User) (string, error) {
	payload := jwt.MapClaims{
		"id":  user.Id,
		"sub": user.Username,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		slog.Warn("JWT token signing")
		return "", err
	}
	return t, err
}

func ParseToken(token string) (int, error) {
	tkn, err := jwt.ParseWithClaims(token, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, err := token.Method.(*jwt.SigningMethodHMAC); !err {
			return nil, errors.New("invalid signing method")
		}
		return []byte(jwtSecretKey), nil
	})
	if err != nil {
		slog.Warn(err.Error())
		return 0, err
	}
	claims := tkn.Claims.(*tokenClaims)
	return claims.UserId, nil

}

func AutorizationUser(conn *sql.DB, user User) (string, error) {
	slog.Info(user.Username)
	slog.Info(user.Password)
	id, err := GetUser(conn, user)
	if err != nil {
		return "", err
	}
	user.Id = strconv.Itoa(id)
	str, err := GenerateToken(user)
	if err != nil {
		return "", err
	}
	slog.Info("Success!")
	return str, err
}

func CreateUser(conn *sql.DB, user User) error {
	slog.Info(user.Username)
	slog.Info(user.Password)
	password_hash, err := HashPassword(user.Password)
	if err != nil {
		return err
	}
	_, err = conn.Exec("INSERT INTO users (username, password_hash) VALUES ($1, $2)", user.Username, password_hash)
	return err
}
