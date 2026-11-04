package movieticket

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type fixture struct {
	system       *BookMyShow
	show         *Show
	seat1, seat2 *Seat
}

func newFixture(t *testing.T) *fixture {
	t.Helper()
	movieRepo := NewMovieRepo()
	screenRepo := NewScreenRepo()
	showRepo := NewInMemoryShowRepository(movieRepo, screenRepo)
	showPricingRepo := NewShowPricingRepo()

	theaterSystem := NewTheaterSystem()
	seatRepo := NewSeatRepo()
	seatLock := NewSeatLock(5 * time.Minute)
	bookingSystem := NewBookMyShow(movieRepo, showRepo, showPricingRepo, screenRepo, theaterSystem, seatRepo, seatLock)

	movie, err := bookingSystem.movieRepository.CreateMovie("Dhrushyam3", ADULT)
	if err != nil {
		t.Fatalf("create movie: %v", err)
	}
	theater, err := bookingSystem.theaterSystem.CreateTheater("Demo", "Dubai", "AL-Barsha")
	if err != nil {
		t.Fatalf("create theater: %v", err)
	}
	screen, err := bookingSystem.screeRepo.CreateScreen(theater.ID, "A", 50)
	if err != nil {
		t.Fatalf("create screen: %v", err)
	}

	today := time.Now()
	show, err := bookingSystem.showRepository.CreateShow(screen.ID, movie.ID, today, today.Add(3*time.Hour))
	if err != nil {
		t.Fatalf("create show: %v", err)
	}

	bookingSystem.showPricingRepo.CreateShowPrice(show.ID, NORMAL, 50)
	bookingSystem.showPricingRepo.CreateShowPrice(show.ID, SILVER, 100)
	bookingSystem.showPricingRepo.CreateShowPrice(show.ID, GOLD, 200)

	return &fixture{
		system: bookingSystem,
		show:   show,
		seat1:  bookingSystem.SeatSystem.CreateSeat(screen.ID, 5, NORMAL),
		seat2:  bookingSystem.SeatSystem.CreateSeat(screen.ID, 5, SILVER),
	}
}

func TestBooking(t *testing.T) {
	f := newFixture(t)

	if err := f.system.Booking(f.show.ID, "A", []string{f.seat1.ID, f.seat2.ID}); err != nil {
		t.Fatalf("Error occured %s", err)
	}
}

func TestBookingSameSeatTwiceFails(t *testing.T) {
	f := newFixture(t)

	if err := f.system.Booking(f.show.ID, "A", []string{f.seat1.ID}); err != nil {
		t.Fatalf("first booking: %v", err)
	}
	err := f.system.Booking(f.show.ID, "B", []string{f.seat1.ID})
	if !errors.Is(err, ErrSeatUnavailable) {
		t.Fatalf("expected ErrSeatUnavailable, got %v", err)
	}
}

func TestBookingDuplicateSeatInRequestFails(t *testing.T) {
	f := newFixture(t)

	err := f.system.Booking(f.show.ID, "A", []string{f.seat1.ID, f.seat1.ID})
	if !errors.Is(err, ErrDuplicateSeat) {
		t.Fatalf("expected ErrDuplicateSeat, got %v", err)
	}
}

func TestBookingIsAllOrNothing(t *testing.T) {
	f := newFixture(t)

	if err := f.system.Booking(f.show.ID, "A", []string{f.seat1.ID}); err != nil {
		t.Fatalf("first booking: %v", err)
	}
	// seat1 is taken, so seat2 must not be left locked by B's failed attempt
	if err := f.system.Booking(f.show.ID, "B", []string{f.seat2.ID, f.seat1.ID}); !errors.Is(err, ErrSeatUnavailable) {
		t.Fatalf("expected ErrSeatUnavailable, got %v", err)
	}
	if err := f.system.Booking(f.show.ID, "C", []string{f.seat2.ID}); err != nil {
		t.Fatalf("seat2 should still be free: %v", err)
	}
}

func TestConcurrentBookingSameSeat(t *testing.T) {
	f := newFixture(t)

	const users = 100
	var (
		wg        sync.WaitGroup
		successes atomic.Int32
		start     = make(chan struct{})
	)
	for i := 0; i < users; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			if err := f.system.Booking(f.show.ID, fmt.Sprintf("user-%d", i), []string{f.seat1.ID}); err == nil {
				successes.Add(1)
			} else if !errors.Is(err, ErrSeatUnavailable) {
				t.Errorf("unexpected error: %v", err)
			}
		}(i)
	}
	close(start)
	wg.Wait()

	if got := successes.Load(); got != 1 {
		t.Fatalf("expected exactly 1 successful booking, got %d", got)
	}
}

func TestSeatLockExpires(t *testing.T) {
	now := time.Now()
	p := NewSeatLock(time.Minute)
	p.now = func() time.Time { return now }

	if err := p.Lock("A", []string{"s1"}, "show"); err != nil {
		t.Fatalf("lock A: %v", err)
	}
	if err := p.Lock("B", []string{"s1"}, "show"); !errors.Is(err, ErrSeatUnavailable) {
		t.Fatalf("expected ErrSeatUnavailable while A holds lock, got %v", err)
	}

	now = now.Add(2 * time.Minute)

	if err := p.Lock("B", []string{"s1"}, "show"); err != nil {
		t.Fatalf("B should get expired seat: %v", err)
	}
	if err := p.Confirm("A", []string{"s1"}, "show"); !errors.Is(err, ErrSeatUnavailable) {
		t.Fatalf("A's confirm after expiry should fail, got %v", err)
	}
	if err := p.Confirm("B", []string{"s1"}, "show"); err != nil {
		t.Fatalf("B confirm: %v", err)
	}

	now = now.Add(time.Hour)
	if err := p.Lock("C", []string{"s1"}, "show"); !errors.Is(err, ErrSeatUnavailable) {
		t.Fatalf("booked seat must never expire, got %v", err)
	}
}
