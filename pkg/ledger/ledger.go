package ledger

import (
	"github.com/asdine/storm/v3"
	"github.com/kotfalya/hulk/pkg/crypto"
	"github.com/kotfalya/hulk/pkg/store"
)

type Ledger struct {
	store *store.LedgerBlockStore
	last  Block
}

func NewLedger(db storm.Node) *Ledger {
	return &Ledger{
		store: store.NewLedgerBlockStore(db),
	}
}

func (l *Ledger) Version() crypto.ID {
	return l.last.ID
}

func (l *Ledger) LastBlock() Block {
	return l.last
}

func (l *Ledger) AddLedgerBlock(block Block) error {
	err := l.store.Add(block.ToStoreBlock())
	if err != nil {
		return err
	}

	l.last = block

	return nil
}

func (l *Ledger) CreateFirstLedgerBlock(netId crypto.ID) error {
	block := store.LedgerBlock{
		ID:   netId.WithSalt(netId[:]),
		Prev: netId,
		Net:  netId,
		Tick: 1,
	}

	err := l.store.Init(block)
	if err != nil {
		return err
	}

	l.last = BlockFromStore(block)

	return nil
}

func (l *Ledger) Load() error {
	block, err := l.store.Last()
	if err != nil {
		return err
	}

	l.last = BlockFromStore(block)

	return nil
}
