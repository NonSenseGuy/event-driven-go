package event

import (
	"context"
	"tickets/entities"

	"github.com/ThreeDotsLabs/go-event-driven/common/log"
)

func (h Handler) DeadNationPostTicketBooking(ctx context.Context, event *entities.BookingMade) error {
	log.FromContext(ctx).Info("booking ticket in dead nation")

	show, err := h.showsRepository.ShowByID(ctx, event.ShowID)
	if err != nil {
		return err
	}

	return h.deadNationAPI.BookInDeadNation(ctx, entities.DeadNationBooking{
		CustomerEmail:     event.CustomerEmail,
		DeadNationEventID: show.DeadNationID,
		NumberOfTickets:   event.NumberOfTickets,
		BookingID:         event.BookingID,
	})
}
