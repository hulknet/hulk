package store

import (
	"github.com/asdine/storm/v3"
	"github.com/kotfalya/hulk/pkg/crypto"
)

type Net struct {
	ID       crypto.ID `storm:"id"`
	AuthorID crypto.ID `storm:"index"`
	Sign     crypto.Signature
}

func LoadNet(db storm.Node) (*Net, error) {
	n := &Net{}
	err := db.Get("host", "net", n)
	if err == storm.ErrNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return n, nil
}
