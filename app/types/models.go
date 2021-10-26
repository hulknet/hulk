package types

import (
	"bytes"
	"fmt"

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
	ID   ID64
	Time Time
	Type string
	Data []byte
	Sign Sign520
}

func UnmarshalBaseMessage(data []byte) (bm BaseMessage, err error) {
	if err = msgpack.Unmarshal(data, &bm); err != nil {
		err = fmt.Errorf("failed to unmarshal BaseMessage: %w \n", err)
	}
	return
}

type NetPartition struct {
	Size byte
}

type Replica struct {
	Min uint64
	Exp uint64
}

type Precursor struct {
	Min uint8
	Exp uint8
}
