package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/hulknet/hulk/app/instruction"
	"github.com/hulknet/hulk/app/types"
)

func main() {
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
		ID: tickId.ID64(),
		//ID256:   tickId,
		Inc:    0,
		Level:  0,
		Status: types.TickStatusHead,
	}
	p := types.Peer{
		Pub:   ecpk.Pub(),
		Token: token,
	}

	var timeSrc []byte
	timeSrc = append(timeSrc, p.Pub.ID().Bytes()...)
	timeSrc = append(timeSrc, block.ID.Bytes()...)
	timeSrc = append(timeSrc, t.ID.Bytes()...)
	timeSrc = append(timeSrc, byte(1))

	time, err := types.DecodeTime(timeSrc)
	if err != nil {
		panic(err)
	}

	var blocks []types.Block
	blocks = append(blocks, block)

	var ticks []types.Tick
	ticks = append(ticks, t)

	s := types.CreateState(time, block, ticks, p, ecpk, token)

	hCont := instruction.NewHandlerContainer()
	hCont.SetState(s)

	log.Error(http.ListenAndServe("127.0.0.1:7004", instruction.NewReceiver(hCont)))
}
