package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

type Source struct {
	data map[string]interface{}
}

type ServeOptions struct {
	host   string
	port   string
	queit  bool
	noCors bool
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

func serve(s *Source, opt *ServeOptions) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if !opt.queit {
			log.Printf("%s %s", r.Method, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Methods", "GET")

		if !opt.noCors {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		path := r.URL.Path[1:]
		response := s.data[path]

		if response == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if r.Method == "GET" {

			result := []interface{}{}
			params := r.URL.Query()
			limit, _ := strconv.Atoi(params.Get("_limit"))

			if len(params) > 0 {
				count := 0

				for _, item := range response.([]interface{}) {
					hasLimitReached := limit > 0 && count >= limit

					if hasLimitReached {
						break
					}

					shouldAdd := true
					item := item.(map[string]interface{})

					for key, value := range params {

						if strings.HasPrefix(key, "_") {
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
		} else if r.Method == "POST" {
			var body map[string]interface{}

			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// check if id is already exists

			for _, item := range response.([]interface{}) {
				item := item.(map[string]interface{})
				if item["id"] == body["id"] {
					w.WriteHeader(http.StatusConflict)
					return
				}
			}

			response = append(response.([]interface{}), body)
			s.data[path] = response

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(body)
		} else if r.Method == "DELETE" {
			// example path: /posts/1
			id := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]

			idInt, _ := strconv.Atoi(id)

			for index, item := range response.([]interface{}) {
				item := item.(map[string]interface{})
				if item["id"] == idInt {
					response = append(response.([]interface{})[:index], response.([]interface{})[index+1:]...)
					s.data[path] = response
					break
				}
			}

			w.WriteHeader(http.StatusNoContent)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

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
		},
		Action: func(c *cli.Context) error {
			host := c.String("host")
			port := c.String("port")
			quiet := c.Bool("quiet")
			noCors := c.Bool("no-cors")

			file := c.Args().First()

			if file == "" {
				return cli.Exit("File name is required.", 1)
			}

			data := read(file)

			opt := &ServeOptions{
				host:   host,
				port:   port,
				queit:  quiet,
				noCors: noCors,
			}

			serve(data, opt)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
