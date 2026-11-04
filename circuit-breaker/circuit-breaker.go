package circuitbreaker

import (
	"fmt"
	"time"
)

type State int

// StateClose is the zero value so a new breaker lets calls through.
const (
	StateClose State = iota
	StateOpen
	StateHalfOpen
)

type CallOutCome struct {
	Timestamp time.Time
	Status    bool
}

type CircuitBreaker struct {
	State State

	Outcome    []CallOutCome
	WindowSize time.Duration
	Threshold  float64 // failure rate (0..1) that trips the breaker

	CooldownPeriod time.Duration
	OpenedAt       time.Time

	HalfOpenTrials    int
	HalfOpenMaxTrials int
	HalfOpenSuccess   int
}

func (cb *CircuitBreaker) Allow() bool {
	switch cb.State {

	case StateClose:
		return true

	case StateOpen:
		if time.Since(cb.OpenedAt) >= cb.CooldownPeriod {
			cb.transitionTo(StateHalfOpen)
			return true
		}
		return false

	case StateHalfOpen:
		if cb.HalfOpenTrials < cb.HalfOpenMaxTrials {
			cb.HalfOpenTrials++
			return true
		}
		return false
	}
	return false
}

func (cb *CircuitBreaker) recordOutcome(val bool) {
	now := time.Now()

	cutOff := now.Add(-cb.WindowSize)

	i := 0
	for i < len(cb.Outcome) && cb.Outcome[i].Timestamp.Before(cutOff) {
		i++
	}

	cb.Outcome = cb.Outcome[i:]

	cb.Outcome = append(cb.Outcome, CallOutCome{Timestamp: time.Now(), Status: val})
}

func (cb *CircuitBreaker) RecordSuccess() {
	switch cb.State {
	case StateClose:
		cb.recordOutcome(true)
	case StateHalfOpen:
		cb.HalfOpenSuccess++
		if cb.HalfOpenSuccess >= cb.HalfOpenMaxTrials {
			cb.transitionTo(StateClose)
		}

	}
}

func (cb *CircuitBreaker) transitionTo(state State) {
	cb.State = state
	switch state {
	case StateClose:
		cb.Outcome = nil
	case StateOpen:
		cb.OpenedAt = time.Now()
	case StateHalfOpen:
		cb.HalfOpenSuccess = 0
		cb.HalfOpenTrials = 0
	}
}

func (cb *CircuitBreaker) ShouldTrip() bool {
	total := len(cb.Outcome)
	if total == 0 {
		return false
	}

	failures := 0

	for _, val := range cb.Outcome {
		if !val.Status {
			failures++
		}
	}

	failureRate := float64(failures) / float64(total)
	return failureRate >= cb.Threshold
}

func (cb *CircuitBreaker) RecordFailure() {
	switch cb.State {
	case StateClose:
		cb.recordOutcome(false)
		if cb.ShouldTrip() {
			cb.transitionTo(StateOpen)
		}
	case StateHalfOpen:
		cb.transitionTo(StateOpen)

	}
}

func (cb *CircuitBreaker) Execute(fn func() error) error {
	if !cb.Allow() {
		return fmt.Errorf("circuit breaker open")
	}

	err := fn()
	if err != nil {
		cb.RecordFailure()
		return err
	}

	cb.RecordSuccess()
	return nil
}
