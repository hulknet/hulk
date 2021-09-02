package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/kotfalya/hulk/research/cpu/types"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/sha3"
	"io/ioutil"
	"net/http"
)

type MessageHeaderModel struct {
	ID    string           `json:"id"`
	Addr  string           `json:"addr"`
	Token string           `json:"token"`
	Body  *json.RawMessage `json:"body"`
	//Part  Partition `json:"part"`
}

type SignModel struct {
	Body *json.RawMessage `json:"body"`
}

type SendModel struct {
	Body string `json:"body"`
	Addr string `json:"addr"`
}

type TokenModel struct {
	Token string `json:"token"`
}

func main() {
	pKey, err := types.DecodeDefaultPrivateKey()
	if err != nil {
		panic(err)
	}

	var token string

	e := echo.New()

	e.GET("/status", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, echo.Map{
			"status": "OK",
		})
	})

	e.POST("/sign", func(ctx echo.Context) error {
		m := new(SignModel)
		if err := ctx.Bind(m); err != nil {
			log.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		msgHash := sha3.Sum256(*m.Body)
		sign, err := crypto.Sign(msgHash[:], pKey)
		if err != nil {
			log.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return ctx.JSON(http.StatusOK, echo.Map{
			"sing": hex.EncodeToString(sign),
		})
	})

	e.POST("/token", func(ctx echo.Context) error {
		m := new(TokenModel)
		if err := ctx.Bind(m); err != nil {
			log.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		token = m.Token

		return ctx.JSON(http.StatusOK, echo.Map{
			"status": "OK",
			"token":  token,
		})
	})

	e.POST("/send", func(ctx echo.Context) error {
		m := new(SendModel)
		if err := ctx.Bind(m); err != nil {
			log.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		msgHash := sha3.Sum256([]byte(m.Body))
		sign, err := crypto.Sign(msgHash[:], pKey)
		if err != nil {
			log.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		client := http.Client{}

		req, err := http.NewRequest("POST", "http://127.0.0.1:7002", bytes.NewBufferString(m.Body))
		if err != nil {
			log.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		req.Header.Add("Token", token)
		req.Header.Add("Addr", m.Addr)
		req.Header.Add("Signature", hex.EncodeToString(sign))
		resp, err := client.Do(req)
		if err != nil {
			log.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		err = resp.Body.Close()
		if err != nil {
			log.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return ctx.JSON(http.StatusOK, echo.Map{
			"status": resp.Status,
		})
	})

	e.POST("/test", func(ctx echo.Context) error {
		m := new(MessageHeaderModel)
		if err := ctx.Bind(m); err != nil {
			log.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		msgHash := sha3.Sum256(*m.Body)
		sign, err := crypto.Sign(msgHash[:], pKey)
		if err != nil {
			log.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		client := http.Client{}

		req, err := http.NewRequest("POST", "http://127.0.0.1:7002", bytes.NewBuffer(*m.Body))
		if err != nil {
			log.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		req.Header.Add("Token", m.Token)
		req.Header.Add("Addr", m.Addr)
		req.Header.Add("Signature", hex.EncodeToString(sign))
		resp, err := client.Do(req)
		if err != nil {
			log.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		err = resp.Body.Close()
		if err != nil {
			log.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return ctx.JSON(http.StatusOK, echo.Map{
			"status":        resp.Status,
			"resp":          string(respBody),
			"bodySignature": sign,
		})
	})

	fmt.Println(e.Start("127.0.0.1:7009"))
}
