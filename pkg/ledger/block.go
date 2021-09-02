package ledger

import (
	"github.com/kotfalya/hulk/pkg/crypto"
	//"github.com/kotfalya/hulk/pkg/msg"
	"github.com/kotfalya/hulk/pkg/store"
)

type Block struct {
	ID   crypto.ID
	Prev crypto.ID
	Net  crypto.ID
	Sign crypto.Signature
}

//func BlockFromMsg(lb msg.Block) Block {
//	return Block{
//		ID:   lb.ID,
//		Prev: lb.Prev,
//		Net:  lb.Net,
//		Sign: lb.Sign,
//	}
//}
//
//func (b Block) ToMsg() msg.Block {
//	return msg.Block{
//		ID:   b.ID,
//		Prev: b.Prev,
//		Net:  b.Net,
//		Sign: b.Sign,
//	}
//}

func BlockFromStore(lb store.LedgerBlock) Block {
	return Block{
		ID:   lb.ID,
		Prev: lb.Prev,
		Net:  lb.Net,
		Sign: lb.Sign,
	}
}

func (b Block) ToStoreBlock() store.LedgerBlock {
	return store.LedgerBlock{
		ID:   b.ID,
		Prev: b.Prev,
		Net:  b.Net,
		Sign: b.Sign,
	}
}
