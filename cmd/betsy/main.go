package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/transeptorlabs/betsy/internal/docker"
	"github.com/transeptorlabs/betsy/internal/server"
	"github.com/transeptorlabs/betsy/logger"
	"github.com/transeptorlabs/betsy/version"
	"github.com/urfave/cli/v2"
)

func main() {
	var err error
	log.Logger, err = logger.GetLogger()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize logger")
	}

	containerManager, err := docker.NewContainerManager()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize container manager")
	}

	parentContext := context.Background()

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
		Usage:                "Your local Account Abstraction development environment toolkit",
		UsageText:            "Your local Account Abstraction development environment toolkit",
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
			log.Info().Msgf("Running preflight checks...")
			_, err := containerManager.PullRequiredImages(
				parentContext,
				[]string{"geth", cCtx.String("bundler")},
			)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to pull required images")
			}
			// TODO: check that docker is installed
			// TODO: check that geth is not already running
			return nil
		},
		After: func(cCtx *cli.Context) error {
			log.Info().Msgf("Tearing down dev environnement!\n")
			_, err := containerManager.StopAndRemoveRunningContainers(parentContext)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to tear down dev environnement!")
			}

			containerManager.Close()
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to close container manager")
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

			log.Error().Err(err).Msgf("WRONG: %#v\n", err)
			return nil
		},
		Action: func(cCtx *cli.Context) error {
			// Create a context that will be canceled when an interrupt signal is caught
			ctx, stop := signal.NotifyContext(parentContext, os.Interrupt, syscall.SIGTERM)
			defer stop()

			// Wait until the eth node container is ready before starting the bundler
			readyChan := make(chan struct{})
			ctxWithReadyChan := context.WithValue(ctx, docker.EthNodeReady, readyChan)

			go containerManager.RunContainerInTheBackground(
				ctxWithReadyChan,
				"geth",
				strconv.Itoa(cCtx.Int("eth.port")),
			)

			select {
			case <-readyChan:
				log.Info().Msg("ETH node is ready, starting bundler...")
			// containerManager.RunContainerInTheBackground(
			// 	cCtx.String("bundler"),
			// 	strconv.Itoa(cCtx.Int("bundler.port")),
			// )
			case <-ctx.Done():
				log.Info().Msg("Received signal, shutting down...")
				return nil
			}

			// Start the server in a goroutine
			httpServer := server.NewHTTPServer(
				net.JoinHostPort("localhost", strconv.Itoa(cCtx.Int("http.port"))),
				cCtx.Bool("debug"),
			)
			go func() {
				if err := httpServer.Run(); err != nil && err != http.ErrServerClosed {
					log.Fatal().Err(err).Msg("HTTP server failed")
				}
			}()

			<-ctx.Done()

			// Create a context with timeout to allow the server to shut down gracefully
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := httpServer.Shutdown(shutdownCtx); err != nil {
				log.Fatal().Err(err).Msg("Server shutdown failed")
			} else {
				log.Info().Msg("Server shutdown completed")
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("Failed to run app")
	}
}
