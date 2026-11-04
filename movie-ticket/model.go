package movieticket

type MovieCategory int
type SeatType int
type BookingStatus int

const (
	ADULT MovieCategory = iota
	CHILD
)

const (
	NORMAL SeatType = iota
	GOLD
	SILVER
	PREMIUM
)

const (
	CONFIRMED BookingStatus = iota
	CANCELLED
)

const (
	SEAT_AVAILABLE SeatStatus = iota
	SEAT_LOCKED
	SEAT_BOOKED
)

type Payment struct {
	ID, BookingID string
	Amount        int64
	Status        string
}

type ShowSeat struct {
	ID, ShowID, SeatID, UserID, BookingID string
	Status                                SeatStatus
}
