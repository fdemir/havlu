package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func HandleBase(w http.ResponseWriter, r *http.Request, s *Source, opt *ServeOptions) {
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

	switch r.Method {
	case http.MethodGet:
		data := GetAll(
			r.URL.Query(),
			response,
		)

		json.NewEncoder(w).Encode(data)
	case http.MethodPost:
		Create(
			r.Body,
			response,
		)

		w.WriteHeader(http.StatusCreated)
	case http.MethodDelete:
		id := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]

		Delete(
			id,
			response,
		)

		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}
