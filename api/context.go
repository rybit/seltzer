package api

import (
	"github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

const (
	tokenKey  = "app.token"
	loggerKey = "app.logger"
)

func getLogger(ctx echo.Context) *logrus.Entry {
	obj := ctx.Value(loggerKey)
	if obj == nil {
		return logrus.NewEntry(logrus.StandardLogger())
	}

	return obj.(*logrus.Entry)
}

func getToken(ctx echo.Context) *jwt.Token {
	obj := ctx.Value(tokenKey)
	if obj == nil {
		return nil
	}

	return obj.(*jwt.Token)
}
