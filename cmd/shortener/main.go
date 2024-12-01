package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/ipcross/urlShortener/config"
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
	mapper.Counter++
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
		return
	}
	if len(body) == 0 {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Base URL not correct"))
		return
	}
	mapper.URL[mapper.Counter] = string(body)
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(config.ServerSettings.AddressBase + "/" + strconv.Itoa(mapper.Counter)))
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
	return r
}

func run() error {
	config.InitSettings()
	return http.ListenAndServe(config.ServerSettings.AddressRun, myRouter())
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
