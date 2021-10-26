package net

import (
	"encoding/hex"

	log "github.com/sirupsen/logrus"

	"github.com/hulknet/hulk/app/types"
)

type MessageChunk struct {
	id   types.ID64
	time types.Time
	data []byte
}

type Message struct {
	id     types.ID64
	time   types.Time
	chunks [][]byte
}

func (m *Message) Update(data []byte) bool {
	if len(m.chunks[data[0]]) > 0 {
		return false
	}
	m.chunks[data[0]] = data
	return true
}

func (m Message) Assembled() bool {
	for _, chunk := range m.chunks {
		if chunk == nil {
			return false
		}
	}
	return true
}

func newMessage(mi MessageChunk, chunksSize byte) (m Message) {
	m.id = mi.id
	m.time = mi.time
	m.chunks = make([][]byte, chunksSize)
	m.chunks[mi.data[0]] = mi.data
	return
}

type ChunkResolver struct {
	chunkSize byte
	messages  map[types.ID64]Message
	resolved  map[types.ID64]struct{}
}

func newChunkResolver(chunkSize byte) *ChunkResolver {
	return &ChunkResolver{
		chunkSize: chunkSize,
		messages:  make(map[types.ID64]Message, 0),
		resolved:  make(map[types.ID64]struct{}, 0),
	}
}

func (s *ChunkResolver) IsMessageResolved(id types.ID64) bool {
	_, ok := s.resolved[id]
	return ok
}

func (s *ChunkResolver) Resolve(id types.ID64) {
	s.resolved[id] = struct{}{}
	delete(s.messages, id)
}

func (s *ChunkResolver) CreateOrUpdate(mi MessageChunk) (m Message) {
	m, ok := s.messages[mi.id]
	if ok {
		m.Update(mi.data)
	} else {
		m = newMessage(mi, s.chunkSize)
		s.messages[m.id] = m
	}

	return
}

type Processor func(m Message)

type GlobalHandler struct {
	messageCh chan MessageChunk
	chunks    *ChunkResolver
	processor func(m Message)
}

func NewGlobalHandler(conf types.NetPartition) *GlobalHandler {
	return &GlobalHandler{
		chunks:    newChunkResolver(conf.Size),
		processor: createProcessor(),
		messageCh: make(chan MessageChunk, 10),
	}
}

func (h *GlobalHandler) Message(header types.NetMessage) {
	h.messageCh <- MessageChunk{header.ID, header.Time, header.Data}
}

func (h *GlobalHandler) Start() {
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

func (h *GlobalHandler) Stop() {
	close(h.messageCh)
}

func createProcessor() Processor {
	return func(m Message) {
		d, err := types.DecryptFromChunks(m.chunks)
		if err != nil {
			log.Errorf("failed to decode message: %v \n", err)
		}

		baseMessage, err := types.UnmarshalBaseMessage(d)
		if err != nil {
			log.Error(err)
		}

		ok, err := baseMessage.Sign.CheckSignature(baseMessage.Data)
		if err != nil {
			log.Errorf("failed to check signature: %v \n", err)
		} else if !ok {
			log.Errorf("signature is invalid: %v \n", baseMessage)
		}

		log.Println(baseMessage.Type)
		log.Println(hex.EncodeToString(baseMessage.Data))
	}
}
