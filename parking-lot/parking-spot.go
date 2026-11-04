package parkinglot

import (
	"fmt"
	"sync"
)

type ParkingSpot interface {
	GetSpotNumber() string
	IsAvailable() bool
	CanFit(v *Vehicle) bool
	Occupy(vehicle *Vehicle) error
	Vacate() (*Vehicle, error)
	GetSize() Vehiclesize
}

type baseSpot struct {
	spotNo   string
	vehicle  *Vehicle
	spotSize Vehiclesize
	sync.Mutex
}

// func NewBaseSpot(spotSize string) *baseSpot {
// 	return &baseSpot{
// 		vehicle: nil,
// 		spotNo:  spotSize,
// 	}
// }

func (bs *baseSpot) GetSpotNumber() string {
	return bs.spotNo
}

func (bs *baseSpot) IsAvailable() bool {
	bs.Lock()
	defer bs.Unlock()
	return bs.vehicle == nil
}

func (bs *baseSpot) CanFit(v *Vehicle) bool {
	return bs.spotSize >= v.GetSize()
}

func (bs *baseSpot) Occupy(vehicle *Vehicle) error {
	bs.Lock()
	defer bs.Unlock()

	if bs.vehicle != nil {
		return fmt.Errorf("spot occupied")
	}

	bs.vehicle = vehicle
	return nil
}

func (bs *baseSpot) Vacate() (*Vehicle, error) {
	bs.Lock()
	defer bs.Unlock()

	if bs.vehicle == nil {
		return nil, fmt.Errorf("already freed")
	}

	v := bs.vehicle
	bs.vehicle = nil
	return v, nil
}

func (bs *baseSpot) GetSize() Vehiclesize {
	return bs.spotSize
}

type SmallSpot struct {
	baseSpot
}

type MediumSpot struct {
	baseSpot
}

type LargeSpot struct {
	baseSpot
}

type Handicaped struct {
	baseSpot
}

func NewSmallSpot(spotNo string) *SmallSpot {
	return &SmallSpot{
		baseSpot{
			spotNo:   spotNo,
			spotSize: SMALL,
		},
	}
}

func NewMediumSpot(spotNo string) *MediumSpot {
	return &MediumSpot{
		baseSpot{
			spotNo:   spotNo,
			spotSize: MEDIUM,
		},
	}
}

func NewLargeSpot(spotNo string) *LargeSpot {
	return &LargeSpot{
		baseSpot{
			spotNo:   spotNo,
			spotSize: LARGE,
		},
	}
}
