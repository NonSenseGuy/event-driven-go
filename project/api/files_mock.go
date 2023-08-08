package api

import (
	"context"
	"sync"
	"tickets/entities"
)

type FilesAPIMock struct {
	mock sync.Mutex

	IssuedReceipts []entities.IssueReceiptRequest
}

func (m *FilesAPIMock) UploadFile(ctx context.Context, fileID string, fileContent string) error {
	m.mock.Lock()
	defer m.mock.Unlock()

	return nil
}
