package config

import (
	"crypto/rand"
	"crypto/rsa"

	"github.com/kotfalya/hulk/pkg/crypto"
	"github.com/kotfalya/hulk/pkg/utils"
)

var DefaultConfig = []Option{
	DefaultSecret,
	DefaultTlsConfig,
	DefaultListenerAddress,
	DefaultHTTPConfig,
	DefaultDataAppConfig,
	DefaultAuthAppConfig,
}

func DefaultListenerAddress(cfg *Config) error {
	if cfg.Transport.Address == "" {
		cfg.Transport.Address = "0.0.0.0:3456"
	}
	return nil
}

func DefaultDataAppConfig(cfg *Config) error {
	if cfg.App.Store.PublicAddress == "" {
		cfg.App.Store.PublicAddress = "0.0.0.0:3457"
	}
	if cfg.App.Store.PrivateAddress == "" {
		cfg.App.Store.PrivateAddress = "0.0.0.0:7543"
	}
	if cfg.App.Store.Token == "" {
		cfg.App.Store.Token = "store_app_token"
	}
	return nil
}

func DefaultAuthAppConfig(cfg *Config) error {
	if cfg.App.Store.PublicAddress == "" {
		cfg.App.Store.PublicAddress = "0.0.0.0:3458"
	}
	if cfg.App.Store.PrivateAddress == "" {
		cfg.App.Store.PrivateAddress = "0.0.0.0:8548"
	}
	if cfg.App.Store.Token == "" {
		cfg.App.Store.Token = "auth_app_token"
	}
	return nil
}

func DefaultSecret(cfg *Config) error {
	if len(cfg.Crypto.Secret) == 0 {
		cfg.Crypto.Secret = []byte(string(rune(utils.Random())))
	}
	return nil
}

func DefaultTlsConfig(cfg *Config) error {
	if cfg.Crypto.PrivateKey == (crypto.PrivateKey)(nil) {
		key, err := rsa.GenerateKey(rand.Reader, 1024)
		if err != nil {
			return err
		}
		cfg.Crypto.PrivateKey = (*crypto.RSAPrivateKey)(key)
	}
	return nil
}

func DefaultHTTPConfig(cfg *Config) error {
	if cfg.HTTP.Username == "" {
		cfg.HTTP.Username = "admin"
	}
	if cfg.HTTP.Password == "" {
		cfg.HTTP.Password = "admin"
	}
	if len(cfg.HTTP.Secret) == 0 {
		cfg.HTTP.Secret = []byte(string(rune(utils.Random())))
	}
	return nil
}
