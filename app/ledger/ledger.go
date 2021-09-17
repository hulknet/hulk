package ledger

import (
	"fmt"

	"github.com/kotfalya/hulk/app/types"
)

type Ledger struct {
	b   map[types.ID]Block
	cb  types.ID
	frb types.ID
}

func NewLedger() *Ledger {
	return new(Ledger)
}

func (l *Ledger) Load(bl []Block, cb types.ID) error {
	for _, b := range bl {
		if _, ok := l.b[b.PID]; len(l.b) > 0 && !ok {
			return fmt.Errorf("block %v has invalid PrevID", b.PID)
		}
		l.b[b.ID] = b
		if b.PID == cb {
			l.frb = b.ID
		}
	}
	l.cb = cb
	return nil
}

// AddBlock TODO: implement nested blocks
func (l *Ledger) AddBlock(b Block) error {
	if b.IsRoot() {
		return l.addRootBlock(b)
	} else {
		return l.addBlock(b)
	}
}

func (l *Ledger) addBlock(b Block) error {
	if _, ok := l.b[b.PID]; !ok {
		return fmt.Errorf("block %v has invalid PrevID", b.PID)
	}
	if _, ok := l.b[b.PPID]; !ok {
		return fmt.Errorf("block %v has invalid ParentID", b.PID)
	}
	l.b[b.ID] = b
	return nil
}

func (l *Ledger) addRootBlock(b Block) error {
	if _, ok := l.b[b.PID]; !ok {
		return fmt.Errorf("block %v has invalid PrevID", b.PID)
	}
	l.b[b.ID] = b
	l.frb = b.ID
	return nil
}

func (l *Ledger) Block() Block {
	return l.b[l.cb]
}

func (l *Ledger) FutureRootBlock() Block {
	return l.b[l.frb]
}

func (l *Ledger) NextRoot() {
	l.cb = l.frb
}

func (l *Ledger) Prev(b Block) Block {
	return l.b[b.PID]
}

func (l *Ledger) Parent(b Block) Block {
	return l.b[b.PPID]
}

func (l *Ledger) Tick() Tick {
	blocks := make([]Block, 0)
	bID := l.cb
	for {
		b := l.b[bID]
		blocks = append(blocks, b)
		if b.IsRoot() {
			break
		}
		bID = b.PPID
	}

	t := Tick{}
	for i := len(blocks) - 1; i >= 0; i-- {
		t = append(t, blocks[i])
	}
	return t
}
