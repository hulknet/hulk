package main

import (
	"fmt"
	"net/http"

	"github.com/kotfalya/hulk/research/cpu/rest"
	"github.com/kotfalya/hulk/research/cpu/types"
	"github.com/labstack/echo/v4"
)

func main() {
	pKey, err := types.DecodeDefaultPublicKey()
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Use(rest.RegisterJWT(pKey))
	e.GET("/status", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, echo.Map{
			"status":  "OK",
			"service": "store",
		})
	})
	fmt.Println(e.Start("127.0.0.1:7004"))
}
