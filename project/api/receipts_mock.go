package api

import (
	"context"
	"sync"
	"tickets/entities"
	"time"
)

type ReceiptsServiceMock struct {
	mock sync.Mutex

	IssuedReceipts []entities.IssueReceiptRequest
	VoidReceipts   []entities.VoidReceipt
}

func (m *ReceiptsServiceMock) IssueReceipt(ctx context.Context, request entities.IssueReceiptRequest) (entities.IssueReceiptResponse, error) {
	m.mock.Lock()
	defer m.mock.Unlock()

	m.IssuedReceipts = append(m.IssuedReceipts, request)

	return entities.IssueReceiptResponse{
		ReceiptNumber: "mocked-receipt-number",
		IssuedAt:      time.Now(),
	}, nil
}

func (m *ReceiptsServiceMock) VoidReceipt(ctx context.Context, request entities.VoidReceipt) error {
	m.mock.Lock()
	defer m.mock.Unlock()

	m.VoidReceipts = append(m.VoidReceipts, request)

	return nil
}
