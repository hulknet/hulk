package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/kotfalya/hulk/app/types"
)

const (
	idHeader        = "ID"
	tokenHeader     = "Token"
	blockHeader     = "Block"
	timeHeader      = "Time"
	toHeader        = "To"
	fromHeader      = "From"
	partHeader      = "Partition"
	signatureHeader = "Signature"
)

func ParseHTTPHeader(header http.Header) (messageHeader types.MessageHeader, err error) {
	if messageHeader.Token, err = parseToken(header.Get(tokenHeader)); err != nil {
		return
	}
	if messageHeader.BlockID, err = parseID(header.Get(blockHeader)); err != nil {
		return
	}
	if messageHeader.Sign, err = parseSignature(header.Get(signatureHeader)); err != nil {
		return
	}
	if messageHeader.Time, err = parseTime(header.Get(timeHeader)); err != nil {
		return
	}
	if messageHeader.To, err = parseID(header.Get(toHeader)); err != nil {
		return
	}
	if messageHeader.From, err = parseID(header.Get(fromHeader)); err != nil {
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

	token, err = types.TokenFromHex(tokenStr)
	if err != nil {
		err = errors.New(types.ErrDecodeToken)
		return
	}

	return
}

func parseSignature(signStr string) (sign []types.Sign520, err error) {
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

func parseTime(timeStr string) (time types.Time, err error) {
	if timeStr == "" {
		err = errors.New(types.ErrGetTime)
		return
	}
	time, parserErr := types.FromHex(timeStr, 0)
	if parserErr != nil {
		err = errors.New(types.ErrDecodeTime)
		return
	}
	if !time.Validate() {
		err = errors.New(types.ErrInvalidTime)
		return
	}
	return
}

func parseID(srcStr string) (id types.ID64, err error) {
	if srcStr == "" {
		err = errors.New(types.ErrGetID)
		return
	}
	id, err = types.ID64FromHex(srcStr)

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
