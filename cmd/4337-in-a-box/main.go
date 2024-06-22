package main

import (
	"fmt"
	"log"
	"os"

	"github.com/transeptorlabs/4337-in-a-box/version"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "4337 In A Box",
		Version: version.Version,
		Authors: []*cli.Author{
			{
				Name:  "Transeptor Labs",
				Email: "transeptorhq@gmail.com",
			},
		},
		Copyright: "(c) 2024 Transeptor Labs",
		Usage:     "Manage your 4337 development environment",
		UsageText: "Manage your 4337 development environment",
		Commands: []*cli.Command{
			{
				Name:        "doo",
				Aliases:     []string{"do"},
				Category:    "motion",
				Usage:       "do the doo",
				UsageText:   "doo - does the dooing",
				Description: "no really, there is a lot of dooing to be done",
				ArgsUsage:   "[arrgh]",
				Flags: []cli.Flag{
					&cli.BoolFlag{Name: "forever", Aliases: []string{"forevvarr"}},
				},
				Subcommands: []*cli.Command{
					{
						Name:   "wop",
						Action: wopAction,
					},
				},
				SkipFlagParsing: false,
				HideHelp:        false,
				Hidden:          false,
				HelpName:        "doo!",
				BashComplete: func(cCtx *cli.Context) {
					fmt.Fprintf(cCtx.App.Writer, "--better\n")
				},
				Before: func(cCtx *cli.Context) error {
					fmt.Fprintf(cCtx.App.Writer, "brace for impact\n")
					return nil
				},
				After: func(cCtx *cli.Context) error {
					fmt.Fprintf(cCtx.App.Writer, "did we lose anyone?\n")
					return nil
				},
				Action: func(cCtx *cli.Context) error {
					cCtx.Command.FullName()
					cCtx.Command.HasName("wop")
					cCtx.Command.Names()
					cCtx.Command.VisibleFlags()
					fmt.Fprintf(cCtx.App.Writer, "dodododododoodododddooooododododooo\n")
					if cCtx.Bool("forever") {
						cCtx.Command.Run(cCtx)
					}
					return nil
				},
				OnUsageError: func(cCtx *cli.Context, err error, isSubcommand bool) error {
					fmt.Fprintf(cCtx.App.Writer, "for shame\n")
					return err
				},
			},
		},
		EnableBashCompletion: true,
		HideVersion:          false,
		HideHelp:             false,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "lang",
				Value: "english",
				Usage: "language for the greeting",
			},
		},
		Before: func(cCtx *cli.Context) error {
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
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func wopAction(cCtx *cli.Context) error {
	fmt.Fprintf(cCtx.App.Writer, ":wave: over here, eh\n")
	return nil
}
