package db

import (
	"context"
	"os"
	"sync"
	"testing"
	"tickets/entities"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

var db *sqlx.DB
var getDbOnce sync.Once

func getDb(t *testing.T) *sqlx.DB {
	getDbOnce.Do(func() {
		var err error
		db, err = sqlx.Open("postgres", os.Getenv("POSTGRES_URL"))
		if err != nil {
			panic(err)
		}

		t.Cleanup(func() {
			db.Close()
		})
	})
	return db
}

func TestIdempotency(t *testing.T) {
	db = getDb(t)

	err := InitializeDBSchema(db)
	require.NoError(t, err)

	dbRepo := NewTicketsRepository(db)

	ticket := entities.Ticket{
		TicketID: uuid.NewString(),
		Price: entities.Money{
			Amount:   "10",
			Currency: "usd",
		},
		CustomerEmail: "customer@email.com",
	}

	err = dbRepo.Add(context.Background(), ticket)
	require.NoError(t, err)

	err = dbRepo.Add(context.Background(), ticket)
	require.NoError(t, err)

	tickets, err := dbRepo.GetTickets(context.Background())
	require.NoError(t, err)
	require.Len(t, tickets, 1)
}
