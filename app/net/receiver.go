package net

import (
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"

	libHttp "github.com/hulknet/hulk/app/http"
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
	netMessage, err := libHttp.ParseHTTPHeader(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Errorf("header is invalid, %s\n", err)
		return
	}

	net, ok := rh.netCont.Net(netMessage.Time.BlockID())
	if !ok || !net.IsActive() {
		http.Error(w, "invalid block", http.StatusForbidden)
		return
	}

	if !net.AllowList().CheckToken(netMessage.Token) {
		http.Error(w, "invalid token", http.StatusForbidden)
		return
	}

	if !net.State().ValidateTime(netMessage.Time) {
		http.Error(w, "invalid time", http.StatusForbidden)
		return
	}

	netMessage.Data, err = ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	correct, err := netMessage.Sign[0].CheckSignature(netMessage.Encode())
	if err != nil || !correct {
		log.Error(err)
		http.Error(w, "invalid signature", http.StatusForbidden)
		return
	}

	nextPeer := net.Table().GetPeer(netMessage.Addr)

	if net.State().Peer().Equal(nextPeer) {
		net.HandleMessage(netMessage)
	} else {
		net.ProxyMessage(netMessage)
	}

	w.WriteHeader(http.StatusOK)
}
