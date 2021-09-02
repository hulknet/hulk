package node

import (
	"context"

	"github.com/kotfalya/hulk/pkg/crypto"
	"github.com/kotfalya/hulk/pkg/ledger"
)

type Node struct {
	ctx    context.Context
	id     crypto.ID
	hostId crypto.ID
	block  ledger.Block
}

func NewNode(ctx context.Context, id crypto.ID) *Node {
	node := &Node{
		ctx: ctx,
		id:  id,
	}

	return node
}

func (n *Node) Status() ([]byte, error) {
	return []byte("Ok"), nil
}

func (n *Node) ID() crypto.ID {
	return n.id
}
