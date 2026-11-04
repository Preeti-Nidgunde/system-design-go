package parkinglot

import (
	"errors"
	"testing"
)

func TestParkVehicle(t *testing.T) {
	baseCal := BaseCalcultor{small: 10, medium: 15, large: 30}
	parkingLot := NewParkingKLot(&baseCal)

	smallSPot := NewSmallSpot("101")
	parkingLot.AddParkingSpot(smallSPot)
	veh := &Vehicle{number: "MH096417", size: SMALL}
	//largeVeh := &Vehicle{number: "msjjd", size: LARGE}

	ticket, err := parkingLot.ParkVehicle(veh)
	if err != nil {
		t.Errorf("got error %v", err)
	}

	t.Logf("ticket details %s %v", ticket.ID, ticket.EntryTime)

	// trying to fit large vehicle
	// ticket, err = parkingLot.ParkVehicle(largeVeh)
	// if err != nil && err.Error() == "parking unavailable" {
	// 	t.Skipf("expected error %v", err)
	// }

	err = parkingLot.LeaveVehicle(ticket)
	if err != nil {
		t.Errorf("\n got error %v", err)
	}

	err = parkingLot.LeaveVehicle(ticket)
	if errors.Is(err, ErrAlreadyParked) {
		t.Logf("Expected error %v", err)
	}
}

// Method	Fails the test?	Stops execution?
// t.Log / t.Logf	❌ No	❌ No
// t.Error / t.Errorf	✅ Yes	❌ No (keeps running)
// t.Fatal / t.Fatalf	✅ Yes	✅ Yes (calls runtime.Goexit)
// t.Fail	✅ Yes	❌ No
// t.FailNow	✅ Yes	✅ Yes
// t.Skip / t.Skipf / t.SkipNow	❌ No (skipped)	✅ Yes
// t.Helper	❌ No	❌ No
// t.Parallel	❌ No	(pauses, resumes later)
// t.Cleanup	❌ No	❌ No
