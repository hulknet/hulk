package net

import "github.com/kotfalya/hulk/app/types"

type MessageHandler interface {
	types.UpdateState
	Message(header types.MessageHeader, data []byte)
}

type MessageHandlerContainer struct {
	tickToHandler map[types.ID64]MessageHandler
}

func (m *MessageHandlerContainer) Message(header types.MessageHeader, data []byte) {

}

func (m *MessageHandlerContainer) UpdateState(state types.State) {

}
