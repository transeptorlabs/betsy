package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/transeptorlabs/betsy/internal/docker"
	"github.com/transeptorlabs/betsy/internal/server"
	"github.com/transeptorlabs/betsy/version"
	"github.com/urfave/cli/v2"
)

func main() {
	containerManager := docker.NewContainerManager()
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
			fmt.Fprintf(cCtx.App.Writer, "Running preflight checks...\n")
			containerManager.PullRequiredImages(
				[]string{"geth", cCtx.String("bundler")},
			)
			// TODO: check that docker is installed
			// TODO: check that geth is not already running
			return nil
		},
		After: func(cCtx *cli.Context) error {
			fmt.Fprintf(cCtx.App.Writer, "Tearing down dev environnement!\n")
			ok, err := containerManager.StopRunningContainers()
			if err != nil {
				panic(err)
			}
			if ok {
				fmt.Fprintf(cCtx.App.Writer, "All containers stopped!\n")
			}

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

			// Run geth in the background
			containerManager.RunContainerInTheBackground(
				"geth",
				strconv.Itoa(cCtx.Int("eth.port")),
			)

			// TODO: Run the ERC 4337 bundler in the background

			// Run the HTTP server
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
