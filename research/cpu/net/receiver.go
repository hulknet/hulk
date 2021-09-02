package net

import (
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/kotfalya/hulk/research/cpu/rest"
	"golang.org/x/crypto/sha3"
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
	token, err := rest.HeaderToToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		log.Errorf("token is invalid, %s\n", err)
		return
	}

	peerIn, ok := rh.net.CheckToken(token)
	if !ok {
		http.Error(w, "token is invalid", http.StatusForbidden)
		return
	}

	sign, err := rest.HeaderToSignature(r.Header)
	if err != nil {
		http.Error(w, "signature is invalid", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	s := string(body)
	fmt.Println(s)
	correct, err := checkSignature(body, sign[:])
	if err != nil || !correct {
		log.Error(err)
		http.Error(w, "signature is invalid", http.StatusBadRequest)
		return
	}

	if !rh.net.CheckPeer(peerIn) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	addr, err := rest.HeaderToAddr(r.Header)
	if err != nil {
		log.Errorf("addr is invalid, %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	peer := rh.net.FindPeer(addr)
	if peer.Equal(rh.net.Self()) {
		w.Write([]byte("self"))
	} else {
		w.Write([]byte("foreign"))
	}

	w.WriteHeader(http.StatusOK)
}
