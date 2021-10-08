package types

import (
	"bytes"

	"github.com/vmihailenco/msgpack/v5"
)

type NetMessage struct {
	ID    ID64
	Addr  ID64
	Time  Time
	Token Token
	Data  []byte
	Sign  []Sign520
}

func (msg *NetMessage) Encode() []byte {
	data := bytes.NewBuffer(msg.ID.Bytes())
	data.Write(msg.Addr.Bytes())
	data.Write(msg.Time.Encode())
	data.Write(msg.Data)
	return data.Bytes()
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
