package cpu

import (
	"github.com/kotfalya/hulk/research/cpu/ledger"
	"github.com/kotfalya/hulk/research/cpu/types"
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
