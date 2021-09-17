package http

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/kotfalya/hulk/research/cpu/types"
)

const (
	idHeader        = "ID"
	tokenHeader     = "Token"
	addrHeader      = "Addr"
	partHeader      = "Partition"
	signatureHeader = "Signature"
)

func ParseHTTPHeader(header http.Header) (messageHeader types.MessageHeader, err error) {
	if messageHeader.Token, err = parseToken(header); err != nil {
		return
	}
	if messageHeader.Sign, err = parseSignature(header); err != nil {
		return
	}
	if messageHeader.To, err = parseAddr(header); err != nil {
		return
	}
	if messageHeader.ID, err = parseID(header); err != nil {
		return
	}
	if messageHeader.Part, err = parsePart(header); err != nil {
		return
	}

	return
}

func parseToken(header http.Header) (token types.Token, err error) {
	tokenHex := header.Get(tokenHeader)
	if tokenHex == "" {
		err = errors.New(types.ErrGetToken)
		return
	}

	token, err = types.FromHex(tokenHex)
	if err != nil {
		err = errors.New(types.ErrDecodeToken)
		return
	}

	return
}

func parseSignature(header http.Header) (sign []types.Sign, err error) {
	signHex := header.Get(signatureHeader)
	if signHex == "" {
		err = errors.New(types.ErrGetSign)
		return
	}

	for _, s := range strings.Split(signHex, ",") {
		signItem, er := types.SignFromHex(s)
		if er != nil {
			err = errors.New(types.ErrDecodeSign)
			return
		}
		sign = append(sign, signItem)
	}

	return
}

func parseAddr(header http.Header) (addr types.Addr, err error) {
	addrStr := header.Get(addrHeader)
	if addrStr == "" {
		err = errors.New(types.ErrGetAddr)
		return
	}

	addrInt, err := strconv.ParseUint(addrStr, 10, 64)
	if err != nil {
		err = errors.New(types.ErrDecodeAddr)
		return
	}
	addr = types.Addr(addrInt)

	return
}

func parseID(header http.Header) (id types.ID, err error) {
	idStr := header.Get(idHeader)
	if idStr == "" {
		err = errors.New(types.ErrGetID)
		return
	}
	id, err = types.FromHex(idStr)

	return
}

func parsePart(header http.Header) (part types.Partition, err error) {
	partStr := header.Get(partHeader)
	if partStr == "" {
		return
	}
	partList := strings.Split(partStr, ",")
	if len(partList) != 2 {
		err = errors.New(types.ErrDecodePart)
		return
	}
	part.Position, err = strconv.ParseUint(partList[0], 10, 64)
	if err != nil {
		err = errors.New(types.ErrDecodePart)
		return
	}
	part.Length, err = strconv.ParseUint(partList[1], 10, 64)
	if err != nil {
		err = errors.New(types.ErrDecodePart)
		return
	}

	return
}
