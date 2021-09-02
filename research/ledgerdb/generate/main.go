package main

import (
	"fmt"

	"github.com/asdine/storm/v3"
	gobCodec "github.com/asdine/storm/v3/codec/gob"
	"github.com/kotfalya/hulk/pkg/store"
	"github.com/kotfalya/hulk/pkg/utils"
)

func main() {
	db, err := storm.Open("../db", storm.Codec(gobCodec.Codec), storm.Batch())

	if err != nil {
		panic(err)
	}

	err = db.Drop("LedgerBlock")
	if err != nil {
		fmt.Println(err)
	}

	n := utils.GenerateSHA()
	for i := 0; i < 100000; i++ {
		block := store.LedgerBlock{
			ID:   n.WithSalt(n[:]),
			Prev: n,
			Tick: uint64(i),
		}
		n = utils.GenerateSHA()

		err := db.Save(&block)
		if err != nil {
			panic(err)
		}
		if i%500 == 0 {
			fmt.Println(i)
		}
	}

}
