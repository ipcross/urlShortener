package repository

import (
	"errors"
	"fmt"
	"sync"
)

type Repository interface {
	GetMapper(req *GetMapperRequest) (*GetMapperResponse, error)
	SetMapper(req *SetMapperRequest)
}

type Store struct {
	mux *sync.Mutex
	s   map[string]string
}

func NewStore() *Store {
	return &Store{
		mux: &sync.Mutex{},
		s:   make(map[string]string),
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
	return nil
}
