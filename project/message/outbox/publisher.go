package outbox

import (
	"context"

	"github.com/ThreeDotsLabs/go-event-driven/common/log"
	wsql "github.com/ThreeDotsLabs/watermill-sql/v2/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/components/forwarder"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jmoiron/sqlx"
)

func NewPublisherForDB(ctx context.Context, db *sqlx.Tx) (message.Publisher, error) {
	var publisher message.Publisher

	logger := log.NewWatermill(log.FromContext(ctx))

	publisher, err := wsql.NewPublisher(
		db,
		wsql.PublisherConfig{
			SchemaAdapter: wsql.DefaultPostgreSQLSchema{},
		},
		logger,
	)
	if err != nil {
		return nil, err
	}
	publisher = forwarder.NewPublisher(publisher, forwarder.PublisherConfig{
		ForwarderTopic: outboxTopic,
	})
	publisher = log.CorrelationPublisherDecorator{publisher}

	return publisher, nil
}
