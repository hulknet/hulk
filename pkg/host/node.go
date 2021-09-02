package host

import (
	"bytes"
	"encoding/binary"

	"github.com/kotfalya/hulk/pkg/crypto"
	"github.com/kotfalya/hulk/pkg/node"
	"github.com/kotfalya/hulk/pkg/store"
)

func (h *Host) AddNode() (*node.Node, error) {
	nodeId, err := h.generatedNextNodeID()
	if err != nil {
		return nil, err
	}

	n := node.NewNode(h.ctx, nodeId)
	h.nodes[nodeId] = n

	return n, store.SaveNode(h.db, store.Node{ID: nodeId})
}

func (h *Host) FindNode(nodeId crypto.ID) (*node.Node, bool) {
	n, ok := h.nodes[nodeId]

	return n, ok
}

func (h *Host) loadNodes() error {
	ns, err := store.LoadNodes(h.db)
	if err != nil {
		return err
	}

	for _, n := range ns {
		h.nodes[n.ID] = node.NewNode(h.ctx, n.ID)
	}

	return nil
}

func (h *Host) generatedNextNodeID() (crypto.ID, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, uint64(len(h.nodes)))
	if err != nil {
		return crypto.ID{}, err
	}

	return h.cfg.Crypto.HostID.WithSalt(buf.Bytes()), nil
}
