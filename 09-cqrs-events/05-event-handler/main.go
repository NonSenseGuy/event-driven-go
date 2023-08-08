package main

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type FollowRequestSent struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type EventsCounter interface {
	CountEvent() error
}

type EventsCounterHandler struct {
	counter int
}

func (h *EventsCounterHandler) CountEvent() error {
	h.counter++
	return nil
}

func (h *EventsCounterHandler) CountEventOnRequestSent(ctx context.Context, event *FollowRequestSent) error {
	return h.CountEvent()
}

func NewFollowRequestSentHandler(counter EventsCounter) cqrs.EventHandler {
	return cqrs.NewEventHandler(
		"FollowRequestSent",
		func(ctx context.Context, event *FollowRequestSent) error {
			return counter.CountEvent()
		},
	)

}
