package net

import (
	"github.com/kotfalya/hulk/pkg/ledger"
)

type LedgerBlockModel struct {
	ID     string `json:"id"`
	PrevID string `json:"prev_id"`
	Sign   string `json:"sign"`
	Size   uint64 `json:"size"`
}

func NewLedgerBlockModel(b ledger.Block) *LedgerBlockModel {
	return &LedgerBlockModel{
		ID:     b.ID.Hex(),
		PrevID: b.Prev.Hex(),
		Sign:   b.Sign.Hex(),
		//Size:   b.Size,
	}
}
