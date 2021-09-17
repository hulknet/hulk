package ledger

import (
	"math/bits"

	"github.com/kotfalya/hulk/pkg/utils"
	"github.com/kotfalya/hulk/research/cpu/types"
)

type Block struct {
	ID      types.ID
	PID     types.ID
	PPID    types.ID
	BitSize uint8
	N       uint64
	U       uint64
}

func (b Block) IsPivot(pk types.PK) bool {
	if !b.IsRoot() {
		return false
	}
	l := bits.Len8(uint8(b.N))
	return utils.Cpl(pk[:1], b.ID[:1]) >= l
}

func (b Block) IsRoot() bool {
	return b.PPID.IsEmpty()
}

func (b Block) IsFirst() bool {
	if b.IsRoot() {
		return b.ID == b.PID
	} else {
		return b.PPID == b.PID
	}
}

type Tick []Block

func (t Tick) NodeBlock() Block {
	return t[len(t)-1]
}

func (t Tick) NetBlock() Block {
	return t[0]
}
