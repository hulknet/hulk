package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	libHttp "github.com/hulknet/hulk/app/http"
	"github.com/hulknet/hulk/app/types"
)

func main() {
	pKey, err := types.DecodeDefaultPublicKey()
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Use(libHttp.RegisterJWT(pKey))
	e.GET("/metrics", libHttp.PrometheusHandler())
	e.GET("/status", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, echo.Map{
			"status":  "OK",
			"service": "memory",
		})
	})
	fmt.Println(e.Start("127.0.0.1:7004"))
}
