package api

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

const (
	tokenKey  = "app.token"
	userKey   = "app.user"
	loggerKey = "app.logger"
)

func UserFromContext(ctx echo.Context) *User {
	user := ctx.Get(userKey)
	return user.(*User)
}

func LoggerFromContext(ctx echo.Context) *logrus.Entry {
	obj := ctx.Get(loggerKey)
	if obj == nil {
		return logrus.NewEntry(logrus.StandardLogger())
	}

	return obj.(*logrus.Entry)
}
