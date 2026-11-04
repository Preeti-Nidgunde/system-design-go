package circuitbreaker

import (
	"errors"
	"testing"
	"time"
)

var errDownstream = errors.New("downstream unavailable")

func newTestBreaker() *CircuitBreaker {
	return &CircuitBreaker{
		WindowSize:        time.Minute,
		Threshold:         0.5,
		CooldownPeriod:    20 * time.Millisecond,
		HalfOpenMaxTrials: 1,
	}
}

func TestCircuitBreaker_TripsOnFailuresAndRecovers(t *testing.T) {
	cb := newTestBreaker()
	fail := func() error { return errDownstream }
	ok := func() error { return nil }

	if err := cb.Execute(fail); !errors.Is(err, errDownstream) {
		t.Fatalf("expected downstream error, got %v", err)
	}
	if cb.State != StateOpen {
		t.Fatalf("expected breaker to open after failure, state=%v", cb.State)
	}

	if err := cb.Execute(ok); err == nil {
		t.Fatal("open breaker must reject calls")
	}

	time.Sleep(30 * time.Millisecond)

	if err := cb.Execute(ok); err != nil {
		t.Fatalf("half-open trial should run, got %v", err)
	}
	if cb.State != StateClose {
		t.Fatalf("expected breaker to close after successful trial, state=%v", cb.State)
	}
}

func TestCircuitBreaker_StaysClosedBelowThreshold(t *testing.T) {
	cb := newTestBreaker()
	for range 3 {
		_ = cb.Execute(func() error { return nil })
	}
	_ = cb.Execute(func() error { return errDownstream })

	if cb.State != StateClose {
		t.Fatalf("1 failure in 4 calls is below 0.5, state=%v", cb.State)
	}
}
