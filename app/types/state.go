package types

import "encoding/hex"

const (
	TickByteLen = 9
	IndexShift  = 1
)

type UpdateState interface {
	UpdateState(state State)
}

type State struct {
	time  Time
	ticks map[ID64]Tick
	block Block
	peer  Peer
	key   *ECKey
	token Token
}

func CreateState(time Time, block Block, ticks []Tick, peer Peer, key *ECKey, token Token) State {
	state := State{
		time:  time,
		ticks: make(map[ID64]Tick, len(ticks)),
		block: block,
		peer:  peer,
		key:   key,
		token: token,
	}

	for _, tick := range ticks {
		state.ticks[tick.ID] = tick
	}

	return state
}

func (s State) Time() Time {
	return s.time
}

func (s State) Token() Token {
	return s.token
}

func (s State) Key() *ECKey {
	return s.key
}

func (s State) Peer() Peer {
	return s.peer
}

func (s State) Block() Block {
	return s.block
}

func (s State) ID() ID64 {
	return s.block.ID
}

func (s State) FindTick(id ID64) (tick Tick, ok bool) {
	tick, ok = s.ticks[id]
	return
}

func (s State) Ticks(onlyActive bool) (ticks []Tick) {
	for _, tick := range s.ticks {
		if onlyActive && !tick.Status.IsActive() {
			continue
		}
		ticks = append(ticks, tick)
	}
	return
}

func (s State) ValidateTime(from ID64, time Time) bool {
	cpl := Cpl(from.Bytes(), s.peer.Pub.ID().Bytes())
	for i, size := range s.block.BitSize {
		tick, ok := s.ticks[time.TickID(i)]
		if !ok || !tick.Status.IsActive() {
			return false
		}
		if cpl <= int(size) {
			break
		}
		cpl -= int(size)
	}
	return true
}

func (s State) LastCommonTick(from ID64, time Time) (tick Tick, ok bool) {

	//cpl := Cpl(from.Bytes(), s.blockIdToStatePeer[block.ID].Peer.Pub.ID().Bytes())
	return
}

type BlockStatus byte

func (b BlockStatus) IsActive() bool {
	return b == BlockStatusHead || b == BlockStatusTail
}

func (b BlockStatus) String() string {
	switch b {
	case BlockStatusNew:
		return "new"
	case BlockStatusHead:
		return "head"
	case BlockStatusTail:
		return "tail"
	case BlockStatusOld:
		return "old"
	default:
		return "unknown"
	}
}

const (
	BlockStatusNew BlockStatus = iota
	BlockStatusHead
	BlockStatusTail
	BlockStatusOld
)

type Block struct {
	ID      ID64
	ID256   ID256
	BitSize []uint8
	Status  BlockStatus
}

type TickStatus byte

func (t TickStatus) IsActive() bool {
	return t == TickStatusHead || t == TickStatusTail
}

const (
	TickStatusNew TickStatus = iota
	TickStatusHead
	TickStatusTail
	TickStatusOld
)

type Tick struct {
	ID            ID64
	ID256         ID256
	IsNode        bool
	Count         uint8
	BitSize       uint8
	BitSizePrefix uint8
	Status        TickStatus
}

type Time []byte

func (t Time) Bytes() []byte {
	return t[:]
}

func (t Time) Hex() string {
	return hex.EncodeToString(t[:])
}

func (t Time) Validate() bool {
	length := len(t) / TickByteLen
	mod := len(t) % TickByteLen
	return length >= 1 && mod == 0
}

func (t Time) CommonPrefix(t1 Time) (tickIDs []ID64) {
	for i, tickID := range t.TickIDs(false) {
		if tickID != t1.TickID(i) {
			return
		}
		tickIDs = append(tickIDs, tickID)
	}
	return
}

func (t Time) TickID(bucket int) (tickId ID64) {
	start := bucket * TickByteLen
	end := start + TickByteLen + IndexShift
	copy(tickId[:], t[start:end])
	return
}

func (t Time) TickIDs(reverse bool) (tickIds []ID64) {
	if reverse {
		for i := (len(t) / TickByteLen) - 1; i >= 0; i-- {
			tickIds = append(tickIds, t.TickID(i))
		}
	} else {
		for i := 0; i < (len(t) / TickByteLen); i++ {
			tickIds = append(tickIds, t.TickID(i))
		}
	}
	return
}
