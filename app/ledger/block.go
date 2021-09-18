package ledger

import (
	"github.com/kotfalya/hulk/app/types"
)

type Block struct {
	ID      types.ID
	PID     types.ID
	PPID    types.ID
	BitSize uint8
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
