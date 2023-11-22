package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:      "havlu",
		HelpName:  "havlu",
		Usage:     "Get a full take mock REST API with zero coding.",
		Version:   "0.0.1",
		UsageText: "havlu [global options]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Value:   "3000",
				Usage:   "port number",
			},
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Value:   "db.json",
				Usage:   "file name",
			},
			&cli.StringFlag{
				Name:    "host",
				Aliases: []string{"H"},
				Value:   "localhost",
				Usage:   "host name",
			},
			&cli.StringFlag{
				Name:    "read-only",
				Usage:   "allow only GET requests",
				Aliases: []string{"ro"},
			},
			&cli.StringFlag{
				Name:    "delay",
				Usage:   "add delay to responses (ms)",
				Value:   "0",
				Aliases: []string{"d"},
			},
		},
		Action: func(c *cli.Context) error {
			port := c.String("port")
			file := c.String("file")

			fmt.Printf("port: %s, file: %s", port, file)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
