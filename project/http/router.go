package http

import (
	"net/http"

	libHttp "github.com/ThreeDotsLabs/go-event-driven/common/http"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/labstack/echo/v4"
)

func NewHttpRouter(
	eventBus *cqrs.EventBus,
	commandBus *cqrs.CommandBus,
	spreadsheetsAPIClient SpreadsheetsAPI,
	ticketsRepository TicketsRepository,
	showsRepository ShowsRepository,
	bookingsRepository BookingsRepository,
) *echo.Echo {
	e := libHttp.NewEcho()

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	handler := Handler{
		eventBus:              eventBus,
		commandBus:            commandBus,
		spreadsheetsAPIClient: spreadsheetsAPIClient,
		ticketsRepository:     ticketsRepository,
		showsRepository:       showsRepository,
		bookingsRepository:    bookingsRepository,
	}

	e.POST("/tickets-status", handler.PostTicketStatus)
	e.PUT("/ticket-refund/:ticket_id", handler.PutTicketRefund)
	e.GET("/tickets", handler.GetTickets)

	e.POST("/shows", handler.PostShow)

	e.POST("/book-tickets", handler.PostBookTickets)

	return e
}
