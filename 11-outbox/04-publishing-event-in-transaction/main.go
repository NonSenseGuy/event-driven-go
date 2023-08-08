package main

import (
	"database/sql"

	"github.com/ThreeDotsLabs/watermill"
	wsql "github.com/ThreeDotsLabs/watermill-sql/v2/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/message"
	_ "github.com/lib/pq"
)

func PublishInTx(
	message *message.Message,
	tx *sql.Tx,
	logger watermill.LoggerAdapter,
) error {
	pub, err := wsql.NewPublisher(
		tx,
		wsql.PublisherConfig{
			SchemaAdapter: wsql.DefaultPostgreSQLSchema{},
		},
		logger,
	)
	if err != nil {
		return err
	}

	pub.Publish("ItemAddedToCart", message)
	return nil
}
