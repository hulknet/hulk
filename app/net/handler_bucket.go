package net

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/kotfalya/hulk/app/types"
)

type MessageItem struct {
	id   types.ID64
	time types.Time
	part types.Partition
	data []byte
}

type Message struct {
	id       types.ID64
	length   byte
	received byte
	messages [][]byte
}

func (m *Message) Update(position byte, data []byte) bool {
	if len(m.messages[position]) > 0 {
		return false
	}

	// waits for msgpack
	data, err := hex.DecodeString(string(data))
	if err != nil {
		fmt.Println(err)
	}

	m.messages[position] = data
	m.received++

	return true
}

func (m Message) Assembled() bool {
	return m.length == m.received
}

func newMessage(mi MessageItem) (m Message) {
	m.id = mi.id
	m.received = 1

	if mi.part.Length > 1 {
		m.length = mi.part.Length
		m.messages = make([][]byte, mi.part.Length)

		// waits for msgpack
		data, err := hex.DecodeString(string(mi.data))
		if err != nil {
			fmt.Println(err)
		}

		m.messages[mi.part.Position] = data
	} else {
		m.length = 1
		m.messages = [][]byte{mi.data}
	}

	return
}

type MessageChunks struct {
	messages      map[types.ID64]Message
	messageToTime map[types.ID64]types.Time
	resolved      map[types.ID64]struct{}
}

func newMessageChunks() *MessageChunks {
	return &MessageChunks{
		messages: make(map[types.ID64]Message, 0),
		resolved: make(map[types.ID64]struct{}, 0),
	}
}

func (s *MessageChunks) IsMessageResolved(id types.ID64) bool {
	_, ok := s.resolved[id]
	return ok
}

func (s *MessageChunks) Resolve(id types.ID64) {
	s.resolved[id] = struct{}{}
	delete(s.messages, id)
}

func (s *MessageChunks) CreateOrUpdate(mi MessageItem) (m Message) {
	m, ok := s.messages[mi.id]
	if ok {
		m.Update(mi.part.Position, mi.data)
	} else {
		m = newMessage(mi)
		s.messages[m.id] = m
	}

	return
}

type Processor func(m Message)

type BucketHandler struct {
	messageCh chan MessageItem
	state     types.State
	chunks    *MessageChunks
	processor func(m Message)
}

func NewMessageHandler(state types.State) *BucketHandler {
	return &BucketHandler{
		state:     state,
		chunks:    newMessageChunks(),
		processor: createProcessor(),
		messageCh: make(chan MessageItem, 10),
	}
}

func (h *BucketHandler) UpdateState(state types.State) {

}

func (h *BucketHandler) Message(header types.MessageHeader, data []byte) {
	h.messageCh <- MessageItem{header.ID, header.Time, header.Part, data}
}

func (h *BucketHandler) Start() error {
	for {
		select {
		case mi := <-h.messageCh:
			if h.chunks.IsMessageResolved(mi.id) {
				continue
			}

			m := h.chunks.CreateOrUpdate(mi)
			if m.Assembled() {
				go h.processor(m)
				h.chunks.Resolve(m.id)
			}
		}
	}
}

func createProcessor() Processor {
	return func(m Message) {
		if m.length > 1 {
			d, err := types.DecryptFromParts(m.messages)
			if err != nil {
				fmt.Printf("error decodingfailed to decode message: %v \n", err)
			}

			var baseMessage types.BaseMessage
			if err = json.Unmarshal(d, &baseMessage); err != nil {
				fmt.Printf("failed to unmarshal BaseMessage: %v \n", err)
			}

			ok, err := types.CheckStringSignature(*baseMessage.Data, baseMessage.Sign)
			if err != nil {
				fmt.Printf("failed to check signature: %v \n", err)
			} else if !ok {
				fmt.Printf("signature is invalid: %v \n", baseMessage)
			}

			fmt.Println(baseMessage.Type)
		} else {
			fmt.Println(string(m.messages[0]))
		}
	}
}
