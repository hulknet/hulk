package types

const (
	BlockShortIDByteLen = 8
	TickByteLen         = 9
	IndexShift          = 1
)

type State struct {
	time   Time
	ticks  map[ID64]Tick
	blocks map[ID64]Block
}

func CreateState(time Time, blocks []Block, ticks []Tick) State {
	state := State{
		time: time,
	}

	for _, block := range blocks {
		state.blocks[block.ID.ID64()] = block
	}
	for _, tick := range ticks {
		state.ticks[tick.ID.ID64()] = tick
	}

	return state
}

func (s *State) Head() Block {
	return s.blocks[s.time.BlockShortID()]
}

func (s *State) Time() Time {
	return s.time
}

func (s *State) Block(shortID ID64) Block {
	return s.blocks[shortID]
}

func (s *State) Tick(shortID ID64) Tick {
	return s.ticks[shortID]
}

func (s *State) Validate(time Time, buckets int) bool {
	block, ok := s.blocks[time.BlockShortID()]
	if !ok || !block.Status.Active() {
		return false
	}
	for i := 0; i < buckets; i++ {
		tick, ok := s.ticks[time.TickShortID(i)]
		if !ok || !tick.Status.Active() {
			return false
		}
	}
	return true
}

type BlockStatus byte

func (b BlockStatus) Active() bool {
	return b == BlockStatusHead || b == BlockStatusTail
}

const (
	BlockStatusNew BlockStatus = iota
	BlockStatusHead
	BlockStatusTail
	BlockStatusOld
)

type Block struct {
	ID      ID256
	PID     ID256
	BitSize []uint8
	Status  BlockStatus
}

type TickStatus byte

func (t TickStatus) Active() bool {
	return t == TickStatusHead || t == TickStatusTail
}

const (
	TickStatusNew TickStatus = iota
	TickStatusHead
	TickStatusTail
	TickStatusOld
)

type Tick struct {
	ID     ID256
	Count  uint8
	Bucket uint8
	Status TickStatus
}

type Time []byte

func (t Time) Validate() bool {
	length := len(t) / BlockShortIDByteLen
	tickNum := length - 1
	mod := len(t) % BlockShortIDByteLen
	return length >= 2 && mod == tickNum
}

func (t Time) BlockShortID() (blockId ID64) {
	copy(blockId[:], t[:BlockShortIDByteLen])
	return
}

func (t Time) TickShortID(bucket int) (tickId ID64) {
	start := BlockShortIDByteLen + bucket*TickByteLen - IndexShift
	end := start + BlockShortIDByteLen + IndexShift
	copy(tickId[:], t[start:end])
	return
}

func (t Time) TickShortIDs() (tickIds []ID64) {
	for i := 0; i < (len(t) / BlockShortIDByteLen); i++ {
		tickIds = append(tickIds, t.TickShortID(i))
	}
	return
}
