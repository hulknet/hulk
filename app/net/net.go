package net

import (
	"sync"

	"github.com/kotfalya/hulk/app/routing"
	"github.com/kotfalya/hulk/app/types"
)

type allowList map[types.Token]types.Peer

type Net struct {
	mu        sync.RWMutex
	self      types.Peer
	table     *routing.Table
	handler   *MessageHandler
	allowList map[types.Token]types.Peer
}

func NewNet(self types.Peer) *Net {
	return &Net{
		mu:        sync.RWMutex{},
		self:      self,
		allowList: createAllowLost(self),
	}
}

func (n *Net) Init(state types.State) {
	n.table = routing.NewRoutingTable(n.self, state.Head().BitSize)
	n.handler = NewMessageHandler(state)
}

func (n *Net) Start() error {
	return n.handler.Start()
}

func (n *Net) AddPeer(peer types.Peer) {
	n.table.SetPeer(peer)
}

func (n *Net) FindPeer(target types.ID64) types.Peer {
	return n.table.GetPeer(target)
}

func (n *Net) CheckToken(token types.Token) bool {
	_, ok := n.allowList[token]
	return ok
}

func (n *Net) IsSelf(token types.Token) bool {
	peer, ok := n.allowList[token]
	return ok && n.self.Equal(peer)
}

func (n *Net) Self() types.Peer {
	return n.self
}

func (n *Net) HandleMessage(header types.MessageHeader, data []byte) error {
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

func createAllowLost(self types.Peer) allowList {
	peers := make(map[types.Token]types.Peer)
	peers[self.Token] = self
	return peers
}
