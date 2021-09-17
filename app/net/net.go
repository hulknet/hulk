package net

import (
	"sync"

	"github.com/kotfalya/hulk/app/ledger"
	"github.com/kotfalya/hulk/app/routing"
	types2 "github.com/kotfalya/hulk/app/types"
)

type allowList map[types2.Token]types2.Peer

type Net struct {
	mu        sync.RWMutex
	self      types2.Peer
	table     *routing.Table
	handler   *MessageHandler
	allowList map[types2.Token]types2.Peer
}

func NewNet(self types2.Peer) *Net {
	return &Net{
		mu:        sync.RWMutex{},
		self:      self,
		allowList: createAllowLost(self),
	}
}

func (n *Net) Init(tick ledger.Tick) {
	n.table = routing.NewRoutingTable(n.self, tick)
	n.handler = NewMessageHandler(tick)
}

func (n *Net) Start() error {
	return n.handler.Start()
}

func (n *Net) SetTick(tick ledger.Tick) {
	//todo: rotate table on tick
	//n.table = routing.NewRoutingTable(n.self, tick)
}

func (n *Net) AddPeer(peer types2.Peer) {
	n.table.SetPeer(peer)
}

func (n *Net) FindPeer(target types2.Addr) types2.Peer {
	return n.table.GetPeer(target)
}

func (n *Net) CheckToken(token types2.Token) bool {
	_, ok := n.allowList[token]
	return ok
}

func (n *Net) IsSelf(token types2.Token) bool {
	peer, ok := n.allowList[token]
	return ok && n.self.Equal(peer)
}

func (n *Net) Self() types2.Peer {
	return n.self
}

func (n *Net) HandleMessage(header types2.MessageHeader, data []byte) error {
	if n.IsSelf(header.Token) {
		n.handler.Message(header.ID, header.Part, data)
	} else {
		//todo: implement proxy client
	}
	//if !rh.net.CheckPeer(peerIn) {
	//	w.WriteHeader(http.StatusForbidden)
	//	return
	//}
	//
	//peer := rh.net.FindPeer(messageHeader.To)
	return nil
}

func createAllowLost(self types2.Peer) allowList {
	peers := make(map[types2.Token]types2.Peer)
	peers[self.Token] = self
	return peers
}
