package main

import (
	"fmt"
	"time"

	"github.com/asdine/storm/v3"
	gobCodec "github.com/asdine/storm/v3/codec/gob"
	"github.com/kotfalya/hulk/pkg/store"
)

func main() {
	db, err := storm.Open("../db", storm.Codec(gobCodec.Codec))
	if err != nil {
		panic(err)
	}

	var blocks []store.LedgerBlock
	startTime := time.Now()
	//err = db.Range("Tick", 44000, 45000, &blocks, storm.Limit(2), storm.Reverse())
	err = db.AllByIndex("Tick", &blocks, storm.Limit(2), storm.Reverse())
	if err != nil {
		panic(err)
	}

	fmt.Println(time.Since(startTime).Microseconds())
	for _, block := range blocks {
		fmt.Println(block.Tick)
	}
}
