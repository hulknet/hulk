package router

import "math"

type Bucket struct {
	peers   map[BucketAddr]Peer
	jury    map[BucketAddr]Peer
	bitMask uint64
}

func newBucket(bitSize int) *Bucket {
	return &Bucket{
		peers: make(map[BucketAddr]Peer, math.Pow(2, float64(bitSize))),
		jury:  make(map[BucketAddr]Peer, bitSize),
	}
}

func (b Bucket) Get(subAddr BucketAddr) (Peer, bool) {
	peer, ok := b.peers[subAddr]

	return peer, ok
}

func (b Bucket) Set(peer Peer) {
	b.peers[peer.BucketAddr] = peer
}
