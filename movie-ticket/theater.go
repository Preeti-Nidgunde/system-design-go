package movieticket

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

var (
	ErrTheaterAlreadyExists = errors.New("theater already exists")
	ErrTheaterNotFound      = errors.New("theater not found")
)

type Theater struct {
	ID       string
	Name     string
	Location string
	Address  string
}

type TheaterSystem struct {
	Theaters map[string]*Theater
	mu       sync.Mutex
}

func NewTheaterSystem() *TheaterSystem {
	return &TheaterSystem{
		Theaters: make(map[string]*Theater),
	}
}

func (s *TheaterSystem) GetTheater(id string) (*Theater, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	theater, ok := s.Theaters[id]
	if !ok {
		return nil, ErrTheaterNotFound
	}

	return theater, nil
}

func (s *TheaterSystem) CreateTheater(name, location, address string) (*Theater, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ID := uuid.NewString()

	s.Theaters[ID] = &Theater{
		ID:       ID,
		Name:     name,
		Location: location,
		Address:  address,
	}

	return s.Theaters[ID], nil
}
