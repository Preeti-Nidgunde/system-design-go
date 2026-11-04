package movieticket

import (
	"sync"

	"github.com/google/uuid"
)

type Seat struct {
	ID, ScreenID string
	Row          int64
	Type         SeatType
}

type SeatRepo struct {
	seats map[string]*Seat
	seMu  sync.Mutex
}

func NewSeatRepo() *SeatRepo {
	return &SeatRepo{
		seats: make(map[string]*Seat),
	}
}
func (s *SeatRepo) CreateSeat(screenID string, row int64, seatType SeatType) *Seat {
	s.seMu.Lock()
	defer s.seMu.Unlock()

	newSeat := &Seat{
		ID:       uuid.NewString(),
		ScreenID: screenID,
		Row:      row,
		Type:     seatType,
	}

	s.seats[newSeat.ID] = newSeat
	return newSeat
}
