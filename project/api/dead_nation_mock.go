package api

import (
	"context"
	"sync"
	"tickets/entities"
)

type DeadNationAPIMock struct {
	lock sync.Mutex

	DeadNationBookings []entities.DeadNationBooking
}

func (c *DeadNationAPIMock) BookInDeadNation(ctx context.Context, request entities.DeadNationBooking) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.DeadNationBookings = append(c.DeadNationBookings, request)

	return nil
}
