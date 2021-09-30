package main

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/sha3"

	"github.com/kotfalya/hulk/app/types"
)

type MessageHeaderModel struct {
	ID      string `json:"id"`
	BlockID string `json:"block"`
	To      string `json:"to"`
	From    string `json:"from"`
	Time    string `json:"time"`
	Token   string `json:"token"`
	Part    string `json:"part"`
	Body    string `json:"body"`
}

type SignModel struct {
	ID   string           `json:"id"`
	Part string           `json:"part"`
	To   string           `json:"to"`
	From string           `json:"from"`
	Time string           `json:"time"`
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

func main() {
	ecpk, err := types.HexToECKey("90313109591dea4b6e4f4145c7f0124ebf05079b43327d06201ae746a2282ef3")
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
		sign, err := crypto.Sign(msgHash[:], ecpk.ECPrivateKey())
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

		chunk, err := types.EncryptToParts(*s.Message, s.Parts)
		if err != nil {
			log.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		_, err = types.DecryptFromParts(chunk)
		if err != nil {
			log.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		res := make([]string, len(chunk))
		for i, v := range chunk {
			res[i] = hex.EncodeToString(v)
		}
		return ctx.JSON(http.StatusOK, echo.Map{
			"status":    "OK",
			"chunks":    res,
			"messageId": types.GenerateSHA().ID64().Hex(),
		})
	})

	e.POST("/send", func(ctx echo.Context) error {
		m := new(MessageHeaderModel)
		if err := ctx.Bind(m); err != nil {
			log.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		sign, err := signMessage(m.Body, m.BlockID, m.ID, m.To, m.From, m.Time, m.Part, ecpk.ECPrivateKey())
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
		req.Header.Add("Block", m.BlockID)
		req.Header.Add("Token", m.Token)
		req.Header.Add("To", m.To)
		req.Header.Add("Time", m.Time)
		req.Header.Add("From", m.From)
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
			"bodySignature": hex.EncodeToString(sign),
		})
	})

	fmt.Println(e.Start("127.0.0.1:7009"))
}

func signMessage(bodyStr, idStr, blockIdStr, toStr, fromStr, timeStr, partStr string, pKey *ecdsa.PrivateKey) ([]byte, error) {
	id, err := parseID(idStr)
	if err != nil {
		return nil, err
	}
	blockId, err := parseID(blockIdStr)
	if err != nil {
		return nil, err
	}
	to, err := parseID(toStr)
	if err != nil {
		return nil, err
	}
	from, err := parseID(fromStr)
	if err != nil {
		return nil, err
	}
	part, err := parsePart(partStr)
	if err != nil {
		return nil, err
	}
	time, err := types.FromHex(timeStr, 0)
	if err != nil {
		return nil, err
	}

	body, err := types.FromHex(bodyStr, 0)
	if err != nil {
		return nil, err
	}

	data := bytes.NewBuffer(id[:])
	data.Write(blockId[:])
	data.Write(to[:])
	data.Write(from[:])
	data.Write(time[:])
	data.WriteByte(part.Position)
	data.WriteByte(part.Length)
	data.Write(body)

	msgHash := sha3.Sum256(data.Bytes())
	return crypto.Sign(msgHash[:], pKey)
}

func parseID(idStr string) (id types.ID64, err error) {
	id, err = types.ID64FromHex(idStr)

	return
}

func parsePart(partStr string) (part types.Partition, err error) {
	if partStr == "" {
		return
	}
	partBytes, err := types.FromHex(partStr, 2)
	if err != nil {
		err = errors.New(types.ErrDecodePart)
		return
	}

	part.Position = partBytes[0]
	part.Position = partBytes[1]

	return
}
