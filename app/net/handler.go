package net

import "github.com/kotfalya/hulk/app/types"

type MessageHandler interface {
	Message(header types.MessageHeader, data []byte)
	Start()
	Stop()
}

type MessageHandlerContainer map[types.ID64]MessageHandler

func NewMessageHandlerContainer(state types.State) MessageHandlerContainer {
	ticks := state.Ticks(true)
	m := make(map[types.ID64]MessageHandler, len(ticks))
	for _, tick := range ticks {
		if tick.IsNode {
			m[tick.ID] = NewNodeHandler()
		} else {
			m[tick.ID] = NewBucketHandler()
		}
		go m[tick.ID].Start()
	}

	return m
}

func (m MessageHandlerContainer) Message(header types.MessageHeader, data []byte) {
	tickIds := header.Time.TickIDs(true)
	for _, tickId := range tickIds {
		handler, ok := m[tickId]
		if !ok {
			continue
		}
		handler.Message(header, data)
	}
}

func (m MessageHandlerContainer) UpdateState(state types.State) {

}
