package types

import "encoding/json"

type MessageType string

type MessageHeader struct {
	ID    ID
	To    Addr
	Token Token
	Part  Partition
	Sign  []Sign
}

type BaseMessage struct {
	Type MessageType      `json:"type"`
	Data *json.RawMessage `json:"data"`
	Part Partition        `json:"part"`
	Sign []Sign           `json:"sign"`
}

type Partition struct {
	Position uint64
	Length   uint64
	Enabled  bool
}

type Replica struct {
	Max uint64
	Min uint64
}

type Shard struct {
	Len uint64
	Num uint64
}
