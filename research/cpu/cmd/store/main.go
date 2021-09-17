package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	libHttp "github.com/kotfalya/hulk/research/cpu/http"
	"github.com/kotfalya/hulk/research/cpu/types"
)

func main() {
	pKey, err := types.DecodeDefaultPublicKey()
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Use(libHttp.RegisterJWT(pKey))
	e.GET("/status", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, echo.Map{
			"status":  "OK",
			"service": "store",
		})
	})
	fmt.Println(e.Start("127.0.0.1:7004"))
}
