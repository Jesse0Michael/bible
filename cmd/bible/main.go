package main

import (
	"log"
	"os"

	"github.com/Jesse0Michael/bible/internal/bible"
	"github.com/urfave/cli/v2"
)

func main() {
	getFlags := []cli.Flag{
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "output format(yaml/json)",
		},
	}
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
		&cli.StringFlag{
			Name:    "associate",
			Aliases: []string{"a"},
			Usage:   "non immediate family member whom this character associated with",
		},
		&cli.StringFlag{
			Name:    "location",
			Aliases: []string{"l"},
			Usage:   "specific location this character was at",
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
				Name:  "get",
				Usage: "get a resource",
				Subcommands: []*cli.Command{
					{
						Name:      "character",
						ArgsUsage: "NAME",
						Usage:     "get a character resource",
						Flags:     getFlags,
						Action:    bible.GetCharacter,
					},
					{
						Name:      "location",
						ArgsUsage: "NAME",
						Usage:     "get a location resource",
						Flags:     getFlags,
						Action:    bible.GetLocation,
					},
				},
			},
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
			{
				Name:  "audit",
				Usage: "audit a resource",
				Subcommands: []*cli.Command{
					{
						Name:   "character",
						Usage:  "audit all character resources",
						Action: bible.AuditCharacters,
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
