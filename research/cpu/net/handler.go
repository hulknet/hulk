package net

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/kotfalya/hulk/research/cpu/ledger"
	"github.com/kotfalya/hulk/research/cpu/types"
)

type MessageItem struct {
	id   types.ID
	part types.Partition
	data []byte
}

type Message struct {
	id       types.ID
	length   uint64
	received uint64
	messages [][]byte
}

func (m *Message) Update(position uint64, data []byte) bool {
	if len(m.messages[position]) > 0 {
		return false
	}

	// waits for msgpack
	data, err := hex.DecodeString(string(data))
	if err != nil {
		fmt.Errorf("%v", err)
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
			fmt.Errorf("%v", err)
		}

		m.messages[mi.part.Position] = data
	} else {
		m.length = 1
		m.messages = [][]byte{mi.data}
	}

	return
}

type MessageState struct {
	messages map[types.ID]Message
	resolved map[types.ID]struct{}
}

func newMessageState() *MessageState {
	return &MessageState{
		messages: make(map[types.ID]Message, 0),
		resolved: make(map[types.ID]struct{}, 0),
	}
}

func (s *MessageState) IsMessageResolved(id types.ID) bool {
	_, ok := s.resolved[id]
	return ok
}

func (s *MessageState) Resolve(id types.ID) {
	s.resolved[id] = struct{}{}
	delete(s.messages, id)
}

func (s *MessageState) CreateOrUpdate(mi MessageItem) (m Message) {
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

type MessageHandler struct {
	tickCh       chan ledger.Tick
	messageCh    chan MessageItem
	tick         ledger.Tick
	state        *MessageState
	stateArchive map[types.ID]*MessageState
	processor    func(m Message)
}

func NewMessageHandler(tick ledger.Tick) *MessageHandler {
	return &MessageHandler{
		tick:      tick,
		state:     newMessageState(),
		processor: createProcessor(tick),
		tickCh:    make(chan ledger.Tick, 1),
		messageCh: make(chan MessageItem, 10),
	}
}

func (h *MessageHandler) SetTick(tick ledger.Tick) {
	h.tickCh <- tick
}

func (h *MessageHandler) Message(id types.ID, part types.Partition, data []byte) {
	h.messageCh <- MessageItem{id, part, data}
}

func (h *MessageHandler) Start() error {
	for {
		select {
		case tick := <-h.tickCh:
			h.stateArchive[h.tick.NodeBlock().ID] = h.state
			h.state = newMessageState()
			h.tick = tick
		case mi := <-h.messageCh:
			if h.state.IsMessageResolved(mi.id) {
				continue
			}

			m := h.state.CreateOrUpdate(mi)
			if m.Assembled() {
				go h.processor(m)
				h.state.Resolve(m.id)
			}
		}
	}
}

func createProcessor(tick ledger.Tick) Processor {
	return func(m Message) {
		if m.length > 1 {
			d, err := types.DecryptFromParts(m.messages)
			if err != nil {
				fmt.Errorf("failed to decode message: %v", err)
			}

			var baseMessage types.BaseMessage
			if err = json.Unmarshal(d, &baseMessage); err != nil {
				fmt.Errorf("failed to unmarshal BaseMessage: %v", err)
			}

			ok, err := types.CheckStringSignature(*baseMessage.Data, baseMessage.Sign)
			if err != nil {
				fmt.Errorf("failed to check signature: %v", err)
			} else if !ok {
				fmt.Errorf("signature is invalid: %v", baseMessage)
			}

			fmt.Println(baseMessage.Type)
		} else {
			fmt.Println(string(m.messages[0]))
		}
	}
}
