package main

import (
	"io"
	"log"
	"net/http"
	"strconv"
)

type Mapper struct {
	Counter int
	URL     map[int]string
}

var mapper Mapper

func myHandler(res http.ResponseWriter, req *http.Request) {
	if mapper.URL == nil {
		mapper = Mapper{}
		mapper.URL = make(map[int]string)
	}

	switch req.Method {
	case http.MethodPost:
		body, err := io.ReadAll(req.Body)
		if err != nil {
			log.Println(err)
			return
		}
		mapper.URL[mapper.Counter] = string(body)
		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte("http://localhost:8080/" + strconv.Itoa(mapper.Counter)))
		mapper.Counter++
		return
	case http.MethodGet:
		i, err := strconv.Atoi(req.URL.String()[1:])
		if err != nil {
			log.Println(err)
		}
		longURL := mapper.URL[i]
		res.Header().Set("Location", longURL)
		res.WriteHeader(http.StatusTemporaryRedirect)
		return
	default:
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("400 StatusBadRequest"))
	}
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	mapper = Mapper{}
	mapper.URL = make(map[int]string)

	mux := http.NewServeMux()
	mux.HandleFunc(`/`, myHandler)

	return http.ListenAndServe(`:8080`, mux)
}
