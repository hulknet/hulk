package main

import (
	"fmt"
	"net/http"

	"github.com/hulknet/hulk/app/net"
	"github.com/hulknet/hulk/app/types"
)

func main() {
	//  tmp ----------------------------------
	//var pk types.Pub = types.GenerateSHA()
	//var token types.Token = types.GenerateSHA()
	ecpk, err := types.HexToECKey("90313109591dea4b6e4f4145c7f0124ebf05079b43327d06201ae746a2282ef3")
	if err != nil {
		panic(err)
	}
	token, err := types.TokenFromHex("051eaf028faeed1e4a1c7acc68c4e1ad2ec49b283632a65dd632010833f6164e")
	if err != nil {
		panic(err)
	}
	blockId, err := types.ID256FromHex("f38157e98d676c4299899118a4a6ecae16f6f1c19013007b35dc7c23f2d52e7a")
	if err != nil {
		panic(err)
	}
	tickId, err := types.ID256FromHex("ecae16f6f1c19013007b35dc7c23f2d52e7af38157e98d676c4299899118a4a6")
	if err != nil {
		panic(err)
	}

	block := types.Block{
		ID:      blockId.ID64(),
		ID256:   blockId,
		BitSize: []byte{1},
		Status:  types.BlockStatusHead,
	}
	t := types.Tick{
		IsLocal: false,
		ID:      tickId.ID64(),
		ID256:   tickId,
		Count:   0,
		Level:   0,
		Status:  types.TickStatusHead,
	}
	p := types.Peer{
		Pub:   ecpk.Pub(),
		Token: token,
	}

	var time types.Time
	time = append(time, t.ID.Bytes()...)
	time = append(time, byte(1))

	var blocks []types.Block
	blocks = append(blocks, block)

	var ticks []types.Tick
	ticks = append(ticks, t)

	s := types.CreateState(time, block, ticks, p, ecpk, token)

	netCont := net.NewNetContainer()
	netCont.SetState(s)

	//pk for jwt internal communication
	pKey, err := types.DecodeDefaultPublicKey()
	if err != nil {
		panic(err)
	}

	r := net.NewRestServer(netCont, "127.0.0.1:7001", pKey)

	errChan := make(chan error)
	go func() {
		errChan <- r.Listen()
	}()
	go func() {
		errChan <- http.ListenAndServe("127.0.0.1:7002", net.NewReceiverHandler(netCont))
	}()

	select {
	case err := <-errChan:
		fmt.Println(err)
	}
}
