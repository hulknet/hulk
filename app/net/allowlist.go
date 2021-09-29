package net

import "github.com/kotfalya/hulk/app/types"

type AllowList map[types.Token]types.Peer

func (al AllowList) CheckToken(token types.Token) bool {
	_, ok := al[token]
	return ok
}

func (al AllowList) AddPeer(peer types.Peer) bool {
	_, ok := al[peer.Token]
	if ok {
		return false
	}
	al[peer.Token] = peer
	return true
}

func (al AllowList) FindPeer(token types.Token) (types.Peer, bool) {
	peer, ok := al[token]
	return peer, ok
}

func createAllowList(self types.Peer) AllowList {
	peers := make(map[types.Token]types.Peer)
	peers[self.Token] = self
	return peers
}
