package http

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	types2 "github.com/kotfalya/hulk/app/types"
)

const (
	idHeader        = "ID"
	tokenHeader     = "Token"
	addrHeader      = "Addr"
	partHeader      = "Partition"
	signatureHeader = "Signature"
)

func ParseHTTPHeader(header http.Header) (messageHeader types2.MessageHeader, err error) {
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

func parseToken(header http.Header) (token types2.Token, err error) {
	tokenHex := header.Get(tokenHeader)
	if tokenHex == "" {
		err = errors.New(types2.ErrGetToken)
		return
	}

	token, err = types2.FromHex(tokenHex)
	if err != nil {
		err = errors.New(types2.ErrDecodeToken)
		return
	}

	return
}

func parseSignature(header http.Header) (sign []types2.Sign, err error) {
	signHex := header.Get(signatureHeader)
	if signHex == "" {
		err = errors.New(types2.ErrGetSign)
		return
	}

	for _, s := range strings.Split(signHex, ",") {
		signItem, er := types2.SignFromHex(s)
		if er != nil {
			err = errors.New(types2.ErrDecodeSign)
			return
		}
		sign = append(sign, signItem)
	}

	return
}

func parseAddr(header http.Header) (addr types2.Addr, err error) {
	addrStr := header.Get(addrHeader)
	if addrStr == "" {
		err = errors.New(types2.ErrGetAddr)
		return
	}

	addrInt, err := strconv.ParseUint(addrStr, 10, 64)
	if err != nil {
		err = errors.New(types2.ErrDecodeAddr)
		return
	}
	addr = types2.Addr(addrInt)

	return
}

func parseID(header http.Header) (id types2.ID, err error) {
	idStr := header.Get(idHeader)
	if idStr == "" {
		err = errors.New(types2.ErrGetID)
		return
	}
	id, err = types2.FromHex(idStr)

	return
}

func parsePart(header http.Header) (part types2.Partition, err error) {
	partStr := header.Get(partHeader)
	if partStr == "" {
		return
	}
	partList := strings.Split(partStr, ",")
	if len(partList) != 2 {
		err = errors.New(types2.ErrDecodePart)
		return
	}
	part.Position, err = strconv.ParseUint(partList[0], 10, 64)
	if err != nil {
		err = errors.New(types2.ErrDecodePart)
		return
	}
	part.Length, err = strconv.ParseUint(partList[1], 10, 64)
	if err != nil {
		err = errors.New(types2.ErrDecodePart)
		return
	}

	return
}
