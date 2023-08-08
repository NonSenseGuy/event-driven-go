package main

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	wsql "github.com/ThreeDotsLabs/watermill-sql/v2/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func SubscribeForMessages(db *sqlx.DB, topic string, logger watermill.LoggerAdapter) (<-chan *message.Message, error) {
	config := wsql.SubscriberConfig{
		SchemaAdapter:  wsql.DefaultPostgreSQLSchema{},
		OffsetsAdapter: wsql.DefaultPostgreSQLOffsetsAdapter{},
	}

	sub, err := wsql.NewSubscriber(db, config, logger)
	if err != nil {
		return nil, err
	}

	err = sub.SubscribeInitialize(topic)
	if err != nil {
		return nil, err
	}

	return sub.Subscribe(context.Background(), topic)
}
