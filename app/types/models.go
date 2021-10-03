package types

import (
	"github.com/vmihailenco/msgpack/v5"
)

type MessageHeader struct {
	ID      ID64
	To      ID64
	From    ID64
	Token   Token
	BlockID ID64
	Time    Time
	Sign    []Sign520
}

type BaseMessage struct {
	Type string
	Data msgpack.RawMessage
	Sign []byte
}

type NetPartition struct {
	Size byte
}

type Partition struct {
	Position byte
	Length   byte
}

type Replica struct {
	Max uint64
	Min uint64
}

type Precursor struct {
	Min uint8
	Mac uint8
}
