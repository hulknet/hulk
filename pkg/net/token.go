package net

import "github.com/kotfalya/hulk/pkg/crypto"

type Token struct {
	ID       crypto.ID
	LedgerID crypto.ID
	NodeID   crypto.ID
	Sign     crypto.Signature
}
