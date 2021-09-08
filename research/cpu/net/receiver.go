package net

import (
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	libHttp "github.com/kotfalya/hulk/research/cpu/http"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/sha3"
	"io/ioutil"
	"net/http"
)

type ReceiverHandler struct {
	net *Net
}

func NewReceiverHandler(net *Net) *ReceiverHandler {
	return &ReceiverHandler{
		net: net,
	}
}

func checkSignature(msg []byte, sign []byte) (bool, error) {
	msgHash := sha3.Sum256(msg)
	pk, err := secp256k1.RecoverPubkey(msgHash[:], sign[:])
	if err != nil {
		return false, err
	}

	return secp256k1.VerifySignature(pk, msgHash[:], sign[:64]), nil
}

func (rh *ReceiverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	messageHeader, err := libHttp.ParseHTTPHeader(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Errorf("header is invalid, %s\n", err)
		return
	}

	peerIn, ok := rh.net.CheckToken(messageHeader.Token)
	if !ok {
		http.Error(w, "token is invalid", http.StatusForbidden)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	correct, err := checkSignature(body, messageHeader.Sign[0][:])
	if err != nil || !correct {
		log.Error(err)
		http.Error(w, "signature is invalid", http.StatusForbidden)
		return
	}

	if !rh.net.CheckPeer(peerIn) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	peer := rh.net.FindPeer(messageHeader.To)
	if peer.Equal(rh.net.Self()) {
		w.Write([]byte("self"))
	} else {
		w.Write([]byte("foreign"))
	}

	w.WriteHeader(http.StatusOK)
}
