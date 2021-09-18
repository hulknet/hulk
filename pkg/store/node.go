package store

import (
	"github.com/asdine/storm/v3"
	"github.com/kotfalya/hulk/pkg/crypto"
)

type Node struct {
	ID    crypto.ID `storm:"id"`
	Index uint64
}

func LoadNodes(db storm.Node) ([]Node, error) {
	var nodes []Node
	err := db.From("nodes").All(&nodes)

	return nodes, err
}

func SaveNode(db storm.Node, node Node) error {
	return db.From("nodes").Save(&node)
}

func DeleteNode(db storm.Node, node Node) error {
	return db.From("nodes").DeleteStruct(&node)
}
