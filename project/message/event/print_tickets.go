package event

import (
	"context"
	"fmt"
	"tickets/entities"

	"github.com/ThreeDotsLabs/go-event-driven/common/log"
)

func (h Handler) PrintTicket(ctx context.Context, event *entities.TicketBookingConfirmed) error {
	log.FromContext(ctx).Info("printing ticket")

	ticketHTML := `
		<html>
			<head>
				<title> Ticket </title>
			</head>
			<body>
				<h1> Ticket: ` + event.TicketID + `</h1>
				<p> Price: ` + event.Price.Amount + ` ` + event.Price.Currency + `</p>
			</body>
		</html>
	`

	ticketsFile := event.TicketID + "-ticket.html"

	err := h.filesAPI.UploadFile(ctx, ticketsFile, ticketHTML)
	if err != nil {
		return fmt.Errorf("failed to upload ticket file: %w", err)
	}

	ticketPrintedEvent := entities.TicketPrinted{
		Header:   event.Header,
		TicketID: event.TicketID,
		FileName: ticketsFile,
	}

	return h.eventBus.Publish(ctx, ticketPrintedEvent)
}
