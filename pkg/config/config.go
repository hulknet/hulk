package config

import (
	"errors"

	"github.com/kotfalya/hulk/pkg/crypto"
)

type Option func(cfg *Config) error

type Config struct {
	Transport *TransportConfig
	DB        *DBConfig
	Crypto    *CryptoConfig
	HTTP      *HTTPConfig
	App       *AppConfigContainer
}

type HTTPConfig struct {
	Username string
	Password string
	Secret   []byte
}

type AppConfigContainer struct {
	Store *AppConfig
	Auth  *AppConfig
}

type AppConfig struct {
	PublicAddress  string
	PrivateAddress string
	Token          string
}

type CryptoConfig struct {
	Secret     []byte
	PrivateKey crypto.PrivateKey
	PublicKey  crypto.PublicKey
	UserID     crypto.ID
	HostID     crypto.ID
}

func (c *CryptoConfig) Init() error {
	if len(c.Secret) == 0 {
		return errors.New("config secret is empty")
	}
	if c.PrivateKey == (crypto.PrivateKey)(nil) {
		return errors.New("config private key is empty")
	}

	c.PublicKey = c.PrivateKey.Public()
	userID, err := c.PublicKey.ID()
	if err != nil {
		return err
	}
	c.UserID = userID
	c.HostID = userID.WithSalt(c.Secret)
	return nil
}

type DBConfig struct {
	Path string
}

type TransportConfig struct {
	BootstrapAddress string
	Address          string
}

func NewConfig() *Config {
	return &Config{
		Transport: &TransportConfig{},
		DB:        &DBConfig{},
		HTTP:      &HTTPConfig{},
		Crypto:    &CryptoConfig{},
		App: &AppConfigContainer{
			Store: &AppConfig{},
			Auth:  &AppConfig{},
		},
	}
}

func (cfg *Config) Apply(opts ...Option) error {
	for _, opt := range opts {
		if err := opt(cfg); err != nil {
			return err
		}
	}
	return nil
}
