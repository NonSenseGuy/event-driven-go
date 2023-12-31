package entities

import "github.com/google/uuid"

type Booking struct {
	BookingID       uuid.UUID `json:"booking_id" db:"booking_id"`
	ShowID          uuid.UUID `json:"show_id" db:"show_id"`
	NumberOfTickets int       `json:"number_of_tickets" db:"number_of_tickets"`
	CustomerEmail   string    `json:"customer_email" db:"customer_email"`
}

type DeadNationBooking struct {
	CustomerEmail     string    `json:"customer_email,omitempty"`
	DeadNationEventID uuid.UUID `json:"dead_nation_event_id,omitempty"`
	NumberOfTickets   int       `json:"number_of_tickets,omitempty"`
	BookingID         uuid.UUID `json:"booking_id,omitempty"`
}
