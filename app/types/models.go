package types

import "encoding/json"

type MessageHeader struct {
	ID    ID
	To    ShortID
	From  ShortID
	Token Token
	Time  Time
	Part  Partition
	Sign  []Sign
}

type BaseMessage struct {
	Type string           `json:"type"`
	Data *json.RawMessage `json:"data"`
	Sign string           `json:"sign"`
}

type Partition struct {
	Position byte
	Length   byte
}

type Replica struct {
	Max uint64
	Min uint64
}

type Shard struct {
	Len uint64
	Num uint64
}
