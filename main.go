package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

type Source struct {
	data map[string]*[]interface{}
}

type ServeOptions struct {
	host   string
	port   string
	queit  bool
	noCors bool
	tmp    bool
}

func read(path string) *Source {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	jsonMap := make(map[string]*[]interface{})

	if err = json.NewDecoder(file).Decode(&jsonMap); err != nil {
		log.Fatal(err)
	}

	return &Source{data: jsonMap}
}

func serve(s *Source, opt *ServeOptions) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		HandleBase(w, r, s, opt)
	})

	fmt.Println(color.GreenString("Havlu is on. Listening on %s:%s!\n", opt.host, opt.port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", opt.host, opt.port), nil))

}

func main() {

	app := &cli.App{
		Name:      "havlu",
		HelpName:  "havlu",
		Usage:     "Get a full take mock REST API with zero coding.",
		UsageText: "havlu [file] [global options]",
		Version:   version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Value:   "3000",
				Usage:   "port number",
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
			&cli.StringFlag{
				Name:    "quiet",
				Usage:   "suppress log messages from output",
				Aliases: []string{"q"},
			},
			&cli.BoolFlag{
				Name:    "no-cors",
				Usage:   "disable CORS headers",
				Aliases: []string{"nc"},
			},
			&cli.BoolFlag{
				Name:    "tmp-a",
				Usage:   "access to the tmp folder",
				Aliases: []string{"tmp"},
			},
		},
		Action: func(c *cli.Context) error {
			host := c.String("host")
			port := c.String("port")
			quiet := c.Bool("quiet")
			noCors := c.Bool("no-cors")
			tmp := c.Bool("tmp-a")

			file := c.Args().First()

			if file == "" {
				return cli.Exit("File name is required.", 1)
			}

			extension := filepath.Ext(file)
			var data *Source

			switch extension {
			case ".json":
				data = read(file)
			case ".hav":
				f, err := os.Open(file)

				if err != nil {
					log.Fatal(err)
				}
				result := Generate(f)

				data = &Source{data: result}

				defer f.Close()

			default:
				return cli.Exit("File extension must be .json", 1)
			}

			opt := &ServeOptions{
				host:   host,
				port:   port,
				queit:  quiet,
				noCors: noCors,
				tmp:    tmp,
			}

			serve(data, opt)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
