package db

import (
	"context"
	"fmt"
	"tickets/entities"

	"github.com/jmoiron/sqlx"
)

type TicketsRepository struct {
	*sqlx.DB
}

func NewTicketsRepository(dbConn *sqlx.DB) *TicketsRepository {
	if dbConn == nil {
		panic("db connection is nil")
	}
	return &TicketsRepository{
		dbConn,
	}
}

func (t TicketsRepository) GetTickets(ctx context.Context) ([]entities.Ticket, error) {
	var tickets []entities.Ticket
	query := `
		SELECT 
			ticket_id, 
			price_amount as "price.amount",
			price_currency as "price.currency",
			customer_email
		from 
			tickets;
	`

	err := t.Select(&tickets, query)
	if err != nil {
		return nil, fmt.Errorf("could not get tickets from db: %w", err)
	}

	return tickets, nil
}

func (t TicketsRepository) Add(ctx context.Context, ticket entities.Ticket) error {
	query := `
        INSERT INTO tickets (ticket_id, price_amount, price_currency, customer_email)
        VALUES (:ticket_id, :price.amount, :price.currency, :customer_email)
		ON CONFLICT DO NOTHING
    `

	_, err := t.NamedExecContext(ctx, query, &ticket)
	if err != nil {
		return fmt.Errorf("could not save ticket: %w", err)
	}

	return nil
}

func (t TicketsRepository) Remove(ctx context.Context, ticketID string) error {
	query := `
		DELETE FROM tickets where ticket_id = $1
	`

	_, err := t.ExecContext(ctx, query, ticketID)
	if err != nil {
		return fmt.Errorf("could not delete ticket: %w", err)
	}

	return nil
}
