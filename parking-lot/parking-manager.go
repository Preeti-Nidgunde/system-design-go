package parkinglot

import (
	"fmt"
)

type ParkingManager struct {
	freeSpots map[Vehiclesize][]ParkingSpot
}

func NewParkingManager() *ParkingManager {
	return &ParkingManager{
		freeSpots: make(map[Vehiclesize][]ParkingSpot),
	}
}

// This is construction phase like adding new spots on any floor.
func (pm *ParkingManager) AddSpot(parkingSpot ParkingSpot) {
	pm.freeSpots[parkingSpot.GetSize()] = append(pm.freeSpots[parkingSpot.GetSize()], parkingSpot)
}

// Add vehicle to existing packing spot
func (pm *ParkingManager) ParkVehicle(vehicle *Vehicle) (ParkingSpot, error) {
	spotLists, exists := pm.freeSpots[vehicle.size]
	if !exists || len(spotLists) == 0 {
		return nil, fmt.Errorf("parking unavailable")
	}

	targetSpot := spotLists[len(spotLists)-1]
	pm.freeSpots[vehicle.size] = spotLists[:len(spotLists)-1]

	if err := targetSpot.Occupy(vehicle); err != nil {
		return nil, err
	}

	return targetSpot, nil
}

func (pm *ParkingManager) UnParkVehicle(spot ParkingSpot) (*Vehicle, error) {
	if spot == nil {
		return nil, fmt.Errorf("invalid spot")
	}

	vehicle, err := spot.Vacate()
	if err != nil {
		return nil, err
	}

	pm.freeSpots[spot.GetSize()] = append(pm.freeSpots[spot.GetSize()], spot)
	return vehicle, nil
}
