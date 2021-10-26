package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/hulknet/hulk/app/types"
)

const (
	idHeader        = "ID"
	tokenHeader     = "Token"
	timeHeader      = "Time"
	addrHeader      = "Addr"
	signatureHeader = "Signature"
)

func ParseHTTPHeader(header http.Header) (netMessage types.NetMessage, err error) {
	if netMessage.Token, err = ParseToken(header.Get(tokenHeader)); err != nil {
		return
	}
	if netMessage.Sign, err = ParseSignature(header.Get(signatureHeader)); err != nil {
		return
	}
	if netMessage.Time, err = ParseTime(header.Get(timeHeader)); err != nil {
		return
	}
	if netMessage.Addr, err = ParseID(header.Get(addrHeader)); err != nil {
		return
	}
	if netMessage.ID, err = ParseID(header.Get(idHeader)); err != nil {
		return
	}

	return
}

func ParseToken(tokenStr string) (token types.Token, err error) {
	if tokenStr == "" {
		err = errors.New(types.ErrGetToken)
		return
	}

	token, err = types.TokenFromHex(tokenStr)
	if err != nil {
		err = errors.New(types.ErrDecodeToken)
		return
	}

	return
}

func ParseSignature(signStr string) (sign []types.Sign520, err error) {
	if signStr == "" {
		err = errors.New(types.ErrGetSign)
		return
	}

	for _, s := range strings.Split(signStr, ",") {
		signItem, er := types.SignFromHex(s)
		if er != nil {
			err = errors.New(types.ErrDecodeSign520)
			return
		}
		sign = append(sign, signItem)
	}

	return
}

func ParseTime(timeStr string) (time types.Time, err error) {
	if timeStr == "" {
		err = errors.New(types.ErrGetTime)
		return
	}
	timeSrc, parserErr := types.FromHex(timeStr, 0)
	if parserErr != nil {
		err = errors.New(types.ErrDecodeTime)
		return
	}
	time, err = types.DecodeTime(timeSrc)
	return
}

func ParseID(srcStr string) (id types.ID64, err error) {
	if srcStr == "" {
		err = errors.New(types.ErrGetID)
		return
	}
	id, err = types.ID64FromHex(srcStr)

	return
}
