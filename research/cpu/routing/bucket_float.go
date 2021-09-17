package routing

import (
	"math"
	"sort"

	"github.com/kotfalya/hulk/research/cpu/types"
)

type FloatBucket struct {
	*BaseBucket
	peers []types.Peer
}

func NewFloatBucket(bitSizePrefix uint8, bitSize uint8) *FloatBucket {
	return &FloatBucket{
		BaseBucket: &BaseBucket{bitSize: bitSize, bitSizePrefix: bitSizePrefix},
	}
}

func (b *FloatBucket) GetPeer(target types.Addr) types.Peer {
	if len(b.peers) == 1 {
		return b.peers[0]
	}
	baseAngle := 1 / float64(len(b.peers))
	index := int(math.Round(target.Normalize(b.bitSizePrefix) / baseAngle))
	return b.peers[index]
}

func (b *FloatBucket) SetPeer(peer types.Peer) {
	b.peers = append(b.peers, peer)
	sort.Sort(peerByAddr(b.peers))
}

type peerByAddr []types.Peer

func (a peerByAddr) Len() int           { return len(a) }
func (a peerByAddr) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a peerByAddr) Less(i, j int) bool { return a[i].PK.ID().Addr() < a[j].PK.ID().Addr() }
