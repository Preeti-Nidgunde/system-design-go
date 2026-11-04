package notifier

import (
	"testing"
	"time"
)

func TestOnEvent(t *testing.T) {

	emailOb := NewEmailObserver([]string{"User1"})
	smsOb := NewSMSObserver([]string{"user2"})

	system := NewEventGenerator()

	system.Subscribe(emailOb)
	system.Subscribe(smsOb)

	event := EventData{
		data:   "event happened",
		sentAt: time.Now(),
	}
	err := system.OnEvent(event)
	if err != nil {
		t.Fatalf("failed")
	}
}
