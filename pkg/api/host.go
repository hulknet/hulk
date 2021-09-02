package api

import (
	"net/http"

	"github.com/kotfalya/hulk/pkg/host"
	"github.com/labstack/echo/v4"
)

func (r *Rest) hostStatus(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, host.NewStatusModel(r.host))
}

func (r *Rest) hostConnect(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "Connected")
}
