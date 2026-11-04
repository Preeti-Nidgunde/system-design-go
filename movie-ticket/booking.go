package movieticket

import "github.com/google/uuid"

type Booking struct {
	ID, ShowID, UserID, PaymentID string
	Status                        BookingStatus
}

func NewBooking(showID, userID, paymentID string, status BookingStatus) *Booking {
	return &Booking{
		ID:        uuid.NewString(),
		ShowID:    showID,
		UserID:    userID,
		PaymentID: paymentID,
		Status:    status,
	}
}
