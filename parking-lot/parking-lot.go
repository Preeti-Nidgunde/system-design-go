package parkinglot

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrAlreadyParked = errors.New("Already used this ticket to unpark")
)

type ParkingLot struct {
	Name            string
	parkingManager  *ParkingManager    // stateful entity
	activateTickets map[string]*Ticket // map to allow mutation hold reference ot Ticket
	feeCalcultor    FeeCalcultor       // interface never use pointer as underlying structure has its pointer
	mu              sync.Mutex         // go map, slice, concurrency primitive dont need *
}

func NewParkingKLot(feeCalc FeeCalcultor) *ParkingLot {
	return &ParkingLot{
		Name:            "Parking lot system",
		parkingManager:  NewParkingManager(),
		activateTickets: map[string]*Ticket{},
		feeCalcultor:    feeCalc,
	}
}

func (pl *ParkingLot) AddParkingSpot(parkingSpot ParkingSpot) {
	pl.parkingManager.AddSpot(parkingSpot)
}

func (pl *ParkingLot) ParkVehicle(vehicle *Vehicle) (*Ticket, error) {
	spot, err := pl.parkingManager.ParkVehicle(vehicle)
	if err != nil {
		return nil, err
	}

	ticket := NewTicket(*vehicle, spot)
	pl.activateTickets[ticket.ID] = ticket
	return ticket, nil
}

func (pl *ParkingLot) LeaveVehicle(ticket *Ticket) error {
	ticket, exist := pl.activateTickets[ticket.ID]
	if !exist {
		return fmt.Errorf("vehicle not parked legally")
	}

	if !ticket.ExitTime.IsZero() {
		return ErrAlreadyParked
	}

	ticket.ExitTime = time.Now()
	pl.parkingManager.UnParkVehicle(ticket.Spot)

	price, err := pl.feeCalcultor.CalculateFee(ticket)
	if err != nil {
		return err
	}

	fmt.Printf("pay the price %d", price)
	delete(pl.activateTickets, ticket.ID) // Delete the ticekt ID

	return nil
}

// Context								Use Pointer (*T)	Use Value (T)			Reason
// Struct with sync.Mutex			✅ Always			❌ Never					Copying locks breaks concurrency safety.
// Methods that modify fields		✅ Always			❌ Never					Value receivers modify temporary copies.
// Interfaces (FeeCalcultor)		❌ Never				✅ Always				Interfaces already hold underlying pointers/values.
// Maps holding mutable structs		✅ map[K]*V			❌ map[K]V				Go disallows map[key].Field = val for values.
// Small read-only data (time.Time)	❌ Rare				✅ Always				Cheap to copy, immutable behavior.
