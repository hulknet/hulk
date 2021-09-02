package node

import "github.com/kotfalya/hulk/pkg/ledger"

func (n *Node) RegisterBlock(block ledger.Block) {
	n.block = block
}

func (n *Node) ActiveBlocks() []ledger.Block {
	return nil
}

func (n *Node) FutureBlock() ledger.Block {
	return ledger.Block{}
}
