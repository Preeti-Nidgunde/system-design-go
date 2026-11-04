package movieticket

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

var (
	ErrShowPricingNotFound = errors.New("Show pricing doesn't exists")
)

type ShowPricing struct {
	ID, ShowID string
	SeatType   SeatType
	Price      int64
}

type ShowPricingRepo struct {
	showPricingRepo map[string]*ShowPricing
	spMu            sync.Mutex
}

func NewShowPricingRepo() *ShowPricingRepo {
	return &ShowPricingRepo{
		showPricingRepo: make(map[string]*ShowPricing),
	}
}

func (c *ShowPricingRepo) CreateShowPrice(showID string, seatType SeatType, price int64) {
	c.spMu.Lock()
	defer c.spMu.Unlock()

	ID := uuid.NewString()
	sp := &ShowPricing{
		ID:       ID,
		ShowID:   showID,
		SeatType: seatType,
		Price:    price,
	}
	c.showPricingRepo[sp.ID] = sp
	fmt.Printf("Length %d", len(c.showPricingRepo))
}

func (c *ShowPricingRepo) GetShowPrice(showID string, seat string) (int64, error) {
	c.spMu.Lock()
	defer c.spMu.Unlock()

	for _, sp := range c.showPricingRepo {
		if sp.ShowID == showID {
			return sp.Price, nil
		}
	}

	return 0, ErrShowPricingNotFound
}
