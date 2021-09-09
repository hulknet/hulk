package net

import (
	"sync"

	"github.com/kotfalya/hulk/research/cpu/ledger"
	"github.com/kotfalya/hulk/research/cpu/routing"
	"github.com/kotfalya/hulk/research/cpu/types"
)

type allowList map[types.Token]types.PeerIn

type Net struct {
	mu        sync.RWMutex
	self      types.PeerOut
	table     *routing.Table
	allowList map[types.Token]types.PeerIn
}

func NewNet(self types.PeerOut) *Net {
	return &Net{
		mu:        sync.RWMutex{},
		self:      self,
		allowList: createAllowLost(types.PeerOutToIn(self)),
	}
}

func (n *Net) SetTick(tick ledger.Tick) {
	n.table = routing.NewRoutingTable(n.self, tick)
}

func (n *Net) AddPeer(peer types.PeerOut) {
	n.table.SetPeer(peer)
}

func (n *Net) FindPeer(target types.Addr) types.PeerOut {
	return n.table.GetPeer(target)
}

func (n *Net) CheckToken(token types.Token) bool {
	_, ok := n.allowList[token]
	return ok
}

func (n *Net) CheckPeer(peer types.PeerIn) bool {
	peerIn, ok := n.allowList[peer.Token]
	return ok && peer.Equal(peerIn)
}

func (n *Net) Self() types.PeerOut {
	return n.self
}

func (n *Net) HandleMessage(header types.MessageHeader, data []byte) error {
	//if !rh.net.CheckPeer(peerIn) {
	//	w.WriteHeader(http.StatusForbidden)
	//	return
	//}
	//
	//peer := rh.net.FindPeer(messageHeader.To)
	return nil
}

func createAllowLost(self types.PeerIn) allowList {
	peers := make(map[types.Token]types.PeerIn)
	peers[self.Token] = self
	return peers
}
