package main

import (
	"context"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/transeptorlabs/betsy/internal/docker"
	"github.com/transeptorlabs/betsy/internal/mempool"
	"github.com/transeptorlabs/betsy/internal/server"
	"github.com/transeptorlabs/betsy/internal/utils"
	"github.com/transeptorlabs/betsy/logger"
	"github.com/transeptorlabs/betsy/version"
	"github.com/transeptorlabs/betsy/wallet"
	"github.com/urfave/cli/v2"
)

type NodeInfo struct {
	EthNodeUrl         string
	BundlerNodeUrl     string
	DashboardServerUrl string
	DevAccounts        []wallet.DevAccount
}

func main() {
	var err error
	var betsyWallet *wallet.Wallet

	containerManager, err := docker.NewContainerManager()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize container manager")
	}

	app := &cli.App{
		Name:    "betsy",
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
			&cli.StringFlag{
				Name:     "log.level",
				Usage:    "Enable debug mode on server",
				Aliases:  []string{"log"},
				Value:    "INFO",
				Required: false,
				Category: "Logger selection:",
			},
			&cli.BoolFlag{
				Name:     "debug",
				Usage:    "Enable debug mode on server",
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
			log.Logger, err = logger.GetLogger(cCtx.String("log.level"))
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to initialize logger")
			}

			err = printWelomeBanner()
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to load welcome banner")
			}

			log.Debug().Msgf("Running preflight checks...")

			// Check that docker is installed and pull required images
			ok := containerManager.IsDockerInstalled()
			if !ok {
				log.Fatal().Err(err).Msg("Docker needs to be installed to use Betsy!")
			}

			_, err := containerManager.PullRequiredImages(
				cCtx.Context,
				[]string{"geth", cCtx.String("bundler")},
			)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to pull required images")
			}

			return nil
		},
		After: func(cCtx *cli.Context) error {
			log.Info().Msgf("Tearing down docker containers!\n")
			_, err := containerManager.StopAndRemoveRunningContainers(cCtx.Context)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to tear down docker containers!")
			}

			containerManager.Close()
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to close container manager")
			}

			err = utils.RemoveDevWallets()
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to remove wallets")
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
			ctx, stop := signal.NotifyContext(cCtx.Context, os.Interrupt, syscall.SIGTERM)
			defer stop()

			// Wait until the eth node container is ready before starting the bundler
			// Using a channel to signal when the container is ready
			// and a channel to signal if there was an error
			readyChan := make(chan struct{})
			readyErrorChan := make(chan error)

			ctxWithReadyChan := context.WithValue(ctx, docker.EthNodeReady, readyChan)

			go func() {
				_, err := containerManager.RunContainerInTheBackground(
					ctxWithReadyChan,
					"geth",
					strconv.Itoa(cCtx.Int("eth.port")),
				)
				if err != nil {
					readyErrorChan <- err
				}
			}()

			/* Handle case where
			- eth node container fails to start
			- eth node container is ready
			- context is canceled with an interrupt signal (ctrl-c)
			*/
			select {
			case err := <-readyErrorChan:
				log.Err(err).Msg("Failed to run ETH node container")
				return nil
			case <-readyChan:
				log.Info().Msg("ETH node is ready, starting bundler and initializing dev wallet...")

				// create dev wallet with default accounts
				betsyWallet, err = wallet.NewWallet(
					ctx,
					strconv.Itoa(cCtx.Int("eth.port")),
					containerManager.CoinbaseKeystoreFile,
				)
				if err != nil {
					log.Err(err).Msg("Failed to create dev wallet")
					return nil
				}

				// Start the bundler container passing a context with the wallet details
				ctxWithBundlerDetails := context.WithValue(ctx, docker.BundlerNodeWalletDetails, betsyWallet.GetBundlerWalletDetails())

				_, err = containerManager.RunContainerInTheBackground(
					ctxWithBundlerDetails,
					cCtx.String("bundler"),
					strconv.Itoa(cCtx.Int("bundler.port")),
				)
				if err != nil {
					log.Err(err).Msgf("Failed to run %s bundler container", cCtx.String("bundler"))
					return nil
				}
			case <-ctx.Done():
				log.Info().Msg("Received signal, shutting down...")
				return nil
			}

			// create a start mempool polling
			mempool := mempool.NewUserOpMempool(
				betsyWallet.GetBundlerWalletDetails().EntryPointAddress,
				betsyWallet.GetEthClient(),
				"http://localhost:"+strconv.Itoa(cCtx.Int("bundler.port")),
			)
			go func() {
				if err := mempool.Run(); err != nil {
					log.Err(err).Msg("mempool failed")
					stop()
				}
			}()

			// create and start http server
			httpServer := server.NewHTTPServer(
				net.JoinHostPort("localhost", strconv.Itoa(cCtx.Int("http.port"))),
				cCtx.Bool("debug"),
				betsyWallet,
				mempool,
			)
			go func() {
				if err := httpServer.Run(); err != nil && err != http.ErrServerClosed {
					log.Err(err).Msg("HTTP server failed")
					stop()
				}
			}()

			accounts, err := betsyWallet.GetDevAccounts(ctx)
			if err != nil {
				log.Err(err).Msg("Failed to get dev accounts")
				stop()
			}

			prefix := "http://localhost:"
			err = printBetsyInfo(NodeInfo{
				EthNodeUrl:         prefix + strconv.Itoa(cCtx.Int("eth.port")),
				BundlerNodeUrl:     prefix + strconv.Itoa(cCtx.Int("bundler.port")),
				DashboardServerUrl: prefix + strconv.Itoa(cCtx.Int("http.port")),
				DevAccounts:        accounts,
			})
			if err != nil {
				log.Err(err).Msg("Failed print Betsy info")
				stop()
			}

			<-ctx.Done()

			// Create a context with timeout to allow the server to shut down gracefully
			shutdownCtx, cancel := context.WithTimeout(cCtx.Context, 5*time.Second)
			defer cancel()

			mempool.Stop()

			if err := httpServer.Shutdown(shutdownCtx); err != nil {
				log.Err(err).Msg("Server shutdown failed")
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

func printWelomeBanner() error {
	var tmplBannerFile = "ui/templates/banner.tmpl"
	tmpl, err := template.New("banner.tmpl").ParseFiles(tmplBannerFile)
	if err != nil {
		return err
	}
	err = tmpl.Execute(os.Stdout, nil)
	if err != nil {
		return err
	}

	return nil
}

func printBetsyInfo(nodeInfo NodeInfo) error {
	var tmplBannerFile = "ui/templates/betsy-info.tmpl"
	tmpl, err := template.New("betsy-info.tmpl").ParseFiles(tmplBannerFile)
	if err != nil {
		return err
	}

	err = tmpl.Execute(os.Stdout, nodeInfo)
	if err != nil {
		return err
	}

	return nil
}
