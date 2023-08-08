package entities

type PaymentRefund struct {
	TicketID       string `json:"ticket_id,omitempty"`
	RefundReason   string `json:"refund_reason,omitempty"`
	IdempotencyKey string `json:"idempotency_key,omitempty"`
}
