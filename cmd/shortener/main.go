package main

import (
	"github.com/go-chi/chi/v5"
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

func PostHandler(res http.ResponseWriter, req *http.Request) {
	if mapper.URL == nil {
		mapper = Mapper{}
		mapper.URL = make(map[int]string)
	}
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
}

func GetHandler(res http.ResponseWriter, req *http.Request) {
	if mapper.URL == nil {
		mapper = Mapper{}
		mapper.URL = make(map[int]string)
	}
	i, err := strconv.Atoi(req.URL.String()[1:])
	if err != nil {
		log.Println(err)
	}
	longURL := mapper.URL[i]
	res.Header().Set("Location", longURL)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func BadRequestHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusBadRequest)
	res.Write([]byte("400 StatusBadRequest"))
}

func myRouter() chi.Router {
	r := chi.NewRouter()
	r.Post("/*", PostHandler)
	r.Get("/{id}", GetHandler)
	r.Put("/*", BadRequestHandler)
	r.Delete("/*", BadRequestHandler)
	r.Options("/*", BadRequestHandler)
	r.Head("/*", BadRequestHandler)
	r.Trace("/*", BadRequestHandler)
	r.Connect("/*", BadRequestHandler)
	r.Patch("/*", BadRequestHandler)
	return r
}

func run() error {
	return http.ListenAndServe(`:8080`, myRouter())
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
