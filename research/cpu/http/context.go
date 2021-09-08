package http

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

const (
	tokenKey   = "app.token"
	serviceKey = "app.service"
	loggerKey  = "app.logger"
)

func ServiceFromContext(ctx echo.Context) *Service {
	user := ctx.Get(serviceKey)
	return user.(*Service)
}

func LoggerFromContext(ctx echo.Context) *logrus.Entry {
	obj := ctx.Get(loggerKey)
	if obj == nil {
		return logrus.NewEntry(logrus.StandardLogger())
	}

	return obj.(*logrus.Entry)
}
