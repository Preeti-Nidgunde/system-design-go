package parkinglot

import (
	"fmt"
	"math"
)

type FeeCalcultor interface {
	CalculateFee(t *Ticket) (int, error)
}

type BaseCalcultor struct {
	small  int
	medium int
	large  int
}

type HourlyRateCalcultor struct {
	multiplier int
}

func (b *BaseCalcultor) CalculateFee(t *Ticket) (int, error) {
	if t.ExitTime.Before(t.EntryTime) {
		return 0, fmt.Errorf("exit time cannot be before entry time")
	}

	duration := t.ExitTime.Sub(t.EntryTime)
	hours := int(math.Ceil(duration.Hours()))

	if hours == 0 {
		hours = 1
	}

	var baseRate int
	switch t.Vehicle.size {
	case SMALL:
		baseRate = b.small
	case MEDIUM:
		baseRate = b.medium
	case LARGE:
		baseRate = b.large
	default:
		return 0, fmt.Errorf("unknown vehicle size %v", t.Vehicle.size)
	}

	return baseRate * hours, nil
}

func (h *HourlyRateCalcultor) CalculateFee(t *Ticket) (int, error) {
	if t.ExitTime.Before(t.EntryTime) {
		return 0, fmt.Errorf("exit time cannot be before entry, vehicle not started waiting %v", t.Vehicle.number)
	}

	duration := int(t.ExitTime.Sub(t.EntryTime).Hours())
	if duration == 0 {
		duration = 1
	}

	return h.multiplier * duration * (int(t.Vehicle.size) + 1), nil
}
