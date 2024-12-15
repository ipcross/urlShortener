package service

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/ipcross/urlShortener/internal/config"
	"github.com/ipcross/urlShortener/internal/repository"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	hashSize    = 6
)

type Repository interface {
	GetMapper(req *repository.GetMapperRequest) (*repository.GetMapperResponse, error)
	SetMapper(req *repository.SetMapperRequest) error
}

type Mapper struct {
	store Repository
	cfg   config.ServerSettings
}

func NewMapper(cfg config.ServerSettings, store Repository) *Mapper {
	return &Mapper{
		store: store,
		cfg:   cfg,
	}
}

type GetMapperRequest struct {
	Key string
}

type GetMapperResponse struct {
	URL string
}

func (f *Mapper) GetMapper(req *GetMapperRequest) (*GetMapperResponse, error) {
	repositoryResp, err := f.store.GetMapper(&repository.GetMapperRequest{
		Key: req.Key,
	})
	if err != nil {
		if errors.Is(err, repository.ErrGetMapperNotFound) {
			return nil, fmt.Errorf("failed to fetch data from the store: %w", err)
		}
	}

	return &GetMapperResponse{
		URL: repositoryResp.URL,
	}, nil
}

type SetMapperRequest struct {
	URL string `json:"url"`
}

type SetMapperResponse struct {
	Link string `json:"result"`
}

func (f *Mapper) SetMapper(req *SetMapperRequest) (*SetMapperResponse, error) {
	newHash := f.generateKey()
	err := f.store.SetMapper(&repository.SetMapperRequest{
		Key: newHash,
		URL: req.URL,
	})
	if err != nil {
		if errors.Is(err, repository.ErrSetMapperKeyExist) {
			return nil, fmt.Errorf("failed to set data in the store: %w", err)
		}
	}

	return &SetMapperResponse{
		Link: f.cfg.AddressBase + "/" + newHash,
	}, nil
}

func (f *Mapper) generateKey() string {
	for {
		hash := randStringBytesRmndr(hashSize)
		_, err := f.store.GetMapper(&repository.GetMapperRequest{
			Key: hash,
		})
		if err != nil {
			if errors.Is(err, repository.ErrGetMapperNotFound) {
				return hash
			}
		}
	}
}

func randStringBytesRmndr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
