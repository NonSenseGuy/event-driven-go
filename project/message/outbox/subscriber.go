package outbox

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	wsql "github.com/ThreeDotsLabs/watermill-sql/v2/pkg/sql"
)

func NewPostgresSubscriber(db *sql.DB, logger watermill.LoggerAdapter) *wsql.Subscriber {
	sub, err := wsql.NewSubscriber(
		db,
		wsql.SubscriberConfig{
			PollInterval:     time.Millisecond * 100,
			InitializeSchema: true,
			SchemaAdapter:    wsql.DefaultPostgreSQLSchema{},
			OffsetsAdapter:   wsql.DefaultPostgreSQLOffsetsAdapter{},
		},
		logger,
	)
	if err != nil {
		panic(fmt.Errorf("failed to create new sql subscriber: %w", err))
	}

	return sub
}
