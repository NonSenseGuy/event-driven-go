package api

import (
	"context"
	"sync"
	"tickets/entities"
)

type PaymentsServiceMock struct {
	lock    sync.Mutex
	Refunds []entities.PaymentRefund
}

func (m *PaymentsServiceMock) RefundPayment(ctx context.Context, refundPayment entities.PaymentRefund) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.Refunds = append(m.Refunds, refundPayment)

	return nil
}
