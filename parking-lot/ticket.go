package parkinglot

import (
	"time"

	"github.com/google/uuid"
)

type Ticket struct {
	ID        string
	Vehicle   Vehicle
	Spot      ParkingSpot
	EntryTime time.Time
	ExitTime  time.Time
}

func NewTicket(vehicle Vehicle, spot ParkingSpot) *Ticket {
	return &Ticket{
		ID:        uuid.NewString(),
		Vehicle:   vehicle,
		EntryTime: time.Now(),
		Spot:      spot,
	}
}
