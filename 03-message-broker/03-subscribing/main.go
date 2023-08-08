package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/redis/go-redis/v9"
)

func main() {
	logger := watermill.NewStdLogger(false, false)
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	suscriber, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client: rdb,
	}, logger)
	if err != nil {
		panic(err)
	}

	messages, err := suscriber.Subscribe(context.Background(), "progress")
	if err != nil {
		panic(err)
	}

	for m := range messages {
		id := string(m.UUID)
		progress := string(m.Payload)
		fmt.Println(fmt.Sprintf("Message ID: %v - %v%", id, progress))
		m.Ack()
	}
}
