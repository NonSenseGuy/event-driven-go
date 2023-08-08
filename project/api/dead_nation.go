package api

import (
	"context"
	"fmt"
	"net/http"
	"tickets/entities"

	"github.com/ThreeDotsLabs/go-event-driven/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/common/clients/dead_nation"
)

type DeadNationAPI struct {
	client *clients.Clients
}

func NewDeadNationAPI(clients *clients.Clients) *DeadNationAPI {
	if clients == nil {
		panic("NewDeadNationServiceClient: client is nil")
	}

	return &DeadNationAPI{clients}
}

func (c DeadNationAPI) BookInDeadNation(ctx context.Context, booking entities.DeadNationBooking) error {
	resp, err := c.client.DeadNation.PostTicketBookingWithResponse(
		ctx,
		dead_nation.PostTicketBookingRequest{
			CustomerAddress: booking.CustomerEmail,
			EventId:         booking.DeadNationEventID,
			NumberOfTickets: booking.NumberOfTickets,
			BookingId:       booking.BookingID,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to book place in dead nation: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected status code from dead nation: %w", err)
	}

	return nil
}
