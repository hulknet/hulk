package routing

import (
	"github.com/kotfalya/hulk/app/types"
)

type Bucket interface {
	GetPeer(addr types.ID64) types.Peer
	SetPeer(peer types.Peer)
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
