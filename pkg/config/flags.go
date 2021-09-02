package config

import "github.com/urfave/cli"

var Flags []cli.Flag

func init() {
	Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "db-path, d",
			Usage:  "Path to database directory",
			EnvVar: "DB_PATH",
			Value:  "./db",
		},
		cli.StringFlag{
			Name:   "pk-path, p",
			Usage:  "Path to private key file",
			EnvVar: "PK_PATH",
		},
		cli.StringFlag{
			Name: "host-id, i",
			Usage: "ID of the host, used for db filename and based on userID and secret, " +
				"auto generated if not provided ",
			EnvVar: "HOST_ID",
		},
		cli.StringFlag{
			Name:   "secret, s",
			Usage:  "Secret string, used as salt, default value is empty rand()",
			EnvVar: "SECRET",
		},
		cli.StringFlag{
			Name:   "address, a",
			Usage:  "Listen Address",
			EnvVar: "ADDRESS",
		},
		cli.StringFlag{
			Name:   "bootstrap, b",
			Usage:  "Bootstrap Address",
			EnvVar: "BOOTSTRAP",
		},
		cli.StringFlag{
			Name:   "http-user",
			Usage:  "HTTP username for admin dashboard",
			EnvVar: "HTTP_USER",
		},
		cli.StringFlag{
			Name:   "http-pass",
			Usage:  "HTTP password for admin dashboard",
			EnvVar: "HTTP_PASSWORD",
		},
		cli.StringFlag{
			Name:   "app-data-address",
			Usage:  "Listen Address for data application",
			EnvVar: "APP_DATA_ADDRESS",
		},
		cli.StringFlag{
			Name:   "app-data-token",
			Usage:  "Secure token for data application",
			EnvVar: "APP_DATA_TOKEN",
		},
	}
}
