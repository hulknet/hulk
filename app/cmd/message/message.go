package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/sha3"

	types2 "github.com/kotfalya/hulk/app/types"
)

type MessageHeaderModel struct {
	ID    string `json:"id"`
	Addr  string `json:"addr"`
	Token string `json:"token"`
	Part  string `json:"part"`
	Body  string `json:"body"`
}

type SignModel struct {
	Body *json.RawMessage `json:"body"`
}

type SplitModel struct {
	Message *json.RawMessage `json:"message"`
	Parts   int              `json:"parts"`
}

type SendModel struct {
	ID    string `json:"id"`
	Token string `json:"token"`
	Body  string `json:"body"`
	Addr  string `json:"addr"`
}

type TokenModel struct {
	Token string `json:"token"`
}

func main() {
	pKey, err := types2.DecodeDefaultPrivateKey()
	if err != nil {
		panic(err)
	}

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
			"body": *m.Body,
		})
	})

	e.POST("/split", func(ctx echo.Context) error {
		s := new(SplitModel)
		if err := ctx.Bind(s); err != nil {
			log.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		chunk, err := types2.EncryptToParts(*s.Message, s.Parts)
		if err != nil {
			log.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		_, err = types2.DecryptFromParts(chunk)
		if err != nil {
			log.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		res := make([]string, len(chunk))
		for i, v := range chunk {
			res[i] = hex.EncodeToString(v)
		}
		return ctx.JSON(http.StatusOK, echo.Map{
			"status": "OK",
			"chunks": res,
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

		req.Header.Add("ID", m.ID)
		req.Header.Add("Token", m.Token)
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

		req.Header.Add("ID", m.ID)
		req.Header.Add("Token", m.Token)
		req.Header.Add("Addr", m.Addr)
		req.Header.Add("Partition", m.Part)
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
