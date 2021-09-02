package host

import (
	"errors"

	"github.com/kotfalya/hulk/pkg/net"
	"github.com/kotfalya/hulk/pkg/store"
)

var (
	ErrNetExists       = errors.New("net already exists")
	ErrNetDoesNotExist = errors.New("net does not exists")
)

func (h *Host) loadNet() error {
	nm, err := store.LoadNet(h.db)
	if err != nil {
		return err
	}

	if nm == nil {
		return nil
	}

	h.net = net.NewNet(nm.ID, h.db.From("net"))

	return h.net.Load()
}

func (h *Host) CreateNet() error {
	if h.net != nil {
		return ErrNetExists
	}

	nm := &store.Net{
		ID:       h.cfg.Crypto.HostID.WithSalt(h.cfg.Crypto.UserID[:]), // can be random
		AuthorID: h.cfg.Crypto.UserID,
	}

	err := h.db.Set("host", "net", nm)
	if err != nil {
		return err
	}

	h.net = net.NewNet(nm.ID, h.db.From("net"))
	if err = h.net.Ledger().CreateFirstLedgerBlock(h.net.ID()); err != nil {
		return err
	}

	return nil
}

// Not implemented
func (h *Host) JoinNet() error {
	if h.net != nil {
		return ErrNetExists
	}

	return nil
}

// Not implemented
func (h *Host) LeaveNet() error {
	if h.net == nil {
		return ErrNetExists
	}

	return nil
}
