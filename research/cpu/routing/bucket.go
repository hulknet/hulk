package routing

import (
	"github.com/kotfalya/hulk/research/cpu/types"
)

type Bucket interface {
	GetPeer(addr types.Addr) types.PeerOut
	SetPeer(peer types.PeerOut)
	BitSize() uint8
	BitSizePrefix() uint8
}

type BaseBucket struct {
	bitSize       uint8
	bitSizePrefix uint8
}

func (b *BaseBucket) BitSize() uint8 {
	return b.bitSize
}

func (b *BaseBucket) BitSizePrefix() uint8 {
	return b.bitSizePrefix
}
