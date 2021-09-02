package main

import (
	"os"

	"github.com/kotfalya/hulk/pkg/api"
	"github.com/kotfalya/hulk/pkg/config"
	"github.com/kotfalya/hulk/pkg/host"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	//log.SetFormatter(new(log.JSONFormatter))
}

func main() {
	app := cli.NewApp()
	app.Flags = config.Flags

	app.Action = func(c *cli.Context) error {
		cfg := config.NewConfig()

		var options []config.Option
		options = append(options, config.DBDirCheckOption(c))
		if c.String("host-id") == "" {
			options = append(options, config.CliOption(c))
			options = append(options, config.DefaultConfig...)
			options = append(options, config.CryptoOption())
			options = append(options, config.DBSaveOption())
		} else {
			options = append(options, config.DBLoadOption(c))
		}

		if err := cfg.Apply(options...); err != nil {
			return err
		}
		db, err := config.GetDatabase(cfg)
		if err != nil {
			return err
		}
		h := host.NewHost(cfg, db)
		if err := h.Load(); err != nil {
			return err
		}

		rest := api.NewRestServer(h, cfg.HTTP)
		return rest.Listen()
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
