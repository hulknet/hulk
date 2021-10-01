package main

import (
	"crypto/subtle"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"

	libHttp "github.com/hulknet/hulk/app/http"
	"github.com/hulknet/hulk/app/types"
)

const DefaultPassword = "password"

type ServiceLogin struct {
	Type     string `json:"type"`
	Password string `json:"password"`
}

func main() {
	pKey, err := types.DecodeDefaultPrivateKey()
	if err != nil {
		panic(err)
	}

	e := echo.New()

	e.GET("/status", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, echo.Map{
			"status": "OK",
		})
	})

	e.POST("/login", func(ctx echo.Context) error {
		s := new(ServiceLogin)
		if err := ctx.Bind(s); err != nil {
			log.Error(err)
			return nil
		}
		if s.Password != DefaultPassword {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
		}

		tokenType := libHttp.TokenTypeFromString(s.Type)
		if tokenType == libHttp.UnknownTokenType {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid type of service")
		}

		t, err := libHttp.GenerateToken(pKey, tokenType)
		if err != nil {
			log.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return ctx.JSON(http.StatusOK, echo.Map{
			"token": t,
		})
	})

	g := e.Group("/admin")
	g.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// Be careful to use constant time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(username), []byte("admin")) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte("admin")) == 1 {
			return true, nil
		}
		return false, nil
	}))

	fmt.Println(e.Start("127.0.0.1:7003"))
}
