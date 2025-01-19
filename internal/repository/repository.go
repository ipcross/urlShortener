package repository

import (
	"errors"
	"fmt"
	"log"
	"sync"

	db "github.com/ipcross/urlShortener/internal/adapters/dbstorage"
	"github.com/ipcross/urlShortener/internal/adapters/filestorage"
	"github.com/ipcross/urlShortener/internal/config"
)

type Repository interface {
	GetMapper(req *GetMapperRequest) (*GetMapperResponse, error)
	SetMapper(req *SetMapperRequest)
}

type Store struct {
	mux *sync.Mutex
	s   map[string]string
	cfg config.ServerSettings
}

func NewStore(cfg config.ServerSettings) *Store {
	mapFromStorage := make(map[string]string)

	if len(cfg.DBStorage) == 0 {
		if err := filestorage.NewConsumer(cfg.FileStorage); err != nil {
			log.Printf("NewStore create consumer: %v", err)
		}
		events, err := filestorage.GetConsumer().GetEvents()
		if err != nil {
			log.Printf("Get events: %v", err)
		}
		for _, event := range events {
			mapFromStorage[event.ShortURL] = event.OriginalURL
		}
		if err := filestorage.NewProducer(cfg.FileStorage); err != nil {
			log.Printf("NewStore create producer: %v", err)
		}
	} else {
		data, err := db.LoadDBData(cfg.DBStorage)
		if err != nil {
			log.Printf("Get events: %v", err)
		}
		for _, event := range data {
			mapFromStorage[event.ShortURL] = event.OriginalURL
		}
	}

	return &Store{
		mux: &sync.Mutex{},
		s:   mapFromStorage,
		cfg: cfg,
	}
}

type GetMapperRequest struct {
	Key string
}

type GetMapperResponse struct {
	URL string
}

var (
	ErrGetMapperNotFound = errors.New("url not found")
	ErrSetMapperKeyExist = errors.New("key exist")
)

func newErrGetMapperNotFound(s string) error {
	return fmt.Errorf("%w for KEY = %s", ErrGetMapperNotFound, s)
}

func newErrSetMapperKeyExist(s string) error {
	return fmt.Errorf("%w for KEY = %s", ErrSetMapperKeyExist, s)
}

func (s *Store) GetMapper(req *GetMapperRequest) (*GetMapperResponse, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	res, ok := s.s[req.Key]
	if !ok || res == "" {
		return nil, newErrGetMapperNotFound(req.Key)
	}
	return &GetMapperResponse{
		URL: res,
	}, nil
}

type SetMapperRequest struct {
	Key string
	URL string
}

func (s *Store) SetMapper(req *SetMapperRequest) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	_, ok := s.s[req.Key]
	if ok {
		return newErrSetMapperKeyExist(req.Key)
	}

	s.s[req.Key] = req.URL
	if len(s.cfg.DBStorage) == 0 {
		if err := saveToFile(req.Key, req.URL); err != nil {
			log.Printf("Error saveToFile: %v", err)
		}
	} else {
		if err := saveToDB(s.cfg.DBStorage, req.Key, req.URL); err != nil {
			log.Printf("Error saveToFile: %v", err)
		}
	}
	return nil
}

func saveToDB(dsn string, key string, url string) error {
	event := db.Event{ShortURL: key, OriginalURL: url}
	if err := db.InsertRecord(dsn, &event); err != nil {
		return fmt.Errorf("failed to saveToFile: %w", err)
	}
	return nil
}

func saveToFile(key string, url string) error {
	event := filestorage.Event{ShortURL: key, OriginalURL: url}
	if err := filestorage.GetProducer().WriteEvent(&event); err != nil {
		return fmt.Errorf("failed to saveToFile: %w", err)
	}
	return nil
}
