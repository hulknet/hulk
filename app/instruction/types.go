package instruction

import (
	"fmt"

	"github.com/vmihailenco/msgpack/v5"

	"github.com/hulknet/hulk/app/types"
)

type (
	MessageType  string
	ResourceType string
	MethodType   string
)

const (
	InitMessageType     = MessageType("init")
	AnnounceMessageType = MessageType("announce")
	InputMessageType    = MessageType("input")
	UnknownMessageType  = MessageType("unknown")
)

func MessageTypeFromString(s string) MessageType {
	switch s {
	case string(InitMessageType):
		return InitMessageType
	case string(AnnounceMessageType):
		return AnnounceMessageType
	case string(InputMessageType):
		return InputMessageType
	default:
		return UnknownMessageType
	}
}

type AnnounceMessage struct {
	Resource ResourceType
	Method   MethodType
	Data     []byte
}

func UnmarshalAnnounceMessage(data []byte) (announceMessage AnnounceMessage, err error) {
	if err = msgpack.Unmarshal(data, &announceMessage); err != nil {
		err = fmt.Errorf("failed to unmarshal AnnounceMessage: %w \n", err)
	}
	return
}

type InputMessage struct {
	Data     []byte
	Position byte
}

func UnmarshalInputMessage(data []byte) (inputMessage InputMessage, err error) {
	if err = msgpack.Unmarshal(data, &inputMessage); err != nil {
		err = fmt.Errorf("failed to unmarshal InputMessage: %w \n", err)
	}
	return
}

type Instruction struct {
	id       types.ID64
	time     types.Time
	resource ResourceType
	method   MessageType
	input    []Input
	data     []byte
	meta     Meta
}

func NewInstruction(id types.ID64, time types.Time) *Instruction {
	return &Instruction{id: id, time: time}
}

func (i Instruction) Announce(am AnnounceMessage) {
}

func (i Instruction) Input(im InputMessage) {
}

type Input struct {
	Position byte
	Shards   [][]byte
}

type Meta struct {
	Shards    []types.Time
	Announces []types.Time
}
