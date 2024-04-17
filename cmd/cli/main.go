package main

import (
	cfg "github.com/NikitosnikN/balance-api/internal/config"
	"github.com/NikitosnikN/balance-api/internal/service"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

var flags = []cli.Flag{
	&cli.StringFlag{
		Name:    "config",
		Aliases: []string{"c"},
		Usage:   "Config file path",
		Value:   "config.yaml",
	},
}

func main() {
	app := &cli.App{
		Name:  "balance-api",
		Flags: flags,
		Action: func(cCtx *cli.Context) error {
			cfgPath := cCtx.String("config")
			config, err := cfg.LoadConfigFromFile(cfgPath)

			if err != nil {
				log.Fatal("failed to load config", config)
			}

			return service.RunApplication(config)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
