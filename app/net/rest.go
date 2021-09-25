package net

import (
	"encoding/hex"
	"net/http"

	"github.com/labstack/echo/v4"

	libHttp "github.com/kotfalya/hulk/app/http"
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

	r.echo.Use(libHttp.RegisterJWT(secret))

	r.echo.GET("/status", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, echo.Map{
			"status": "OK",
		})
	})

	r.echo.GET("/token", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, echo.Map{
			"id":   types.ToHex(libHttp.ServiceFromContext(ctx).ID),
			"type": libHttp.TokenToString(libHttp.ServiceFromContext(ctx).Type),
		})
	})

	r.echo.GET("/self", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, echo.Map{
			"pk":       types.ToHex(r.net.self.PK),
			"pkPrefix": hex.EncodeToString(r.net.self.PK.Prefix().Bytes()),
			"addr":     hex.EncodeToString(r.net.self.PK.ID().Prefix().Bytes()),
			"token":    types.ToHex(r.net.self.Token),
		})
	})

	return r
}

func (r *Rest) Listen() error {
	return r.echo.Start(r.addr)
}
