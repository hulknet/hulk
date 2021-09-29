package net

import (
	"bytes"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"

	libHttp "github.com/kotfalya/hulk/app/http"
)

type ReceiverHandler struct {
	netCont *Container
}

func NewReceiverHandler(netCont *Container) *ReceiverHandler {
	return &ReceiverHandler{
		netCont: netCont,
	}
}

func (rh *ReceiverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	messageHeader, err := libHttp.ParseHTTPHeader(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Errorf("header is invalid, %s\n", err)
		return
	}

	net, ok := rh.netCont.Net(messageHeader.BlockID)
	if !ok || !net.IsActive() {
		http.Error(w, "invalid block", http.StatusForbidden)
		return
	}

	if !net.AllowList().CheckToken(messageHeader.Token) {
		http.Error(w, "invalid token", http.StatusForbidden)
		return
	}

	if !net.State().ValidateTime(messageHeader.From, messageHeader.Time) {
		http.Error(w, "invalid time", http.StatusForbidden)
		return
	}

	messageBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	data := bytes.NewBuffer(messageHeader.ID.Bytes())
	data.Write(messageHeader.BlockID.Bytes())
	data.Write(messageHeader.To.Bytes())
	data.Write(messageHeader.From.Bytes())
	data.Write(messageHeader.Time.Bytes())
	data.WriteByte(messageHeader.Part.Position)
	data.WriteByte(messageHeader.Part.Length)
	data.Write(messageBody)

	correct, err := messageHeader.Sign[0].CheckSignature(data.Bytes())
	if err != nil || !correct {
		log.Error(err)
		http.Error(w, "invalid signature", http.StatusForbidden)
		return
	}

	nextPeer := net.Table().GetPeer(messageHeader.To)

	if net.State().Peer().Equal(nextPeer) {
		net.HandleMessage(messageHeader, messageBody)
	} else {
		net.ProxyMessage(messageHeader, messageBody)
	}

	w.WriteHeader(http.StatusOK)
}
