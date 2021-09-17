package routing

import (
	"github.com/kotfalya/hulk/app/ledger"
	"github.com/kotfalya/hulk/app/types"
)

type Table struct {
	self    types.Peer
	tick    ledger.Tick
	buckets []Bucket
}

func NewRoutingTable(self types.Peer, tick ledger.Tick) *Table {
	t := &Table{
		self:    self,
		tick:    tick,
		buckets: createBuckets(tick),
	}
	t.SetPeer(self)
	return t
}

func (rt *Table) GetPeer(target types.Addr) types.Peer {
	return rt.bucket(target).GetPeer(target)
}

func (rt *Table) SetPeer(peer types.Peer) {
	rt.bucket(peer.PK.ID().Addr()).SetPeer(peer)
}

func (rt *Table) bucket(target types.Addr) Bucket {
	cpl := rt.self.PK.ID().Addr().Cpl(target)
	for _, b := range rt.buckets {
		if cpl <= int(b.BitSizePrefix()+b.BitSize()) {
			return b
		}
	}
	return rt.buckets[len(rt.buckets)-1]
}

func createBuckets(tick ledger.Tick) []Bucket {
	buckets := make([]Bucket, len(tick))
	bitSizePrefix := uint8(0)
	for i, t := range tick {
		if i == len(tick)-1 {
			buckets[i] = NewFloatBucket(bitSizePrefix, t.BitSize)
		} else {
			buckets[i] = NewFixedBucket(bitSizePrefix, t.BitSize)
		}
		bitSizePrefix += t.BitSize
	}
	return buckets
}
