package net

import (
	"net/http"

	"github.com/labstack/echo/v4"

	libHttp "github.com/hulknet/hulk/app/http"
	"github.com/hulknet/hulk/app/types"
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

	r.echo.GET("/block/:blockId", func(ctx echo.Context) error {
		blockId, err := types.ID64FromHex(ctx.Param("blockId"))
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}

		net, ok := r.netCont.blockToNet[blockId]
		if !ok {
			return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "unknown block id"})
		}

		return ctx.JSON(http.StatusOK, echo.Map{
			"pub":     net.State().Peer().Pub.ID256().Hex(),
			"id":      net.State().Peer().Pub.ID().Hex(),
			"blockId": net.State().ID().Hex(),
			"token":   net.State().Peer().Token.Hex(),
			"time":    net.State().Time().Hex(),
		})
	})

	r.echo.GET("/block", func(ctx echo.Context) error {
		resp := echo.Map{}
		for id, net := range r.netCont.blockToNet {
			resp[id.Hex()] = net.state.Block().Status.String()
		}

		return ctx.JSON(http.StatusOK, resp)
	})

	return r
}

func (r *Rest) Listen() error {
	return r.echo.Start(r.addr)
}
