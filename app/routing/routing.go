package routing

import (
	"github.com/kotfalya/hulk/app/types"
)

type Table struct {
	id      types.ID64
	buckets []Bucket
}

func NewTable(block types.Block, self types.Peer) *Table {
	c := &Table{
		id:      block.ID,
		buckets: createBuckets(block.BitSize),
	}
	c.SetPeer(self)

	return c
}

func (rt *Table) GetPeer(target types.ID64) types.Peer {
	return rt.bucket(target).GetPeer(target)
}

func (rt *Table) SetPeer(peer types.Peer) {
	rt.bucket(peer.Pub.ID()).SetPeer(peer)
}

func (rt *Table) bucket(target types.ID64) Bucket {
	cpl := types.Cpl(rt.id.Bytes(), target.Bytes())
	for _, b := range rt.buckets {
		if cpl <= int(b.BitSizePrefix()+b.BitSize()) {
			return b
		}
	}
	return rt.buckets[len(rt.buckets)-1]
}

func createBuckets(bitSizeList []uint8) []Bucket {
	buckets := make([]Bucket, len(bitSizeList))
	bitSizePrefix := uint8(0)
	for i, s := range bitSizeList {
		if i == len(bitSizeList)-1 {
			buckets[i] = NewFloatBucket(bitSizePrefix, s)
		} else {
			buckets[i] = NewFixedBucket(bitSizePrefix, s)
		}
		bitSizePrefix += s
	}
	return buckets
}
