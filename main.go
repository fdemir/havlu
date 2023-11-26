package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/urfave/cli/v2"
)

type Source struct {
	data map[string]interface{}
}

func read(path string) *Source {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	jsonMap := make(map[string]interface{})

	if err = json.NewDecoder(file).Decode(&jsonMap); err != nil {
		log.Fatal(err)
	}

	return &Source{data: jsonMap}
}

func serve(s *Source, host string, port string, queit bool) {
	// todo: syncronize data access with mutex or channels
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !queit {
			log.Printf("%s %s", r.Method, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")

		path := r.URL.Path[1:]
		response := s.data[path]

		if response == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// unmarshal query params
		result := []interface{}{}
		params := r.URL.Query()

		limit, err := strconv.Atoi(params.Get("_limit"))
		if err != nil {
			// do nothing
			limit = 0
		}

		if len(params) > 0 {

			count := 0

			for _, item := range response.([]interface{}) {
				if count >= limit {
					break
				}

				shouldAdd := true
				item := item.(map[string]interface{})

				for key, value := range params {

					// if params params does start with underscore, skip
					if key[0] == '_' {
						continue
					}

					if intValue, err := strconv.Atoi(value[0]); err == nil {
						if item[key] != intValue {
							shouldAdd = false
						}
					} else if boolValue, err := strconv.ParseBool(value[0]); err == nil {
						if item[key] != boolValue {
							shouldAdd = false
						}
					} else {
						if item[key] != value[0] {
							shouldAdd = false
						}
					}
				}

				if shouldAdd {
					result = append(result, item)
				}

				count += 1
			}

		} else {
			result = response.([]interface{})
		}

		json.NewEncoder(w).Encode(result)
	})

	fmt.Printf("Havlu is on. Listening on %s:%s!\n", host, port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), nil))

}

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
			&cli.StringFlag{
				Name:    "quiet",
				Usage:   "suppress log messages from output",
				Aliases: []string{"q"},
			},
		},
		Action: func(c *cli.Context) error {
			host := c.String("host")
			port := c.String("port")
			file := c.String("file")
			quiet := c.Bool("quiet")

			data := read(file)
			serve(data, host, port, quiet)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
