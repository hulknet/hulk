package net

import (
	"encoding/hex"
	"fmt"

	"github.com/vmihailenco/msgpack/v5"

	"github.com/hulknet/hulk/app/types"
)

type MessageChunk struct {
	id   types.ID64
	data []byte
}

type Message struct {
	id     types.ID64
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

type BucketHandler struct {
	messageCh chan MessageChunk
	chunks    *ChunkResolver
	processor func(m Message)
}

func NewBucketHandler(chunkSize byte) *BucketHandler {
	return &BucketHandler{
		chunks:    newChunkResolver(chunkSize),
		processor: createProcessor(),
		messageCh: make(chan MessageChunk, 10),
	}
}

func (h *BucketHandler) Message(header types.MessageHeader, data []byte) {
	h.messageCh <- MessageChunk{header.ID, data}
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
		d, err := types.DecryptFromChunks(m.chunks)
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
