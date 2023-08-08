package entities

import "time"

type VoidReceipt struct {
	TicketID       string
	Reason         string
	IdempotencyKey string
}

type TicketReceiptIssued struct {
	Header        EventHeader `json:"header,omitempty"`
	TicketID      string      `json:"ticket_id,omitempty"`
	ReceiptNumber string      `json:"receipt_number,omitempty"`
	IssuedAt      time.Time   `json:"issued_at,omitempty"`
}

type IssueReceiptRequest struct {
	TicketID       string `json:"ticket_id"`
	Price          Money  `json:"price"`
	IdempotencyKey string `json:"idempotency_key"`
}

type IssueReceiptResponse struct {
	ReceiptNumber string    `json:"number"`
	IssuedAt      time.Time `json:"issued_at"`
}
