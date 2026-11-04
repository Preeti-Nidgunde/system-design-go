package movieticket

import (
	"sync"

	"github.com/google/uuid"
)

type BookMyShow struct {
	movieRepository *MovieRepo
	showRepository  ShowRepository
	showPricingRepo *ShowPricingRepo
	screeRepo       *ScreenRepo
	theaterSystem   *TheaterSystem
	SeatSystem      *SeatRepo
	seatLock        SeatLockProvider
	booking         map[string]*Booking
	bMu             sync.Mutex
	payments        map[string]*Payment
	pMu             sync.RWMutex
}

func NewBookMyShow(movieRepo *MovieRepo, showRepo ShowRepository, showPrice *ShowPricingRepo, screenRepo *ScreenRepo, theaterSystem *TheaterSystem, seatRepo *SeatRepo, seatLock SeatLockProvider) *BookMyShow {
	return &BookMyShow{
		movieRepository: movieRepo,
		showRepository:  showRepo,
		showPricingRepo: showPrice,
		screeRepo:       screenRepo,
		theaterSystem:   theaterSystem,
		SeatSystem:      seatRepo,
		seatLock:        seatLock,
		booking:         make(map[string]*Booking),
		payments:        make(map[string]*Payment),
	}
}

func (b *BookMyShow) Booking(showID, userID string, seats []string) error {
	var total int64
	_, err := b.showRepository.GetShowByID(showID)
	if err != nil {
		return err
	}

	// hold the seats before taking money; fails if any seat is locked/booked
	if err := b.seatLock.Lock(userID, seats, showID); err != nil {
		return err
	}

	for _, st := range seats {
		price, err := b.showPricingRepo.GetShowPrice(showID, st)
		if err != nil {
			b.seatLock.UnLock(userID, seats, showID)
			return err
		}
		total += price
	}

	pay := &Payment{
		ID:     uuid.NewString(),
		Amount: total,
		Status: "DONE",
	}

	b.pMu.Lock()
	b.payments[pay.ID] = pay
	b.pMu.Unlock()

	// lock may have expired while paying; in a real system this triggers a refund
	if err := b.seatLock.Confirm(userID, seats, showID); err != nil {
		b.seatLock.UnLock(userID, seats, showID)
		return err
	}

	newBooking := NewBooking(showID, userID, pay.ID, CONFIRMED)
	b.bMu.Lock()
	b.booking[newBooking.ID] = newBooking
	b.bMu.Unlock()
	return nil
}
