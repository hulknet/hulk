package config

import (
	"encoding/gob"
	"path/filepath"

	"github.com/asdine/storm/v3"
	gobCodec "github.com/asdine/storm/v3/codec/gob"
	"github.com/kotfalya/hulk/pkg/crypto"
)

const DBPathPrefix = "host."

func GetDatabase(cfg *Config) (*storm.DB, error) {
	path := filepath.Join(cfg.DB.Path, DBPathPrefix+cfg.Crypto.HostID.Hex())
	return storm.Open(path, storm.Codec(gobCodec.Codec))
}

func saveConfig(cfg *Config) error {
	db, err := GetDatabase(cfg)
	if err != nil {
		return err
	}
	defer db.Close()
	initGob()

	return db.Set("host", "config", cfg)
}

func loadConfig(cfg *Config, path string) error {
	db, err := storm.Open(path, storm.Codec(gobCodec.Codec))
	if err != nil {
		return err
	}
	defer db.Close()
	initGob()
	return db.WithCodec(gobCodec.Codec).Get("host", "config", cfg)
}

func initGob() {
	gob.Register(&crypto.RSAPrivateKey{})
	gob.Register(&crypto.RSAPublicKey{})
}
