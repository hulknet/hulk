package memory

import "github.com/hulknet/hulk/app/types"

type Memory struct {
	state types.State
}

func (m *Memory) UpdateState(state types.State) {
	m.state = state
}

func (m *Memory) CreateBucket(id types.ID64) {

}
