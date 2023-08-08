package message

import (
	"tickets/db"
	"tickets/message/command"
	"tickets/message/event"
	"tickets/message/outbox"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
)

const brokenMesageID = "2beaf5bc-d5e4-4653-b075-2b36bbf28949"

func NewWatermillRouter(
	postgresSub message.Subscriber,
	pub message.Publisher,
	eventProcessorConfig cqrs.EventProcessorConfig,
	eventHandler event.Handler,
	commandProcessorConfig cqrs.CommandProcessorConfig,
	commandHandler command.Handler,
	opsBookingReadModel db.OpsBookingReadModel,
	watermillLogger watermill.LoggerAdapter,
) *message.Router {
	router, err := message.NewRouter(message.RouterConfig{}, watermillLogger)
	if err != nil {
		panic(err)
	}

	useMiddleware(router, watermillLogger)

	outbox.AddForwarderHandler(postgresSub, pub, router, watermillLogger)

	ep, err := cqrs.NewEventProcessorWithConfig(router, eventProcessorConfig)
	if err != nil {
		panic(err)
	}

	ep.AddHandlers(
		cqrs.NewEventHandler(
			"AppendToTracker",
			eventHandler.AppendToTracker,
		),
		cqrs.NewEventHandler(
			"TicketRefundToSheet",
			eventHandler.TicketRefundToSheet,
		),
		cqrs.NewEventHandler(
			"IssueReceipt",
			eventHandler.IssueReceipt,
		),
		cqrs.NewEventHandler(
			"PrintTicketHandler",
			eventHandler.PrintTicket,
		),
		cqrs.NewEventHandler(
			"StoreTickets",
			eventHandler.StoreTicket,
		),
		cqrs.NewEventHandler(
			"RemoveCanceledTicket",
			eventHandler.RemoveTicket,
		),
		cqrs.NewEventHandler(
			"PostBookingDeadNation",
			eventHandler.DeadNationPostTicketBooking,
		),
		cqrs.NewEventHandler(
			"ops_reading_model.OnBookingMade",
			opsBookingReadModel.OnBookingMade,
		),
	)

	cp, err := cqrs.NewCommandProcessorWithConfig(router, commandProcessorConfig)
	if err != nil {
		panic(err)
	}

	cp.AddHandlers(
		cqrs.NewCommandHandler(
			"TicketRefund",
			commandHandler.RefundTicket,
		),
	)

	return router
}
