package net

import (
	"github.com/kotfalya/hulk/app/routing"
	"github.com/kotfalya/hulk/app/types"
)

type Container struct {
	blockToNet map[types.ID64]*Net
}

func NewNetContainer() *Container {
	return &Container{
		make(map[types.ID64]*Net),
	}
}

func (c *Container) SetState(state types.State) {
	net, ok := c.blockToNet[state.Block().ID]
	if !ok {
		c.blockToNet[state.Block().ID] = NewNet(state)
	} else {
		net.UpdateState(state)
	}
}

func (c *Container) Net(id types.ID64) (net *Net, ok bool) {
	net, ok = c.blockToNet[id]
	return
}

type Net struct {
	state     types.State
	table     *routing.Table
	handler   MessageHandlerContainer
	allowList AllowList
}

func NewNet(state types.State) *Net {
	return &Net{
		state:     state,
		table:     routing.NewTable(state.Block(), state.Peer()),
		handler:   NewMessageHandlerContainer(state),
		allowList: createAllowList(state.Peer()),
	}
}

func (n *Net) UpdateState(state types.State) {
	n.handler.UpdateState(state)
	n.state = state
}

func (n *Net) State() types.State {
	return n.state
}

func (n *Net) IsActive() bool {
	return n.state.Block().Status.IsActive()
}

func (n *Net) Table() *routing.Table {
	return n.table
}

func (n *Net) AllowList() AllowList {
	return n.allowList
}

func (n *Net) HandleMessage(header types.MessageHeader, data []byte) {
	n.handler.Message(header, data)
}

func (n *Net) ProxyMessage(header types.MessageHeader, data []byte) {

}
