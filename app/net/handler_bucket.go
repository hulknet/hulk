package net

import (
	"encoding/hex"
	"fmt"

	"github.com/vmihailenco/msgpack/v5"

	"github.com/hulknet/hulk/app/types"
)

type MessageItem struct {
	id   types.ID64
	part types.Partition
	data []byte
}

type Message struct {
	id       types.ID64
	length   byte
	received byte
	items    [][]byte
}

func (m *Message) Update(position byte, data []byte) bool {
	if len(m.items[position]) > 0 {
		return false
	}

	m.items[position] = data
	m.received++

	return true
}

func (m Message) Assembled() bool {
	return m.length == m.received
}

func newMessage(mi MessageItem) (m Message) {
	m.id = mi.id
	m.received = 1
	m.length = mi.part.Length
	m.items = make([][]byte, mi.part.Length)
	m.items[mi.part.Position] = mi.data
	return
}

type MessageChunks struct {
	messages map[types.ID64]Message
	resolved map[types.ID64]struct{}
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
	chunks    *MessageChunks
	processor func(m Message)
}

func NewBucketHandler() *BucketHandler {
	return &BucketHandler{
		chunks:    newMessageChunks(),
		processor: createProcessor(),
		messageCh: make(chan MessageItem, 10),
	}
}

func (h *BucketHandler) Message(header types.MessageHeader, data []byte) {
	h.messageCh <- MessageItem{header.ID, header.Part, data}
}

func (h *BucketHandler) Start() {
	for {
		select {
		case mi, ok := <-h.messageCh:
			if !ok {
				return
			}
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

func (h *BucketHandler) Stop() {
	close(h.messageCh)
}

func createProcessor() Processor {
	return func(m Message) {
		d, err := types.DecryptFromParts(m.items)
		if err != nil {
			fmt.Printf("error decodingfailed to decode message: %v \n", err)
		}

		var baseMessage types.BaseMessage
		if err = msgpack.Unmarshal(d, &baseMessage); err != nil {
			fmt.Printf("failed to unmarshal BaseMessage: %v \n", err)
		}

		ok, err := types.CheckSignature(baseMessage.Data, baseMessage.Sign)
		if err != nil {
			fmt.Printf("failed to check signature: %v \n", err)
		} else if !ok {
			fmt.Printf("signature is invalid: %v \n", baseMessage)
		}

		fmt.Println(baseMessage.Type)
		fmt.Println(hex.EncodeToString(baseMessage.Data))
	}
}
