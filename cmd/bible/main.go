package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Jesse0Michael/bible/internal/bible"
	"github.com/urfave/cli/v2"
)

func main() {
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
						Flags: []cli.Flag{
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
						},
						Action: bible.CreateCharacter,
					},
					{
						Name:  "location",
						Usage: "create a new location resource",
						Action: func(c *cli.Context) error {
							fmt.Println("removed task template: ", c.Args().First())
							return nil
						},
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
