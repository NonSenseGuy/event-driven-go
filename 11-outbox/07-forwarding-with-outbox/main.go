package main

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	wsql "github.com/ThreeDotsLabs/watermill-sql/v2/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/components/forwarder"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	_ "github.com/lib/pq"
)

func RunForwarder(
	db *sqlx.DB,
	rdb *redis.Client,
	outboxTopic string,
	logger watermill.LoggerAdapter,
) error {
	sub, err := wsql.NewSubscriber(db, wsql.SubscriberConfig{
		SchemaAdapter:    wsql.DefaultPostgreSQLSchema{},
		OffsetsAdapter:   wsql.DefaultPostgreSQLOffsetsAdapter{},
		InitializeSchema: true,
	}, logger)
	if err != nil {
		return err
	}

	if err = sub.SubscribeInitialize(outboxTopic); err != nil {
		return err
	}

	_, err = sub.Subscribe(context.Background(), outboxTopic)
	if err != nil {
		return err
	}

	pub, err := redisstream.NewPublisher(redisstream.PublisherConfig{
		Client: rdb,
	}, logger)
	if err != nil {
		return err
	}

	fwd, err := forwarder.NewForwarder(sub, pub, logger, forwarder.Config{
		ForwarderTopic: outboxTopic,
	})
	if err != nil {
		return err
	}

	go func() {
		err := fwd.Run(context.Background())
		if err != nil {
			panic(err)
		}
	}()

	<-fwd.Running()

	return nil
}
