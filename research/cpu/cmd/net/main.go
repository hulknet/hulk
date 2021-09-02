package main

import (
	"fmt"
	"net/http"

	"github.com/kotfalya/hulk/research/cpu/ledger"
	"github.com/kotfalya/hulk/research/cpu/net"
	"github.com/kotfalya/hulk/research/cpu/types"
)

func main() {
	//  tmp ----------------------------------
	var pk types.PK = types.GenerateSHA()
	var token types.Token = types.GenerateSHA()

	id := pk.ID()

	b := ledger.Block{
		ID:      id,
		PID:     id,
		BitSize: 1,
		N:       2,
		U:       1,
	}
	t := ledger.Tick{b}
	pOut := types.PeerOut{
		PK:    pk,
		Token: token,
	}

	n := net.NewNet(pOut)
	n.SetTick(t)

	//var pk1 types.PK = types.GenerateSHA()
	//var token1 types.Token = types.GenerateSHA()
	//pOut1 := types.PeerOut{
	//	PK:    pk1,
	//	Token: token1,
	//}
	//n.AddPeer(pOut1)

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
