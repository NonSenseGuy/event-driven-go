package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ThreeDotsLabs/go-event-driven/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/common/clients/spreadsheets"
)

type SpreadsheetServiceClient struct {
	clients *clients.Clients
}

func NewSpreadsheetServiceClient(clients *clients.Clients) *SpreadsheetServiceClient {
	if clients == nil {
		panic("NewSpreadsheetServiceClient: clients is nil")
	}

	return &SpreadsheetServiceClient{clients: clients}
}

func (c SpreadsheetServiceClient) AppendRow(ctx context.Context, sheetName string, row []string) error {
	resp, err := c.clients.Spreadsheets.PostSheetsSheetRowsWithResponse(ctx, sheetName, spreadsheets.PostSheetsSheetRowsJSONRequestBody{
		Columns: row,
	})
	if err != nil {
		return fmt.Errorf("failed to post row to spreadsheet %v: %w", sheetName, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("failed to post row: unexpected status code %d", resp.StatusCode())
	}

	return nil
}
