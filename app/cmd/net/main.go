package main

import (
	"fmt"
	"net/http"

	"github.com/kotfalya/hulk/app/ledger"
	"github.com/kotfalya/hulk/app/net"
	"github.com/kotfalya/hulk/app/types"
)

func main() {
	//  tmp ----------------------------------
	//var pk types.PK = types.GenerateSHA()
	//var token types.Token = types.GenerateSHA()
	pk, err := types.FromHex("f38157e98d676c4299899118a4a6ecae16f6f1c19013007b35dc7c23f2d52e7a")
	if err != nil {
		panic(err)
	}
	token, err := types.FromHex("051eaf028faeed1e4a1c7acc68c4e1ad2ec49b283632a65dd632010833f6164e")
	if err != nil {
		panic(err)
	}

	id := types.PK(pk).ID()

	b := ledger.Block{
		ID:      id,
		PID:     id,
		BitSize: 1,
	}
	t := ledger.Tick{b}
	pOut := types.Peer{
		PK:    pk,
		Token: token,
	}

	n := net.NewNet(pOut)
	n.Init(t)
	go func() {
		panic(n.Start())
	}()

	pKey, err := types.DecodeDefaultPublicKey()
	if err != nil {
		panic(err)
	}

	r := net.NewRestServer(n, "127.0.0.1:7001", pKey)

	errChan := make(chan error)
	go func() {
		errChan <- r.Listen()
	}()
	go func() {
		errChan <- http.ListenAndServe("127.0.0.1:7002", net.NewReceiverHandler(n))
	}()

	select {
	case err := <-errChan:
		fmt.Println(err)
	}
}
