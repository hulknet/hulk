package app

import (
	"github.com/kotfalya/hulk/app/ledger"
	"github.com/kotfalya/hulk/app/types"
)

type Cpu struct {
	addr  types.Addr
	pks   []types.PK
	block ledger.Block
}

type Register struct {
	addr types.Addr
	cpu  types.Addr
	pks  []types.PK
}
