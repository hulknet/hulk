package store

import (
	"github.com/asdine/storm/v3"
	"github.com/kotfalya/hulk/pkg/crypto"
)

type LedgerBlock struct {
	ID   crypto.ID `storm:"id"`
	Prev crypto.ID `storm:"index"`
	Net  crypto.ID `storm:"index"`
	Sign crypto.Signature
	Tick uint64 `storm:"index"`
	Size uint64
}

type LedgerBlockStore struct {
	db storm.Node
}

func NewLedgerBlockStore(db storm.Node) *LedgerBlockStore {
	return &LedgerBlockStore{
		db: db,
	}
}

func (lbs *LedgerBlockStore) Add(block LedgerBlock) error {
	prevBlock, err := lbs.FindOne(block.Prev)
	if err != nil {
		return err
	}
	block.Tick = prevBlock.Tick + 1

	return lbs.db.Save(&block)
}

func (lbs *LedgerBlockStore) Init(block LedgerBlock) error {
	return lbs.db.Save(&block)
}

func (lbs *LedgerBlockStore) FindOne(id crypto.ID) (LedgerBlock, error) {
	var block LedgerBlock
	err := lbs.db.One("ID", id, &block)

	return block, err
}

func (lbs *LedgerBlockStore) FindOneByTick(tick uint64) (LedgerBlock, error) {
	var block LedgerBlock
	err := lbs.db.One("Tick", tick, &block)

	return block, err
}

func (lbs *LedgerBlockStore) Last() (LedgerBlock, error) {
	var blocks []LedgerBlock

	err := lbs.db.AllByIndex("Tick", &blocks, storm.Limit(1), storm.Reverse())
	if err != nil {
		return LedgerBlock{}, err
	}

	if len(blocks) == 0 {
		return LedgerBlock{}, storm.ErrNotFound
	}

	return blocks[0], nil
}

func (lbs *LedgerBlockStore) Chain(block LedgerBlock, num int) ([]LedgerBlock, error) {
	chain := []LedgerBlock{block}
	for i := 1; i < num; i++ {
		prevBlock, err := lbs.FindOne(chain[len(chain)-1].Prev)
		if err != nil {
			return chain, err
		}

		chain = append(chain, prevBlock)
	}

	return chain, nil
}
