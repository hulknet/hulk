package instruction

import (
	"io/ioutil"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/hulknet/hulk/app/types"
)

const (
	ErrGetMessageType     = "failed to get MessageType from request header"
	ErrUnknownMessageType = "unknown MessageType"
)

type Container struct {
	mux            sync.RWMutex
	blockToHandler map[types.ID64]*Handler
}

func (c *Container) Handler(blockId types.ID64) *Handler {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.blockToHandler[blockId]
}

func NewHandlerContainer() *Container {
	return &Container{
		blockToHandler: make(map[types.ID64]*Handler),
	}
}

func (c *Container) SetState(state types.State) {
	h, ok := c.blockToHandler[state.Block().ID]
	if !ok {
		c.blockToHandler[state.Block().ID] = NewHandler(10)
	} else {
		h.UpdateState(state)
	}
}

type Receiver struct {
	cont *Container
}

func NewReceiver(cont *Container) *Receiver {
	return &Receiver{
		cont: cont,
	}
}

func (rh *Receiver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// todo: check internal auth token

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		log.Error(err)
		return
	}

	baseMessage, err := types.UnmarshalBaseMessage(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error(err)
		return
	}

	h := rh.cont.Handler(baseMessage.Time.BlockID())
	h.Message(baseMessage)
}
