package net

import "github.com/hulknet/hulk/app/types"

type HandlerLocalTree struct {
	level   byte
	tick    types.Tick
	locals  map[byte]*HandlerLocalTree
	globals map[byte]*HandlerGlobalTree
}

type HandlerGlobalTree struct {
	level    byte
	inc      byte
	req      chan types.NetMessage
	children map[byte]*HandlerGlobalTree
}

type MessageHandler interface {
	Message(msg types.NetMessage)
	Start()
	Stop()
}

type MessageHandlerContainer struct {
	state   types.State
	handler map[types.ID64]MessageHandler
}

func NewMessageHandlerContainer(state types.State) *MessageHandlerContainer {
	m := &MessageHandlerContainer{
		state:   state,
		handler: make(map[types.ID64]MessageHandler),
	}
	m.handler[m.state.TimeToHandlerID(state.Now())] = NewLocalHandler()
	go m.handler[m.state.TimeToHandlerID(state.Now())].Start()
	return m
}

func (m *MessageHandlerContainer) Message(msg types.NetMessage) {
	index := m.state.TimeToHandlerID(msg.Time)
	_, ok := m.handler[index]
	if !ok {
		m.handler[index] = NewGlobalHandler(m.state.NetPartition())
		go m.handler[index].Start()
	}
	m.handler[index].Message(msg)
}

func (m *MessageHandlerContainer) UpdateState(state types.State) {

}
