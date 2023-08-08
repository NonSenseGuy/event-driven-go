package event

import (
	"context"
	"tickets/entities"
)

func (h Handler) CreateReadModelOnBookingMade(ctx context.Context, event *entities.BookingMade) error {
	return nil
}
