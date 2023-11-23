package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

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

	result := bytes.NewBuffer(nil)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		result.Write(scanner.Bytes())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	jsonMap := make(map[string]interface{})

	json.Unmarshal(result.Bytes(), &jsonMap)

	return &Source{data: jsonMap}
}

func serve(s Source, host string, port string) {
	// todo: syncronize data access with mutex or channels
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		path := r.URL.Path[1:]
		response := s.data[path]

		if response == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(response)
	})

	fmt.Printf("Listening on %s:%s\n", host, port)
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
		},
		Action: func(c *cli.Context) error {
			host := c.String("host")
			port := c.String("port")
			file := c.String("file")

			data := read(file)
			serve(*data, host, port)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
