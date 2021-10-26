package instruction

import (
	log "github.com/sirupsen/logrus"

	"github.com/hulknet/hulk/app/types"
)

type worker struct {
	in   chan types.BaseMessage
	inst map[types.ID64]*Instruction
}

func (w worker) run() {
	for {
		msg, ok := <-w.in
		if !ok {
			return
		}
		i := w.getOrCreateInstruction(msg)

		switch MessageTypeFromString(msg.Type) {
		case AnnounceMessageType:
			a, err := UnmarshalAnnounceMessage(msg.Data)
			if err != nil {
				log.Error(err)
			} else {
				i.Announce(a)
			}
		case InputMessageType:
			in, err := UnmarshalInputMessage(msg.Data)
			if err != nil {
				log.Error(err)
			} else {
				i.Input(in)
			}
		}
	}
}

func (w worker) getOrCreateInstruction(msg types.BaseMessage) *Instruction {
	if i, ok := w.inst[msg.ID]; ok {
		return i
	}
	return NewInstruction(msg.ID, msg.Time)
}

type Handler struct {
	wpSize byte // Worker pool size
	input  []chan types.BaseMessage
}

func NewHandler(wpSize byte) *Handler {
	h := &Handler{
		wpSize: wpSize,
		input:  make([]chan types.BaseMessage, wpSize),
	}

	return h
}

func (h *Handler) Close() {
	for _, i := range h.input {
		close(i)
	}
}

func (h *Handler) Message(msg types.BaseMessage) {
	h.input[msg.ID.Bytes()[0]%h.wpSize] <- msg
}

func (h *Handler) UpdateState(state types.State) {

}
