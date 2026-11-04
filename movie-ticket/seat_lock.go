package movieticket

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type SeatStatus int

var (
	ErrDuplicateSeat   = errors.New("duplicate seat provided")
	ErrSeatUnavailable = errors.New("seat is not free or taken")
)

const (
	FREE SeatStatus = iota
	OCCUPIED
	LOCKED
)

type SeatLockProvider interface {
	Lock(userID string, seats []string, showID string) error
	Confirm(userID string, seats []string, showID string) error
	UnLock(userID string, seats []string, showID string) error
}

type SeatLockValue struct {
	userID    string
	status    SeatStatus
	expiresAt time.Time
}

type SeatLock struct {
	showSeats map[string]map[string]*SeatLockValue
	ttl       time.Duration
	now       func() time.Time
	mu        sync.Mutex
}

var _ (SeatLockProvider) = (*SeatLock)(nil)

func NewSeatLock(ttl time.Duration) *SeatLock {
	return &SeatLock{
		showSeats: make(map[string]map[string]*SeatLockValue),
		ttl:       ttl,
		now:       time.Now,
	}
}

func isTaken(v *SeatLockValue, now time.Time) bool {
	return v != nil && (v.status == OCCUPIED || (v.status == LOCKED && now.Before(v.expiresAt)))
}

func (sl *SeatLock) Lock(userID string, seats []string, showID string) error {
	visited := make(map[string]struct{}, len(seats))

	for _, seat := range seats {
		if _, ok := visited[seat]; ok {
			return ErrDuplicateSeat
		}
		visited[seat] = struct{}{}
	}

	sl.mu.Lock()
	defer sl.mu.Unlock()

	userSeats, ok := sl.showSeats[showID]
	if !ok {
		userSeats = make(map[string]*SeatLockValue, len(seats))
		sl.showSeats[showID] = userSeats
	}

	now := sl.now() // better to take here
	for _, seat := range seats {
		if isTaken(userSeats[seat], now) {
			return ErrSeatUnavailable
		}
	}

	for _, seat := range seats {
		userSeats[seat] = &SeatLockValue{
			userID:    userID,
			status:    LOCKED,
			expiresAt: sl.now().Add(sl.ttl),
		}
	}
	return nil
}

func isLockHeld(v *SeatLockValue, userID string, t time.Time) bool {
	return v != nil && v.userID == userID && v.status == LOCKED && t.Before(v.expiresAt)
}

func (sl *SeatLock) Confirm(userID string, seats []string, showID string) error {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	userSeats, ok := sl.showSeats[showID]
	if !ok {
		return fmt.Errorf("show seats doesnt exist")
	}

	for _, seat := range seats {
		if !isLockHeld(userSeats[seat], userID, sl.now()) {
			return ErrSeatUnavailable
		}
	}

	for _, seat := range seats {
		userSeats[seat].status = OCCUPIED
	}

	return nil
}

func (sl *SeatLock) UnLock(userID string, seats []string, showID string) error {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	userSeat, ok := sl.showSeats[showID]
	if !ok {
		return ErrShowNotFound
	}

	for _, seat := range seats {
		if userSeat[seat] != nil && userSeat[seat].status == LOCKED && userSeat[seat].userID == userID {
			delete(userSeat, seat)
		}
	}
	return nil
}

// package movieticket

// import (
// 	"errors"
// 	"sync"
// 	"time"
// )

// type SEATLOCK int

// const (
// 	FREE SEATLOCK = iota
// 	LOCK
// 	OCCUPIED
// )

// var (
// 	ErrDuplicateSeat   = errors.New("seat already booked")
// 	ErrSeatUnavailable = errors.New("seat unavailable")
// 	ErrSeatLockNotHeld = errors.New("seat lock not hold")
// )

// type SeatLockProvider interface {
// 	Lock(userID, showID string, seats []string) error
// 	Confirm(userID, showID string, seats []string) error
// 	UnLock(userID, showID string, seats []string) error
// }

// type SeatLockValue struct {
// 	userID    string
// 	status    SEATLOCK
// 	expiresAt time.Time
// }

// type InMemorySeatLock struct {
// 	showSeatMap map[string]map[string]*SeatLockValue // show->seat->SeatValue
// 	ttl         time.Duration
// 	now         func() time.Time
// 	mu          sync.Mutex
// }

// func NewInMemorySeatLock(ttl time.Duration) *InMemorySeatLock {
// 	return &InMemorySeatLock{
// 		showSeatMap: make(map[string]map[string]*SeatLockValue),
// 		ttl:         ttl,
// 		now:         time.Now,
// 	}
// }

// var _ SeatLockProvider = (*InMemorySeatLock)(nil)

// func isFree(seatLockval *SeatLockValue, now time.Time) bool {
// 	return seatLockval == nil || (seatLockval.status == LOCK && !now.Before(seatLockval.expiresAt))
// }

// func (lp *InMemorySeatLock) Lock(userID, showID string, seats []string) error {
// 	visited := make(map[string]struct{}, len(seats))

// 	// To sanitize the input if they are sending same seat twice
// 	for _, seat := range seats {
// 		if _, ok := visited[seat]; ok {
// 			return ErrDuplicateSeat
// 		}
// 		visited[seat] = struct{}{}
// 	}

// 	lp.mu.Lock()
// 	defer lp.mu.Unlock()
// 	now := lp.now() // moved here as goroutine may wait if lock unavialble

// 	showSeats, ok := lp.showSeatMap[showID]
// 	if !ok {
// 		showSeats = make(map[string]*SeatLockValue)
// 		lp.showSeatMap[showID] = showSeats
// 	}

// 	// First check to avoid concurrency issue of holding mutex or half lock
// 	for _, seat := range seats {
// 		if !isFree(showSeats[seat], now) { // seat free is taken by ==nil in isFree method only
// 			return ErrSeatUnavailable
// 		}
// 	}

// 	// Now finally update
// 	for _, seat := range seats {
// 		showSeats[seat] = &SeatLockValue{
// 			userID:    userID,
// 			status:    LOCK,
// 			expiresAt: now.Add(lp.ttl),
// 		}
// 	}

// 	return nil
// }

// func isLockHeld(seatLockval *SeatLockValue, userID, seatID string, now time.Time) bool {
// 	return seatLockval != nil && seatLockval.userID == userID && seatLockval.status == LOCK && now.Before(seatLockval.expiresAt)
// }

// func (lp *InMemorySeatLock) Confirm(userID, showID string, seats []string) error {
// 	lp.mu.Lock()
// 	defer lp.mu.Unlock()

// 	showSeat, ok := lp.showSeatMap[showID]
// 	if !ok {
// 		return ErrSeatUnavailable
// 	}

// 	for _, seat := range seats {
// 		if !isLockHeld(showSeat[seat], userID, seat, lp.now()) {
// 			return ErrSeatLockNotHeld
// 		}
// 	}

// 	for _, seat := range seats {
// 		showSeat[seat].status = OCCUPIED
// 	}

// 	return nil
// }

// func (lp *InMemorySeatLock) UnLock(userID, showID string, seats []string) error {
// 	showSeat := lp.showSeatMap[showID]
// 	for _, seats := range seats {
// 		if val := showSeat[seats]; val != nil && val.userID == userID && val.status == LOCK {
// 			delete(showSeat, seats)
// 		}
// 	}
// 	return nil
// }
