package http

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/kotfalya/hulk/app/types"
)

const (
	idHeader        = "ID"
	tokenHeader     = "Token"
	addrHeader      = "Addr"
	toHeader        = "To"
	fromHeader      = "From"
	partHeader      = "Partition"
	signatureHeader = "Signature"
)

func ParseHTTPHeader(header http.Header) (messageHeader types.MessageHeader, err error) {
	if messageHeader.Token, err = parseToken(header.Get(tokenHeader)); err != nil {
		return
	}
	if messageHeader.Sign, err = parseSignature(header.Get(signatureHeader)); err != nil {
		return
	}
	if messageHeader.To, err = parseIDPrefix(header.Get(toHeader)); err != nil {
		return
	}
	if messageHeader.From, err = parseIDPrefix(header.Get(fromHeader)); err != nil {
		return
	}
	if messageHeader.ID, err = parseID(header.Get(idHeader)); err != nil {
		return
	}
	if messageHeader.Part, err = parsePart(header.Get(partHeader)); err != nil {
		return
	}

	return
}

func parseToken(tokenStr string) (token types.Token, err error) {
	if tokenStr == "" {
		err = errors.New(types.ErrGetToken)
		return
	}

	token, err = types.IDFromHex(tokenStr)
	if err != nil {
		err = errors.New(types.ErrDecodeToken)
		return
	}

	return
}

func parseSignature(signStr string) (sign []types.Sign, err error) {
	if signStr == "" {
		err = errors.New(types.ErrGetSign)
		return
	}

	for _, s := range strings.Split(signStr, ",") {
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

func parseID(idStr string) (id types.ID, err error) {
	if idStr == "" {
		err = errors.New(types.ErrGetID)
		return
	}
	id, err = types.IDFromHex(idStr)

	return
}

func parseIDPrefix(srcStr string) (id types.IDPrefix, err error) {
	if srcStr == "" {
		err = errors.New(types.ErrGetIDPrefix)
		return
	}
	id, err = types.IDPrefixFromHex(srcStr)

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
