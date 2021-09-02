package router

import "github.com/kotfalya/hulk/pkg/crypto"

type Peer struct {
	ID         crypto.ID
	Addr       Addr
	BucketAddr BucketAddr
	Rank       uint64
	Dest       string
	Token      string
}
