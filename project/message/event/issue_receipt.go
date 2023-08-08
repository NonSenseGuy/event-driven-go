package event

import (
	"context"
	"fmt"
	"tickets/entities"

	"github.com/ThreeDotsLabs/go-event-driven/common/log"
)

func (h Handler) IssueReceipt(ctx context.Context, event *entities.TicketBookingConfirmed) error {
	log.FromContext(ctx).Info("issuing receipt")

	request := entities.IssueReceiptRequest{
		TicketID: event.TicketID,
		Price:    event.Price,
	}

	response, err := h.receiptsService.IssueReceipt(ctx, request)
	if err != nil {
		return fmt.Errorf("failed to issue receipt: %w", err)
	}

	err = h.eventBus.Publish(ctx, entities.TicketReceiptIssued{
		Header:        entities.NewEventHeaderWithIdempotencyKey(event.Header.IdempotencyKey),
		TicketID:      event.TicketID,
		IssuedAt:      response.IssuedAt,
		ReceiptNumber: response.ReceiptNumber,
	})
	if err != nil {
		return fmt.Errorf("failed to publish TicketReceiptIssued: %w", err)
	}

	return nil
}
