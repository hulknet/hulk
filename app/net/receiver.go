package net

import (
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"

	libHttp "github.com/kotfalya/hulk/app/http"
	"github.com/kotfalya/hulk/app/types"
)

type ReceiverHandler struct {
	net *Net
}

func NewReceiverHandler(net *Net) *ReceiverHandler {
	return &ReceiverHandler{
		net: net,
	}
}

func (rh *ReceiverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	messageHeader, err := libHttp.ParseHTTPHeader(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Errorf("header is invalid, %s\n", err)
		return
	}

	if !rh.net.CheckToken(messageHeader.Token) {
		http.Error(w, "unknown token", http.StatusForbidden)
		return
	}

	messageBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	correct, err := types.CheckSignature(messageBody, messageHeader.Sign[0][:])
	if err != nil || !correct {
		log.Error(err)
		http.Error(w, "signature is invalid", http.StatusForbidden)
		return
	}

	if err = rh.net.HandleMessage(messageHeader, messageBody); err != nil {
		log.Error(err)
		http.Error(w, "failed to handle message", http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
