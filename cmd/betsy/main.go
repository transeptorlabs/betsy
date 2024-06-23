package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/transeptorlabs/betsy/internal/eth"
	"github.com/transeptorlabs/betsy/internal/server"
	"github.com/transeptorlabs/betsy/version"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "Betsy",
		Version: version.Version,
		Authors: []*cli.Author{
			{
				Name:  "Transeptor Labs",
				Email: "transeptorhq@gmail.com",
			},
		},
		Copyright:            "(c) 2024 Transeptor Labs",
		Usage:                "Your local 4337 development environment",
		UsageText:            "Your local 4337 development environment",
		EnableBashCompletion: true,
		HideVersion:          false,
		HideHelp:             false,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:     "debug",
				Usage:    "Enable debug server",
				Aliases:  []string{"d"},
				Value:    false,
				Required: false,
				Category: "Http server selection:",
			},
			&cli.UintFlag{
				Name:     "http.port",
				Usage:    "HTTP server listening port",
				Required: false,
				Value:    8080,
				Category: "Http server selection:",
			},
			&cli.UintFlag{
				Name:     "eth.port",
				Usage:    "ETH client network port",
				Required: false,
				Value:    8545,
				Category: "ETH client selection:",
			},
			&cli.StringFlag{
				Name:     "bundler",
				Usage:    "ERC 4337 bundler",
				Required: false,
				Value:    "transeptor",
				Category: "ERC 4337 bundler selection:",
			},
			&cli.UintFlag{
				Name:     "bundler.port",
				Usage:    "ERC 4337 bundler listening port",
				Required: false,
				Value:    4337,
				Category: "ERC 4337 bundler selection:",
			},
		},
		Before: func(cCtx *cli.Context) error {
			// TODO: check that docker is installed
			// TODO: Check in geth image is available if not pull it
			// TODO: check that geth is not already running
			fmt.Fprintf(cCtx.App.Writer, "HEEEERE GOES\n")
			return nil
		},
		After: func(cCtx *cli.Context) error {
			fmt.Fprintf(cCtx.App.Writer, "Phew!\n")
			return nil
		},
		CommandNotFound: func(cCtx *cli.Context, command string) {
			fmt.Fprintf(cCtx.App.Writer, "Thar be no %q here.\n", command)
		},
		OnUsageError: func(cCtx *cli.Context, err error, isSubcommand bool) error {
			if isSubcommand {
				return err
			}

			fmt.Fprintf(cCtx.App.Writer, "WRONG: %#v\n", err)
			return nil
		},
		Action: func(cCtx *cli.Context) error {
			eth.StartGethNode()

			httpServer := server.NewHTTPServer(
				net.JoinHostPort("localhost", strconv.Itoa(cCtx.Int("http.port"))),
				cCtx.Bool("debug"),
			)

			httpServer.Run()
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
