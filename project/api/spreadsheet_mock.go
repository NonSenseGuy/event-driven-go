package api

import (
	"context"
	"sync"
)

type SpreadsheetsServiceMock struct {
	lock sync.Mutex

	Rows map[string][][]string
}

func (m *SpreadsheetsServiceMock) AppendRow(ctx context.Context, sheetName string, row []string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.Rows == nil {
		m.Rows = make(map[string][][]string)
	}

	m.Rows[sheetName] = append(m.Rows[sheetName], row)

	return nil
}
