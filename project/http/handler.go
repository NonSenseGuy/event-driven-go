package http

import (
	"context"
	"tickets/entities"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/google/uuid"
)

type Handler struct {
	eventBus              *cqrs.EventBus
	commandBus            *cqrs.CommandBus
	spreadsheetsAPIClient SpreadsheetsAPI
	ticketsRepository     TicketsRepository
	showsRepository       ShowsRepository
	bookingsRepository    BookingsRepository
}

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, spreadsheetName string, row []string) error
}

type TicketsRepository interface {
	GetTickets(ctx context.Context) ([]entities.Ticket, error)
}

type ShowsRepository interface {
	AddShow(context.Context, entities.Show) error
	ShowByID(context.Context, uuid.UUID) (entities.Show, error)
	AllShows(context.Context) ([]entities.Show, error)
}

type BookingsRepository interface {
	AddBooking(context.Context, entities.Booking) error
}
