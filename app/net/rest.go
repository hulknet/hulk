package net

import (
	"net/http"

	"github.com/labstack/echo/v4"

	http2 "github.com/kotfalya/hulk/app/http"
	"github.com/kotfalya/hulk/app/types"
)

type Rest struct {
	echo   *echo.Echo
	net    *Net
	addr   string
	secret interface{}
}

func NewRestServer(net *Net, addr string, secret interface{}) *Rest {
	r := &Rest{
		echo:   echo.New(),
		net:    net,
		addr:   addr,
		secret: secret,
	}

	r.echo.Use(http2.RegisterJWT(secret))

	r.echo.GET("/status", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, echo.Map{
			"status": "OK",
		})
	})

	r.echo.GET("/token", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, echo.Map{
			"id":   types.ToHex(http2.ServiceFromContext(ctx).ID),
			"type": http2.TokenToString(http2.ServiceFromContext(ctx).Type),
		})
	})

	r.echo.GET("/self", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, echo.Map{
			"pk":    types.ToHex(r.net.self.PK),
			"addr":  r.net.self.PK.ID().Addr(),
			"token": types.ToHex(r.net.self.Token),
		})
	})

	return r
}

func (r *Rest) Listen() error {
	return r.echo.Start(r.addr)
}
