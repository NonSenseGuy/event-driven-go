package api

import (
	"context"
	"fmt"
	"net/http"
	"tickets/entities"

	"github.com/ThreeDotsLabs/go-event-driven/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/common/clients/payments"
)

type PaymentsServiceClient struct {
	client *clients.Clients
}

func NewPaymentsServiceClient(clients *clients.Clients) *PaymentsServiceClient {
	if clients == nil {
		panic("NewPaymentsAPI: client is nil")
	}

	return &PaymentsServiceClient{clients}
}

func (c PaymentsServiceClient) RefundPayment(ctx context.Context, refundPayment entities.PaymentRefund) error {
	resp, err := c.client.Payments.PutRefundsWithResponse(
		ctx,
		payments.PaymentRefundRequest{
			PaymentReference: refundPayment.TicketID,
			Reason:           "customer requested refund",
			DeduplicationId:  &refundPayment.IdempotencyKey,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to refund ticket: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected status code from payments: %w", err)
	}

	return nil
}
