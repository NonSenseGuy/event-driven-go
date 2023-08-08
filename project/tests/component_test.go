package tests_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"tickets/api"
	"tickets/entities"
	"tickets/message"
	"tickets/service"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lithammer/shortuuid/v3"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TicketStatusRequest struct {
	Tickets []TicketStatus `json:"tickets,omitempty"`
}

type TicketStatus struct {
	TicketID  string         `json:"ticket_id,omitempty"`
	Status    string         `json:"status,omitempty"`
	Price     entities.Money `json:"price,omitempty"`
	Email     string         `json:"email,omitempty"`
	BookingID string         `json:"booking_id,omitempty"`
}

func TestComponent(t *testing.T) {
	db, err := sqlx.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err)
	}
	redisClient := message.NewRedisClient(os.Getenv("REDIS_ADDR"))
	defer redisClient.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	spreadsheetServiceMock := &api.SpreadsheetsServiceMock{}
	receiptsServiceMock := &api.ReceiptsServiceMock{}
	filesAPIMock := &api.FilesAPIMock{}
	deadNationAPIMock := &api.DeadNationAPIMock{}
	paymentsServiceMock := &api.PaymentsServiceMock{}

	go func() {
		svc := service.New(
			db,
			redisClient,
			spreadsheetServiceMock,
			receiptsServiceMock,
			filesAPIMock,
			deadNationAPIMock,
			paymentsServiceMock,
		)

		assert.NoError(t, svc.Run(ctx))
	}()

	waitForHttpServer(t)

	ticket := TicketStatus{
		TicketID: uuid.NewString(),
		Status:   "confirmed",
		Price: entities.Money{
			Amount:   "20",
			Currency: "USD",
		},
		Email:     "email@example.com",
		BookingID: uuid.NewString(),
	}

	sendTicketStatus(t, TicketStatusRequest{Tickets: []TicketStatus{ticket}})
	assertReceiptForTicketSend(t, receiptsServiceMock, ticket)
	assertRowAddedToSheet(t, spreadsheetServiceMock, ticket, "tickets-to-print")

	ticket.Status = "canceled"
	sendTicketStatus(t, TicketStatusRequest{Tickets: []TicketStatus{ticket}})
	assertRowAddedToSheet(t, spreadsheetServiceMock, ticket, "tickets-to-refund")
}

func waitForHttpServer(t *testing.T) {
	t.Helper()

	require.EventuallyWithT(
		t,
		func(t *assert.CollectT) {
			resp, err := http.Get("http://localhost:8080/health")
			if !assert.NoError(t, err) {
				return
			}
			defer resp.Body.Close()

			if assert.Less(t, resp.StatusCode, 300, "API not ready, http status: %d", resp.StatusCode) {
				return
			}
		},
		time.Second*10,
		time.Millisecond*50,
	)
}

func sendTicketStatus(t *testing.T, req TicketStatusRequest) {
	t.Helper()

	payload, err := json.Marshal(req)
	require.NoError(t, err)

	correlationID := shortuuid.New()

	ticketIDs := make([]string, 0, len(req.Tickets))
	for _, ticket := range req.Tickets {
		ticketIDs = append(ticketIDs, ticket.TicketID)
	}

	httpReq, err := http.NewRequest(
		http.MethodPost,
		"http://localhost:8080/tickets-status",
		bytes.NewBuffer(payload),
	)
	require.NoError(t, err)

	httpReq.Header.Set("correlation-id", correlationID)
	httpReq.Header.Set("content-type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func assertReceiptForTicketSend(t *testing.T, receiptsService *api.ReceiptsServiceMock, ticket TicketStatus) {
	t.Helper()

	assert.EventuallyWithT(
		t,
		func(collect *assert.CollectT) {
			issuedReceipts := len(receiptsService.IssuedReceipts)
			t.Log("issued_receipts", issuedReceipts)

			assert.Greater(collect, issuedReceipts, 0, "no receipts issued")
		},
		10*time.Second,
		100*time.Millisecond,
	)

	receipt, ok := lo.Find(receiptsService.IssuedReceipts, func(r entities.IssueReceiptRequest) bool {
		return r.TicketID == ticket.TicketID
	})
	require.Truef(t, ok, "receipt for ticket %d is not found", ticket.TicketID)

	assert.Equal(t, ticket.TicketID, receipt.TicketID)
	assert.Equal(t, ticket.Price.Amount, receipt.Price.Amount)
	assert.Equal(t, ticket.Price.Currency, receipt.Price.Currency)
}

func assertRowAddedToSheet(t *testing.T, spreadsheetService *api.SpreadsheetsServiceMock, ticket TicketStatus, sheetName string) {
	t.Helper()

	assert.EventuallyWithT(
		t,
		func(collect *assert.CollectT) {
			rows, ok := spreadsheetService.Rows[sheetName]
			if !assert.True(t, ok, "sheet %s not found", sheetName) {
				return
			}

			allValues := []string{}

			for _, row := range rows {
				for _, col := range row {
					allValues = append(allValues, col)
				}
			}

			assert.Contains(t, allValues, ticket.TicketID, "ticket id %s not found in sheet", ticket.TicketID)
		},
		10*time.Second,
		100*time.Millisecond,
	)
}
