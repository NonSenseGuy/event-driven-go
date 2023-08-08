package event

import (
	"context"
	"tickets/entities"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/google/uuid"
)

type Handler struct {
	deadNationAPI       DeadNationAPI
	spreadsheetsService SpreadsheetsService
	receiptsService     ReceiptsService
	filesAPI            FilesAPI
	ticketsRepository   TicketsRepository
	showsRepository     ShowsRepository
	eventBus            *cqrs.EventBus
}

func NewHandler(
	spreadsheetsService SpreadsheetsService,
	receiptsService ReceiptsService,
	showsRepository ShowsRepository,
	ticketsRepository TicketsRepository,
	filesAPI FilesAPI,
	deadNationAPI DeadNationAPI,
	eventBus *cqrs.EventBus,
) Handler {
	if spreadsheetsService == nil {
		panic("missing spreadsheetsService")
	}
	if receiptsService == nil {
		panic("missing receiptsService")
	}
	if ticketsRepository == nil {
		panic("missing tickets repository")
	}
	if showsRepository == nil {
		panic("missing shows repository")
	}
	if filesAPI == nil {
		panic("missing files api")
	}
	if deadNationAPI == nil {
		panic("missing dead nation api")
	}
	if eventBus == nil {
		panic("missing event bus")
	}

	return Handler{
		spreadsheetsService: spreadsheetsService,
		receiptsService:     receiptsService,
		ticketsRepository:   ticketsRepository,
		showsRepository:     showsRepository,
		filesAPI:            filesAPI,
		deadNationAPI:       deadNationAPI,
		eventBus:            eventBus,
	}
}

type SpreadsheetsService interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, request entities.IssueReceiptRequest) (entities.IssueReceiptResponse, error)
}

type TicketsRepository interface {
	Add(ctx context.Context, ticket entities.Ticket) error
	Remove(ctx context.Context, ticketID string) error
}

type ShowsRepository interface {
	ShowByID(context.Context, uuid.UUID) (entities.Show, error)
}

type FilesAPI interface {
	UploadFile(ctx context.Context, fileID string, fileContent string) error
}

type DeadNationAPI interface {
	BookInDeadNation(ctx context.Context, booking entities.DeadNationBooking) error
}
