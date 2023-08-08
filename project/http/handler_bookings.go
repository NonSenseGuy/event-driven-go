package http

import (
	"net/http"
	"tickets/db"
	"tickets/entities"

	"github.com/ThreeDotsLabs/go-event-driven/common/log"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type bookTicketRequest struct {
	ShowID          uuid.UUID `json:"show_id"`
	NumberOfTickets int       `json:"number_of_tickets"`
	CustomerEmail   string    `json:"customer_email"`
}

type bookTicketResponse struct {
	BookingID uuid.UUID   `json:"booking_id"`
	TicketIds []uuid.UUID `json:"ticket_ids"`
}

func (h Handler) PostBookTickets(c echo.Context) error {
	var request bookTicketRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	if request.NumberOfTickets < 1 {
		return echo.NewHTTPError(http.StatusBadRequest, "number of tickets must be greater than 0")
	}

	bookingID := uuid.New()
	booking := entities.Booking{
		BookingID:       bookingID,
		ShowID:          request.ShowID,
		NumberOfTickets: request.NumberOfTickets,
		CustomerEmail:   request.CustomerEmail,
	}

	err = h.bookingsRepository.AddBooking(c.Request().Context(), booking)
	if err == db.ErrNotEnoughTicketsAvailable {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err != nil {
		log.FromContext(c.Request().Context()).Error(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, bookTicketResponse{BookingID: bookingID})
}
