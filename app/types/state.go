package types

import (
	"bytes"
	"errors"
)

type UpdateState interface {
	UpdateState(state State)
}

//State should be immutable and represent single time state todo
type State struct {
	head  Time
	ticks map[ID64]Tick
	block Block

	// should be in another place
	peer    Peer
	key     *ECKey
	token   Token
	netPart NetPartition
}

func CreateState(time Time, block Block, ticks []Tick, peer Peer, key *ECKey, token Token) State {
	state := State{
		head:    time,
		ticks:   make(map[ID64]Tick, len(ticks)),
		block:   block,
		peer:    peer,
		key:     key,
		token:   token,
		netPart: NetPartition{Size: 2},
	}

	for _, tick := range ticks {
		state.ticks[tick.ID] = tick
	}

	return state
}

func (s State) Now() Time {
	return s.head
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
	return s.Block().ID
}

func (s State) NetPartition() NetPartition {
	return s.netPart
}

func (s State) CommonLevel(target ID64) byte {
	cpl := Cpl(s.peer.Pub.ID().Bytes(), target.Bytes())
	for i, size := range s.block.BitSize {
		if cpl <= int(size) {
			return byte(i)
		}
	}
	return byte(len(s.block.BitSize)) // if distribution works well, it shouldn't happen.
}

func (s State) IsLocalTime(time Time) bool {
	return s.CommonLevel(time.Addr()) == byte(len(s.block.BitSize))
}

func (s State) TimeToHandlerID(time Time) ID64 {
	level := s.CommonLevel(time.Addr())
	data := bytes.NewBuffer([]byte{})
	for _, t := range time.Ticks() {
		if t.Level > level {
			data.Write([]byte{t.Inc})
		} else {
			data.Write(t.ID.Bytes())
		}
	}
	return GenerateID256FromSource(data.Bytes()).ID64()
}

func (s State) ValidateTime(time Time) bool {
	level := s.CommonLevel(time.Addr())
	for _, t := range time.Ticks() {
		if t.Level > level {
			break
		}
		tick, ok := s.ticks[t.ID]
		if !ok || !tick.Status.IsActive() {
			return false
		}
	}

	return true
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

const (
	TickStatusNew TickStatus = iota
	TickStatusHead
	TickStatusTail
	TickStatusOld
)

type TickStatus byte

func (t TickStatus) IsActive() bool {
	return t == TickStatusHead || t == TickStatusTail
}

const (
	IDByteLength   = 8
	TickIDByteLen  = 8
	TickIncByteLen = 1
	TickByteLen    = 9
	IndexShift     = 1
)

type Tick struct {
	ID     ID64
	Inc    byte
	Level  byte
	Status TickStatus
}

type TickGlobal struct {
	ID    ID64
	Inc   byte
	Level byte
}

type Time struct {
	blockId ID64
	addr    ID64
	ticks   []TickGlobal
}

func (t Time) Head() TickGlobal {
	return t.ticks[len(t.ticks)-1]
}

func (t Time) BlockID() ID64 {
	return t.blockId
}

func (t Time) Addr() ID64 {
	return t.addr
}

func (t Time) Ticks() []TickGlobal {
	return t.ticks
}

func (t Time) IncHash() string {
	hash := make([]byte, len(t.ticks))
	for i, t := range t.ticks {
		hash[i] = t.Inc
	}
	return string(hash)
}

func (t Time) Encode() []byte {
	data := bytes.NewBuffer(t.addr.Bytes())
	data.Write(t.blockId.Bytes())
	for _, t := range t.ticks {
		data.Write(t.ID.Bytes())
		data.Write([]byte{t.Inc})
	}
	return data.Bytes()
}

// DecodeTime src[--- addr[8] ---,--- blockId[8] ---,--- (tickId[8],tickInc[1])[9] ---, ... ]
func DecodeTime(src []byte) (time Time, err error) {
	if !validateTimeBytes(src) {
		err = errors.New("invalid time source")
		return
	}

	time.addr = parseAddr(src)
	time.blockId = parseBlockId(src)
	time.ticks = parseTicks(src)

	return
}

func parseBlockId(src []byte) (id ID64) {
	_ = copy(id[:], src[IDByteLength:IDByteLength*2])
	return
}

func parseAddr(src []byte) (id ID64) {
	_ = copy(id[:], src[:IDByteLength])
	return
}

func parseTicks(src []byte) []TickGlobal {
	start := IDByteLength * 2
	count := (len(src) - start) / TickByteLen
	ticks := make([]TickGlobal, count)
	for i := 0; i < count; i++ {
		var tickId ID64
		_ = copy(tickId[:], src[start:start+TickIDByteLen])

		ticks[i] = TickGlobal{
			ID:    tickId,
			Level: byte(i),
			Inc:   src[start+TickByteLen-1],
		}

		start += TickByteLen
	}

	return ticks
}

func validateTimeBytes(src []byte) bool {
	length := len(src) - IDByteLength - IDByteLength // length minus BlockID and Addr
	count := length / TickByteLen
	mod := length % TickByteLen
	return count >= 1 && mod == 0
}

//
//type TimeByte []byte
//type TickInc byte
//type TimeInc []TickInc
//
//func (ti TimeInc) String() string {
//	return string(ti)
//}
//
//func (t TimeByte) Bytes() []byte {
//	return t[:]
//}
//
//func (t TimeByte) Hex() string {
//	return hex.EncodeToString(t[:])
//}
//
//func (t TimeByte) Validate(level byte) bool {
//	length := len(t) / TickByteLen
//	mod := len(t) % TickByteLen
//	return length == int(level) && mod == 0
//}
//
//func (t TimeByte) CommonPrefix(t1 TimeByte) (tickIDs []ID64) {
//	for i, tickID := range t.TickIDs(false) {
//		if tickID != t1.TickID(byte(i)) {
//			return
//		}
//		tickIDs = append(tickIDs, tickID)
//	}
//	return
//}
//
//func (t TimeByte) Level() byte {
//	return byte(len(t) / TickByteLen)
//}
//
//func (t TimeByte) TickID(level byte) (tickId ID64) {
//	start := level * TickByteLen
//	end := start + TickByteLen - IndexShift
//	copy(tickId[:], t[start:end])
//	return
//}
//
//func (t TimeByte) TickInc(level byte) TickInc {
//	return TickInc(t[(level+1)*TickByteLen-IndexShift])
//}
//
//func (t TimeByte) TimeInc() (inc TimeInc) {
//	for i := byte(0); i < t.Level(); i++ {
//		inc = append(inc, t.TickInc(i))
//	}
//	return
//}
//
//func (t TimeByte) LastTickInc() TickInc {
//	return t.TickInc(t.Level())
//}
//
//func (t TimeByte) LastTickID() ID64 {
//	return t.TickID(t.Level())
//}
//
//func (t TimeByte) TickIDs(reverse bool) (tickIds []ID64) {
//	if reverse {
//		for i := (len(t) / TickByteLen) - 1; i >= 0; i-- {
//			tickIds = append(tickIds, t.TickID(byte(i)))
//		}
//	} else {
//		for i := 0; i < (len(t) / TickByteLen); i++ {
//			tickIds = append(tickIds, t.TickID(byte(i)))
//		}
//	}
//	return
//}
