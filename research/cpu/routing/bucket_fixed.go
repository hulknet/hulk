package routing

import (
	"math"

	"github.com/kotfalya/hulk/research/cpu/types"
)

type FixedBucket struct {
	*BaseBucket
	peers  map[byte]types.Peer
	bitmap types.Bitmap256
}

func NewFixedBucket(bitSizePrefix uint8, bitSize uint8) *FixedBucket {
	return &FixedBucket{
		BaseBucket: &BaseBucket{bitSize: bitSize, bitSizePrefix: bitSizePrefix},
		peers:      make(map[byte]types.Peer, int(math.Pow(2, float64(bitSize)))),
	}
}

func (b *FixedBucket) GetPeer(target types.Addr) types.Peer {
	return b.peers[b.bucketAddr(target)]
}

func (b *FixedBucket) SetPeer(peer types.Peer) {
	bucketAddr := b.bucketAddr(peer.PK.ID().Addr())
	if b.bitmap.IsSet(bucketAddr) {
		return
	}
	b.peers[b.bucketAddr(peer.PK.ID().Addr())] = peer
	b.bitmap.Set(bucketAddr)
}

func (b *FixedBucket) Bitmap() types.Bitmap256 {
	return b.bitmap
}

func (b *BaseBucket) bucketAddr(target types.Addr) byte {
	leftShift := b.BitSizePrefix()
	rightShift := 64 - leftShift - b.BitSize()
	return byte((target << leftShift) >> rightShift)
}
