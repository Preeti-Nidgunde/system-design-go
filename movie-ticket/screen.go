package movieticket

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

var (
	ErrScreenUnavailable = errors.New("screen not found")
)

type Screen struct {
	ID        string
	TheaterID string
	Name      string
	Capacity  int64
}

type ScreenRepo struct {
	screens map[string]*Screen
	scMu    sync.Mutex
}

func NewScreenRepo() *ScreenRepo {
	return &ScreenRepo{
		screens: make(map[string]*Screen),
	}
}
func (c *ScreenRepo) CreateScreen(theaterID, name string, cap int64) (*Screen, error) {

	c.scMu.Lock()
	defer c.scMu.Unlock()

	screen := &Screen{
		ID:        uuid.NewString(),
		Name:      name,
		TheaterID: theaterID,
		Capacity:  cap,
	}

	c.screens[screen.ID] = screen
	return screen, nil
}

func (c *ScreenRepo) GetScreen(screenID string) (*Screen, error) {
	c.scMu.Lock()
	defer c.scMu.Unlock()

	for _, screen := range c.screens {
		if screen.ID == screenID {
			return screen, nil
		}
	}

	return nil, ErrScreenUnavailable

}
