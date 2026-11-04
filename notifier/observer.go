package notifier

import (
	"context"
	"fmt"
)

type Observers interface {
	Update(ctx context.Context, data EventData) error
}

type SMSObserver struct {
	recipients []string
}

type EmailObserver struct {
	recipients []string
}

func NewSMSObserver(rec []string) *SMSObserver {
	return &SMSObserver{
		recipients: rec,
	}
}

func NewEmailObserver(rec []string) *EmailObserver {
	return &EmailObserver{
		recipients: rec,
	}
}

func (so *SMSObserver) Update(ctx context.Context, data EventData) error {
	for _, client := range so.recipients {
		fmt.Printf("Sending to client %s: %s", client, data.data)
	}
	return nil
}

func (eo *EmailObserver) Update(ctx context.Context, data EventData) error {
	for _, client := range eo.recipients {
		fmt.Printf("Sending EMAIL to  client %s body %s", client, data.data)
	}
	return nil
}
