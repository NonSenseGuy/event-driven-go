package db

import (
	"context"
	"fmt"
	"tickets/entities"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ShowsRepository struct {
	*sqlx.DB
}

func NewShowsRepository(dbConn *sqlx.DB) *ShowsRepository {
	if dbConn == nil {
		panic("db connection is nil")
	}
	return &ShowsRepository{
		dbConn,
	}
}

func (s ShowsRepository) AddShow(ctx context.Context, show entities.Show) error {
	query := `
		INSERT INTO 
			shows (show_id, dead_nation_id, number_of_tickets,start_time, title, venue)
		VALUES (:show_id, :dead_nation_id, :number_of_tickets, :start_time, :title, :venue)
		ON CONFLICT DO NOTHING
	`

	_, err := s.NamedExecContext(ctx, query, &show)
	if err != nil {
		return fmt.Errorf("could not save show: %w", err)
	}

	return nil
}

func (s ShowsRepository) AllShows(ctx context.Context) ([]entities.Show, error) {
	var shows []entities.Show
	err := s.DB.SelectContext(ctx, &shows, `
		SELECT 
			*
		FROM 
			shows
	`)
	if err != nil {
		return nil, fmt.Errorf("could not get all shows from db: %w", err)
	}

	return shows, nil
}

func (s ShowsRepository) ShowByID(ctx context.Context, showID uuid.UUID) (entities.Show, error) {
	var show entities.Show

	err := s.DB.GetContext(ctx, &show,
		`
		SELECT 
			*
		FROM 
			shows
		WHERE 
			show_id = $1
		`,
		showID)
	if err != nil {
		return entities.Show{}, fmt.Errorf("could not get show from db: %w", err)
	}

	return show, nil
}
