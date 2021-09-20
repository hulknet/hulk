package routing

import (
	"github.com/kotfalya/hulk/app/types"
)

type Table struct {
	self    types.Peer
	buckets []Bucket
}

func NewRoutingTable(self types.Peer, bitSize []uint8) *Table {
	t := &Table{
		self:    self,
		buckets: createBuckets(bitSize),
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

func createBuckets(bitSize []uint8) []Bucket {
	buckets := make([]Bucket, len(bitSize))
	bitSizePrefix := uint8(0)
	for i, s := range bitSize {
		if i == len(bitSize)-1 {
			buckets[i] = NewFloatBucket(bitSizePrefix, s)
		} else {
			buckets[i] = NewFixedBucket(bitSizePrefix, s)
		}
		bitSizePrefix += s
	}
	return buckets
}
