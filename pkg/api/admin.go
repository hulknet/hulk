package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (r *Rest) createNet(ctx echo.Context) error {
	err := r.host.CreateNet()
	if err != nil {
		return err
	}
	return ctx.NoContent(http.StatusCreated)
}
