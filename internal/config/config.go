package config

import "github.com/urfave/cli/v2"

// BetsyFlags is the list of flags that can be passed to the Betsy command line
var BetsyFlags = []cli.Flag{
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
}
