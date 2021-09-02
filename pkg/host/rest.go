package host

import (
	"github.com/kotfalya/hulk/pkg/net"
)

type StatusModel struct {
	ID         string                `json:"id"`
	Net        string                `json:"net_id"`
	Ledger     *net.LedgerBlockModel `json:"ledger"`
	NodesCount int                   `json:"nodes_count"`
}

func NewStatusModel(h *Host) *StatusModel {
	var netId string
	var ledger *net.LedgerBlockModel
	if h.net != nil {
		netId = h.net.ID().Hex()
		ledger = net.NewLedgerBlockModel(h.net.Ledger().LastBlock())
	}
	return &StatusModel{
		ID:         h.id.Hex(),
		Net:        netId,
		Ledger:     ledger,
		NodesCount: len(h.nodes),
	}
}

type NodeListModel struct {
	Nodes []*NodeIDModel `json:"nodes"`
}

func NewNodeListModel(h *Host) *NodeListModel {
	var nodes []*NodeIDModel
	for id, _ := range h.nodes {
		nodes = append(nodes, &NodeIDModel{
			ID: id.Hex(),
		})
	}
	return &NodeListModel{
		Nodes: nodes,
	}
}

type NodeIDModel struct {
	ID string `json:"id"`
}
