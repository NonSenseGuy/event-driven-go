package db

import (
	"errors"
	"fmt"
	"tickets/entities"
	"tickets/message/event"
	"tickets/message/outbox"

	"github.com/jmoiron/sqlx"
	"golang.org/x/net/context"
)

var ErrNotEnoughTicketsAvailable = errors.New("tickets sold out")

type BookingsRepository struct {
	*sqlx.DB
}

func NewBookingsRepository(dbConn *sqlx.DB) *BookingsRepository {
	if dbConn == nil {
		panic("db connection is nil")
	}

	return &BookingsRepository{
		dbConn,
	}
}

func (b BookingsRepository) AddBooking(ctx context.Context, booking entities.Booking) error {
	query := `
		INSERT INTO bookings (booking_id, show_id, number_of_tickets, customer_email)
		VALUES (:booking_id, :show_id, :number_of_tickets, :customer_email)
		ON CONFLICT DO NOTHING
	`

	tx, err := b.DB.Beginx()
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			err = errors.Join(err, rollbackErr)
			return
		}

		err = tx.Commit()
	}()

	availableSeats := 0
	err = tx.GetContext(ctx, &availableSeats, `
		SELECT 
			number_of_tickets AS available_seats
		FROM 
			shows
		WHERE 
			show_id = $1
	`, booking.ShowID)
	if err != nil {
		return fmt.Errorf("could not get available seats: %w", err)
	}

	alreadyBookedSeats := 0
	err = tx.GetContext(ctx, &alreadyBookedSeats, `
		SELECT 
			coalesce(SUM(number_of_tickets), 0) AS already_booked_seats
		FROM 
			bookings
		WHERE 
			show_id = $1
	`, booking.ShowID)
	if err != nil {
		return fmt.Errorf("could not get already booked seats: %w", err)
	}

	if availableSeats-alreadyBookedSeats < booking.NumberOfTickets {
		return ErrNotEnoughTicketsAvailable
	}

	_, err = tx.NamedExecContext(ctx, query, &booking)
	if err != nil {
		return fmt.Errorf("could not save booking %w", err)
	}

	outboxPublisher, err := outbox.NewPublisherForDB(ctx, tx)
	if err != nil {
		return fmt.Errorf("could not create event bus: %w", err)
	}

	err = event.NewEventBus(outboxPublisher).Publish(ctx, entities.BookingMade{
		Header:          entities.NewEventHeader(),
		BookingID:       booking.BookingID,
		NumberOfTickets: booking.NumberOfTickets,
		CustomerEmail:   booking.CustomerEmail,
		ShowID:          booking.ShowID,
	})
	if err != nil {
		return fmt.Errorf("could not publish event: %w", err)
	}

	return nil
}
