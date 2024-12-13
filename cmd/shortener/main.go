package main

import (
	"fmt"
	"log"

	"github.com/ipcross/urlShortener/internal/config"
	"github.com/ipcross/urlShortener/internal/handlers"
	"github.com/ipcross/urlShortener/internal/repository"
	"github.com/ipcross/urlShortener/internal/service"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg := config.GetConfig()

	store := repository.NewStore()
	mapperService := service.NewMapper(store)

	err := handlers.Serve(cfg, mapperService)
	return fmt.Errorf("run wrap: %w", err)
}
