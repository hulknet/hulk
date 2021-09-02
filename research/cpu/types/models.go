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
	Sign []Sign           `json:"sign"`
}

type Partition struct {
	Length   byte `json:"length"`
	Position byte `json:"position"`
}

type Replica struct {
	Max byte
	Min byte
}

type Shard struct {
	Len byte
	Num byte
}
