package rest

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/kotfalya/hulk/research/cpu/types"
)

const (
	tokenHeader     = "Token"
	addrHeader      = "Addr"
	signatureHeader = "Signature"
)

func HeaderToToken(header http.Header) (types.Token, error) {
	tokenHex := header.Get(tokenHeader)
	if tokenHex == "" {
		return types.Token{}, errors.New(types.ErrGetToken)
	}

	token, err := types.FromHex(tokenHex)
	if err != nil {
		return types.Token{}, errors.New(types.ErrDecodeToken)
	}

	return token, nil
}

func HeaderToSignature(header http.Header) (types.Sign, error) {
	signHex := header.Get(signatureHeader)
	if signHex == "" {
		return types.Sign{}, errors.New(types.ErrGetToken)
	}

	sign, err := types.SignFromHex(signHex)
	if err != nil {
		return sign, errors.New(types.ErrDecodeToken)
	}

	return sign, nil
}

func HeaderToAddr(header http.Header) (types.Addr, error) {
	addrStr := header.Get(addrHeader)
	addr, err := strconv.ParseUint(addrStr, 10, 64)
	if err != nil {
		return 0, errors.New(types.ErrDecodeAddr)
	}
	return types.Addr(addr), nil
}
