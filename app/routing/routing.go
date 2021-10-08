package routing

import (
	"github.com/hulknet/hulk/app/types"
)

type Table struct {
	id      types.ID64
	buckets []Bucket
}

func NewTable(state types.State) *Table {
	c := &Table{
		id:      state.Block().ID,
		buckets: createBuckets(state.Block().BitSize),
	}
	c.SetPeer(state.Peer())

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
