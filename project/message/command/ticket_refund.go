package command

import (
	"context"
	"fmt"
	"tickets/entities"
)

func (h Handler) RefundTicket(ctx context.Context, ticketRefund *entities.RefundTicket) error {
	idempotencyKey := ticketRefund.Header.IdempotencyKey
	if idempotencyKey == "" {
		return fmt.Errorf("idempotency key is required")
	}

	err := h.receiptsServiceClient.VoidReceipt(ctx, entities.VoidReceipt{
		TicketID:       ticketRefund.TicketID,
		Reason:         "ticket refunded",
		IdempotencyKey: idempotencyKey,
	})
	if err != nil {
		return fmt.Errorf("could not void receipt: %w", err)
	}

	err = h.paymentsServiceClient.RefundPayment(ctx, entities.PaymentRefund{
		TicketID:       ticketRefund.TicketID,
		RefundReason:   "ticket refunded",
		IdempotencyKey: idempotencyKey,
	})
	if err != nil {
		return fmt.Errorf("could not refund payment: %w", err)
	}

	err = h.eventBus.Publish(ctx, entities.TicketRefunded{
		Header:   entities.NewEventHeader(),
		TicketID: ticketRefund.TicketID,
	})
	if err != nil {
		return fmt.Errorf("failed to publish TicketRefunded event: %w", err)
	}

	return nil
}
