package notifier

import (
	"context"
	"errors"
	"sync"
	"time"
)

type EventData struct {
	data   string
	sentAt time.Time
}

type Producer interface {
	OnEvent(data EventData) error
}

type EventGenerator struct {
	obs []Observers
	mu  sync.Mutex
}

func NewEventGenerator() *EventGenerator {
	return &EventGenerator{
		obs: make([]Observers, 0, 3),
	}
}

func (eg *EventGenerator) Subscribe(o Observers) {
	eg.mu.Lock()
	defer eg.mu.Unlock()

	eg.obs = append(eg.obs, o)
}

func (eg *EventGenerator) Unsubscribe(o Observers) error {
	eg.mu.Lock()
	defer eg.mu.Unlock()

	for i, ob := range eg.obs {
		if ob == o {
			eg.obs = append(eg.obs[:i], eg.obs[i+1:]...)
			return nil
		}
	}
	return errors.New("observer systm not found")
}

func (eg *EventGenerator) OnEvent(data EventData) error {
	eg.mu.Lock()
	defer eg.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	for _, ob := range eg.obs {
		ob.Update(ctx, data)
	}

	return nil
}
