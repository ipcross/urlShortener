package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ipcross/urlShortener/internal/config"
	l "github.com/ipcross/urlShortener/internal/logger"
	"github.com/ipcross/urlShortener/internal/service"
	"go.uber.org/zap"
)

func Serve(cfg config.ServerSettings, mapper Mapper) error {
	logger, err := l.Initialize(cfg.LogLevel)
	if err != nil {
		return fmt.Errorf("logger.Initialize: %w", err)
	}
	defer l.Sync(logger)

	logger.Info("Running server", zap.String("address", cfg.AddressRun))

	h := NewHandlers(mapper, cfg)
	router := myRouter(h, logger)

	srv := &http.Server{
		Addr:    cfg.AddressRun,
		Handler: router,
	}

	err = srv.ListenAndServe()
	return fmt.Errorf("handlers.Serve wrap: %w", err)
}

func myRouter(h *handlers, logger *zap.Logger) chi.Router {
	r := chi.NewRouter()
	r.Use(l.RequestLogger(logger))
	r.Post("/*", h.PostHandler)
	r.Get("/{key}", h.GetHandler)
	r.Put("/*", h.BadRequestHandler)
	return r
}

type Mapper interface {
	GetMapper(req *service.GetMapperRequest) (*service.GetMapperResponse, error)
	SetMapper(req *service.SetMapperRequest) (*service.SetMapperResponse, error)
}

type handlers struct {
	mapper Mapper
	config config.ServerSettings
}

func NewHandlers(mapper Mapper, cfg config.ServerSettings) *handlers {
	return &handlers{
		mapper: mapper,
		config: cfg,
	}
}

func (h *handlers) BadRequestHandler(res http.ResponseWriter, _ *http.Request) {
	res.WriteHeader(http.StatusBadRequest)
	_, err := res.Write([]byte("400 StatusBadRequest"))
	if err != nil {
		log.Printf("Error writing to response: %v", err)
		return
	}
}

func (h *handlers) GetHandler(res http.ResponseWriter, req *http.Request) {
	keyStr := req.PathValue("key")
	if len(keyStr) == 0 {
		res.WriteHeader(http.StatusBadRequest)
		_, err := res.Write([]byte("URL not correct"))
		if err != nil {
			log.Printf("Error writing to response: %v", err)
			return
		}
		return
	}

	resp, err := h.mapper.GetMapper(&service.GetMapperRequest{
		Key: keyStr,
	})
	if err != nil {
		log.Printf("failed to get URL: %v", err)
		http.Error(res, "Not found", http.StatusBadRequest)
		return
	}

	res.Header().Set("Location", resp.URL)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *handlers) PostHandler(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil || len(body) == 0 {
		log.Printf("PostHandler err: %v", err)
		http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	resp, err := h.mapper.SetMapper(&service.SetMapperRequest{
		URL: string(body),
	})
	if err != nil {
		log.Printf("failed to save URL: %v", err)
		http.Error(res, "Failed to save URL", http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	_, err = res.Write([]byte(h.config.AddressBase + "/" + resp.Key))
	if err != nil {
		log.Printf("Error writing to response: %v", err)
		return
	}
}
