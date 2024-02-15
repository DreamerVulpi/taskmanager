package server

import (
	"log/slog"
	"strings"
	"taskmanager/pkg/db"

	"github.com/gin-gonic/gin"
)

func userIdentity(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" {
		return
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		return
	}
	userId, err := db.ParseToken(headerParts[1])
	if err != nil {
		slog.Warn(err.Error())
		return
	}
	c.Set("id", userId)
}
