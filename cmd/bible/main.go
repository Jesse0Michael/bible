package main

import (
	"log"
	"os"

	"github.com/Jesse0Michael/bible/internal/bible"
	"github.com/urfave/cli/v2"
)

func main() {
	characterFlags := []cli.Flag{
		&cli.StringFlag{
			Name:    "sex",
			Aliases: []string{"s"},
			Usage:   "character sex",
		},
		&cli.StringFlag{
			Name:    "parent",
			Aliases: []string{"p"},
			Usage:   "character parent (name or reference)",
		},
		&cli.StringFlag{
			Name:    "spouse",
			Aliases: []string{"sp"},
			Usage:   "character spouse (name or reference)",
		},
		&cli.StringFlag{
			Name:    "note",
			Aliases: []string{"n"},
			Usage:   "character note to be stored in info",
		},
		&cli.StringFlag{
			Name:    "reference",
			Aliases: []string{"r", "ref"},
			Usage:   "reference for note (requires note)",
		},
		&cli.StringFlag{
			Name:    "commentary",
			Aliases: []string{"c"},
			Usage:   "commentary for note (requires note)",
		},
	}
	locationFlags := []cli.Flag{
		&cli.StringFlag{
			Name:    "note",
			Aliases: []string{"n"},
			Usage:   "character note to be stored in info",
		},
		&cli.StringFlag{
			Name:    "reference",
			Aliases: []string{"r", "ref"},
			Usage:   "reference for note (requires note)",
		},
		&cli.StringFlag{
			Name:    "commentary",
			Aliases: []string{"c"},
			Usage:   "commentary for note (requires note)",
		},
	}
	app := cli.App{
		Name:    "bible",
		Version: "0.0.0",
		Commands: []*cli.Command{
			{
				Name:  "create",
				Usage: "create a new resource",
				Subcommands: []*cli.Command{
					{
						Name:      "character",
						ArgsUsage: "NAME",
						Usage:     "create a new character resource",
						Flags:     characterFlags,
						Action:    bible.CreateCharacter,
					},
					{
						Name:      "location",
						ArgsUsage: "NAME",
						Usage:     "create a new location resource",
						Flags:     locationFlags,
						Action:    bible.CreateLocation,
					},
				},
			},
			{
				Name:  "update",
				Usage: "update a resource",
				Subcommands: []*cli.Command{
					{
						Name:      "character",
						ArgsUsage: "NAME",
						Usage:     "update a character resource",
						Flags:     characterFlags,
						Action:    bible.UpdateCharacter,
					},
					{
						Name:      "location",
						ArgsUsage: "NAME",
						Usage:     "update a location resource",
						Flags:     locationFlags,
						Action:    bible.UpdateLocation,
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
