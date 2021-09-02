package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/kotfalya/hulk/pkg/crypto"
	"github.com/urfave/cli"
)

func CliOption(c *cli.Context) Option {
	return func(cfg *Config) error {
		if c.String("pk-path") != "" {
			pk, err := pkLoad(c)
			if err != nil {
				return err
			}
			cfg.Crypto.PrivateKey = pk
		}
		cfg.Crypto.Secret = []byte(c.String("secret"))
		cfg.Transport.Address = c.String("address")
		cfg.Transport.BootstrapAddress = c.String("bootstrap")
		cfg.DB.Path = c.String("db-path")
		cfg.HTTP.Username = c.String("http-user")
		cfg.HTTP.Password = c.String("http-pass")
		cfg.App.Store.PublicAddress = c.String("app-data-address")
		cfg.App.Store.Token = c.String("app-data-token")
		return nil
	}
}

func CryptoOption() Option {
	return func(cfg *Config) error {
		return cfg.Crypto.Init()
	}
}

func DBLoadOption(c *cli.Context) Option {
	return func(cfg *Config) error {
		path := filepath.Join(c.String("db-path"), DBPathPrefix+c.String("host-id"))
		return loadConfig(cfg, path)
	}
}
func DBSaveOption() Option {
	return func(cfg *Config) error {
		return saveConfig(cfg)
	}
}

func DBDirCheckOption(c *cli.Context) Option {
	return func(cfg *Config) error {
		if st, err := os.Stat(c.String("db-path")); os.IsNotExist(err) || !st.IsDir() {
			return errors.New(fmt.Sprintf("dp-path must be a directory, %s", c.String("db-path")))
		}
		return nil
	}
}

func pkLoad(c *cli.Context) (crypto.PrivateKey, error) {
	pkData, err := ioutil.ReadFile(c.String("pk-path"))
	if err != nil {
		return nil, err
	}
	return crypto.UnmarshalPrivateKey(pkData)
}
