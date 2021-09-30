package net

import (
	"encoding/json"
	"fmt"

	"github.com/kotfalya/hulk/app/types"
)

type NodeMessage struct {
	id   types.ID64
	data []byte
}

type NodeHandler struct {
	messageCh chan NodeMessage
	processor func(m NodeMessage)
	resolved  map[types.ID64]struct{}
}

func NewNodeHandler() *NodeHandler {
	return &NodeHandler{
		processor: createNodeProcessor(),
		messageCh: make(chan NodeMessage, 10),
		resolved:  make(map[types.ID64]struct{}),
	}
}

func (h *NodeHandler) Start() {
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
		}
	}
}

func (h *NodeHandler) Stop() {
	close(h.messageCh)
}

func (h *NodeHandler) Message(header types.MessageHeader, data []byte) {
	h.messageCh <- NodeMessage{header.ID, data}
}

type NodeProcessor func(m NodeMessage)

func createNodeProcessor() NodeProcessor {
	return func(m NodeMessage) {
		var baseMessage types.BaseMessage
		if err := json.Unmarshal(m.data, &baseMessage); err != nil {
			fmt.Printf("failed to unmarshal BaseMessage: %v \n", err)
		}

		ok, err := types.CheckStringSignature(*baseMessage.Data, baseMessage.Sign)
		if err != nil {
			fmt.Printf("failed to check signature: %v \n", err)
		} else if !ok {
			fmt.Printf("signature is invalid: %v \n", baseMessage)
		}

		fmt.Println(baseMessage.Type)
	}
}
