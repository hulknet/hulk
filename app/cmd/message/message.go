package main

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/crypto/sha3"

	"github.com/hulknet/hulk/app/types"
)

type MessageHeaderModel struct {
	ID    string
	Block string
	To    string
	From  string
	Time  string
	Token string
	Part  string
	Body  string
}

type SignModel struct {
	Body msgpack.RawMessage
}

type EncodeModel struct {
	Type string
	Sign string
	Data string
}

type EncodeResponseModel struct {
	Type string
	Sign []byte
	Data msgpack.RawMessage
}

type SplitModel struct {
	Message string
	Parts   int
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

		msgHash := sha3.Sum256(m.Body)
		sign, err := crypto.Sign(msgHash[:], ecpk.ECPrivateKey())
		if err != nil {
			log.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return ctx.JSON(http.StatusOK, echo.Map{
			"sing": hex.EncodeToString(sign),
			"body": hex.EncodeToString(m.Body),
		})
	})

	e.POST("/encode", func(ctx echo.Context) error {
		m := new(EncodeModel)
		if err := ctx.Bind(m); err != nil {
			log.Error(err)
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		data, err := types.FromHex(m.Data, 0)
		if err != nil {
			return err
		}

		sign, err := types.FromHex(m.Sign, 0)
		if err != nil {
			return err
		}
		mr := &EncodeResponseModel{
			Type: m.Type,
			Data: data,
			Sign: sign,
		}

		mrByte, err := msgpack.Marshal(mr)

		return ctx.JSON(http.StatusOK, echo.Map{
			"message": hex.EncodeToString(mrByte),
		})
	})

	e.POST("/split", func(ctx echo.Context) error {
		s := new(SplitModel)
		if err := ctx.Bind(s); err != nil {
			return err
		}

		message, err := types.FromHex(s.Message, 0)
		if err != nil {
			return err
		}

		chunk, err := types.EncryptToParts(message, s.Parts)
		if err != nil {
			return err
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
			return err
		}

		sign, err := signMessage(m, ecpk.ECPrivateKey())
		if err != nil {
			return err
		}

		client := http.Client{}

		body, err := types.FromHex(m.Body, 0)
		if err != nil {
			return err
		}

		req, err := http.NewRequest("POST", "http://127.0.0.1:7002", bytes.NewBuffer(body))
		if err != nil {
			return err
		}

		req.Header.Add("ID", m.ID)
		req.Header.Add("Block", m.Block)
		req.Header.Add("Token", m.Token)
		req.Header.Add("To", m.To)
		req.Header.Add("Time", m.Time)
		req.Header.Add("From", m.From)
		req.Header.Add("Partition", m.Part)
		req.Header.Add("Signature", hex.EncodeToString(sign))
		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		err = resp.Body.Close()
		if err != nil {
			return err
		}

		return ctx.JSON(http.StatusOK, echo.Map{
			"status":        resp.Status,
			"resp":          string(respBody),
			"bodySignature": hex.EncodeToString(sign),
		})
	})

	fmt.Println(e.Start("127.0.0.1:7009"))
}

func signMessage(m *MessageHeaderModel, pKey *ecdsa.PrivateKey) ([]byte, error) {
	id, err := parseID(m.ID)
	if err != nil {
		return nil, err
	}
	blockId, err := parseID(m.Block)
	if err != nil {
		return nil, err
	}
	to, err := parseID(m.To)
	if err != nil {
		return nil, err
	}
	from, err := parseID(m.From)
	if err != nil {
		return nil, err
	}
	part, err := parsePart(m.Part)
	if err != nil {
		return nil, err
	}
	time, err := types.FromHex(m.Time, 0)
	if err != nil {
		return nil, err
	}

	body, err := types.FromHex(m.Body, 0)
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
