package net

import (
	"github.com/asdine/storm/v3"
	"github.com/kotfalya/hulk/pkg/crypto"
	"github.com/kotfalya/hulk/pkg/ledger"
	log "github.com/sirupsen/logrus"
)

type Net struct {
	id     crypto.ID
	log    *log.Entry
	ledger *ledger.Ledger
}

func NewNet(id crypto.ID, db storm.Node) *Net {
	return &Net{
		id:     id,
		ledger: ledger.NewLedger(db),
	}
}

func (n *Net) Ledger() *ledger.Ledger {
	return n.ledger
}

func (n *Net) ID() crypto.ID {
	return n.id
}

func (n *Net) Load() error {
	return n.ledger.Load()
}
