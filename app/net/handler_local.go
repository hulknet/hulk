package net

import (
	"fmt"

	"github.com/vmihailenco/msgpack/v5"

	"github.com/hulknet/hulk/app/types"
)

type LocalMessage struct {
	id   types.ID64
	data []byte
}

type LocalHandler struct {
	messageCh chan LocalMessage
	processor func(m LocalMessage)
	resolved  map[types.ID64]struct{}
}

func NewLocalHandler() *LocalHandler {
	return &LocalHandler{
		processor: createLocalProcessor(),
		messageCh: make(chan LocalMessage, 10),
		resolved:  make(map[types.ID64]struct{}),
	}
}

func (h *LocalHandler) Start() {
	for {
		select {
		case mi, ok := <-h.messageCh:
			if !ok {
				return
			}
			if _, resolved := h.resolved[mi.id]; resolved {
				continue
			}
			go h.processor(mi)
			h.resolved[mi.id] = struct{}{}
		}
	}
}

func (h *LocalHandler) Stop() {
	close(h.messageCh)
}

func (h *LocalHandler) Message(header types.NetMessage) {
	h.messageCh <- LocalMessage{header.ID, header.Data}
}

type LocalProcessor func(m LocalMessage)

func createLocalProcessor() LocalProcessor {
	return func(m LocalMessage) {
		var baseMessage types.BaseMessage
		if err := msgpack.Unmarshal(m.data, &baseMessage); err != nil {
			fmt.Printf("failed to unmarshal BaseMessage: %v \n", err)
		}

		ok, err := types.CheckSignature(baseMessage.Data, baseMessage.Sign)
		if err != nil {
			fmt.Printf("failed to check signature: %v \n", err)
		} else if !ok {
			fmt.Printf("signature is invalid: %v \n", baseMessage)
		}

		fmt.Println(baseMessage.Type)
	}
}
