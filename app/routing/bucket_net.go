package routing

import (
	"math"

	"github.com/kotfalya/hulk/app/types"
)

type NetBucket struct {
	*BaseBucket
	peers  map[byte]types.Peer
	bitmap types.Bitmap256
}

func NewFixedBucket(bitSizePrefix uint8, bitSize uint8) *NetBucket {
	return &NetBucket{
		BaseBucket: &BaseBucket{bitSize: bitSize, bitSizePrefix: bitSizePrefix},
		peers:      make(map[byte]types.Peer, int(math.Pow(2, float64(bitSize)))),
	}
}

func (b *NetBucket) GetPeer(target types.ID64) types.Peer {
	return b.peers[b.bucketAddr(target.Uint64())]
}

func (b *NetBucket) SetPeer(peer types.Peer) {
	bucketAddr := b.bucketAddr(peer.Pub.ID256().Uint64())
	if b.bitmap.IsSet(bucketAddr) {
		return
	}
	b.peers[b.bucketAddr(peer.Pub.ID256().Uint64())] = peer
	b.bitmap.Set(bucketAddr)
}

func (b *NetBucket) Bitmap() types.Bitmap256 {
	return b.bitmap
}

func (b *BaseBucket) bucketAddr(target uint64) byte {
	leftShift := b.BitSizePrefix()
	rightShift := 64 - leftShift - b.BitSize()
	return byte((target << leftShift) >> rightShift)
}
