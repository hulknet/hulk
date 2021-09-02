package host

import (
	"context"

	"github.com/asdine/storm/v3"
	"github.com/kotfalya/hulk/pkg/config"
	"github.com/kotfalya/hulk/pkg/crypto"
	"github.com/kotfalya/hulk/pkg/net"
	"github.com/kotfalya/hulk/pkg/node"
)

type Host struct {
	id    crypto.ID
	ctx   context.Context
	cfg   *config.Config
	db    storm.Node
	net   *net.Net
	nodes map[crypto.ID]*node.Node
}

func NewHost(cfg *config.Config, db storm.Node) *Host {
	return &Host{
		ctx:   context.Background(),
		cfg:   cfg,
		id:    cfg.Crypto.HostID,
		db:    db,
		nodes: make(map[crypto.ID]*node.Node),
	}
}

func (h *Host) Load() error {
	err := h.loadNet()
	if err != nil {
		return err
	}

	if h.net == nil {
		return nil
	}

	err = h.loadNodes()
	if err != nil {
		return err
	}

	if len(h.nodes) == 0 {
		_, err = h.AddNode()
		if err != nil {
			return err
		}
	}

	for _, n := range h.nodes {
		n.RegisterBlock(h.net.Ledger().LastBlock())
	}

	return nil
}

func (h *Host) ID() crypto.ID {
	return h.id
}

func (h *Host) Addr() string {
	return h.cfg.Transport.Address
}
