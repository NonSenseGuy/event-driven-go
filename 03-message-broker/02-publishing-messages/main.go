package main

import (
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	redis "github.com/redis/go-redis/v9"
)

const topic = "progress"

func main() {
	logger := watermill.NewStdLogger(false, false)
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	publisher, err := redisstream.NewPublisher(redisstream.PublisherConfig{
		Client: rdb,
	}, logger)
	if err != nil {
		panic(err)
	}

	msg := message.NewMessage(watermill.NewUUID(), []byte("50"))
	err = publisher.Publish(topic, msg)
	if err != nil {
		panic(err)
	}

	msg = message.NewMessage(watermill.NewUUID(), []byte("100"))
	err = publisher.Publish(topic, msg)
	if err != nil {
		panic(err)
	}

}
