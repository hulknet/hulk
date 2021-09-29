package net

import (
	"net/http"

	"github.com/labstack/echo/v4"

	libHttp "github.com/kotfalya/hulk/app/http"
	"github.com/kotfalya/hulk/app/types"
)

type Rest struct {
	echo    *echo.Echo
	netCont *Container
	addr    string
	secret  interface{}
}

func NewRestServer(netCont *Container, addr string, secret interface{}) *Rest {
	r := &Rest{
		echo:    echo.New(),
		netCont: netCont,
		addr:    addr,
		secret:  secret,
	}

	r.echo.Use(libHttp.RegisterJWT(secret))

	r.echo.GET("/status", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, echo.Map{
			"status": "OK",
		})
	})

	r.echo.GET("/token", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, echo.Map{
			"id":   types.ID256ToHex(libHttp.ServiceFromContext(ctx).ID),
			"type": libHttp.TokenToString(libHttp.ServiceFromContext(ctx).Type),
		})
	})

	r.echo.GET("/self", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, echo.Map{
			//"pk":       types.ID256ToHex(r.net.self.Pub.ID256()),
			//"pkPrefix": hex.EncodeToString(r.net.self.Pub.ID().Bytes()),
			//"addr":     hex.EncodeToString(r.net.self.Pub.ID256().ID().Bytes()),
			//"token":    types.ID256ToHex(r.net.self.Token),
		})
	})

	return r
}

func (r *Rest) Listen() error {
	return r.echo.Start(r.addr)
}
