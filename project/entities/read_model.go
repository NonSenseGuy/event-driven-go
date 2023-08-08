package entities

import (
	"time"

	"github.com/google/uuid"
)

type OpsBooking struct {
	BookingID uuid.UUID `json:"booking_id,omitempty"`
	BookedAt  time.Time `json:"booked_at,omitempty"`

	Tickets map[string]OpsTicket `json:"tickets,omitempty"`

	LastUpdate time.Time `json:"last_update,omitempty"`
}

type OpsTicket struct {
	PriceAmount   string `json:"price_amount,omitempty"`
	PriceCurrency string `json:"price_currency,omitempty"`
	CustomerEmail string `json:"customer_email,omitempty"`

	Status string `json:"status,omitempty"`

	PrintedAt       time.Time `json:"printed_at,omitempty"`
	PrintedFileName string    `json:"printed_file_name,omitempty"`

	ReceiptIssuedAt time.Time `json:"receipt_issued_at,omitempty"`
	ReceiptNumber   string    `json:"receipt_number,omitempty"`
}
