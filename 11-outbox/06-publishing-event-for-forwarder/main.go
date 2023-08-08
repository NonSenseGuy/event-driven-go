package main

import (
	"database/sql"

	"github.com/ThreeDotsLabs/watermill"
	wsql "github.com/ThreeDotsLabs/watermill-sql/v2/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/components/forwarder"
	"github.com/ThreeDotsLabs/watermill/message"
	_ "github.com/lib/pq"
)

var outboxTopic = "events_to_forward"

func PublishInTx(
	msg *message.Message,
	tx *sql.Tx,
	logger watermill.LoggerAdapter,
) error {
	pub, err := wsql.NewPublisher(tx, wsql.PublisherConfig{
		SchemaAdapter: wsql.DefaultPostgreSQLSchema{},
	}, logger)
	if err != nil {
		return err
	}

	publisher := forwarder.NewPublisher(pub, forwarder.PublisherConfig{
		ForwarderTopic: outboxTopic,
	})

	return publisher.Publish("ItemAddedToCart", msg)
}
