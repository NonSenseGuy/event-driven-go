package service

import (
	"context"
	"fmt"
	"net/http"
	"tickets/db"
	libHttp "tickets/http"
	"tickets/message"
	"tickets/message/command"
	"tickets/message/event"
	"tickets/message/outbox"

	"github.com/ThreeDotsLabs/go-event-driven/common/log"
	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func init() {
	log.Init(logrus.InfoLevel)
}

type Service struct {
	dbConn          *sqlx.DB
	watermillRouter *watermillMessage.Router
	echoRouter      *echo.Echo
}

type ReceiptService interface {
	event.ReceiptsService
	command.ReceiptsService
}

func New(
	dbConn *sqlx.DB,
	redisClient *redis.Client,
	spreadsheetsService event.SpreadsheetsService,
	receiptsService ReceiptService,
	filesAPI event.FilesAPI,
	deadNationAPI event.DeadNationAPI,
	paymentsService command.PaymentsService,
) Service {
	ticketsRepo := db.NewTicketsRepository(dbConn)
	showsRepo := db.NewShowsRepository(dbConn)
	bookingRepo := db.NewBookingsRepository(dbConn)
	readModelRepo := db.NewOpsBookingReadModel(dbConn)

	watermillLogger := log.NewWatermill(log.FromContext(context.Background()))

	var redisPublisher watermillMessage.Publisher
	redisPublisher = message.NewRedisPublisher(redisClient, watermillLogger)
	redisPublisher = log.CorrelationPublisherDecorator{Publisher: redisPublisher}

	eventBus := event.NewEventBus(redisPublisher)

	eventsHandler := event.NewHandler(
		spreadsheetsService,
		receiptsService,
		showsRepo,
		ticketsRepo,
		filesAPI,
		deadNationAPI,
		eventBus,
	)

	commandHandler := command.NewHandler(eventBus, receiptsService, paymentsService)
	commandBus := command.NewCommandBus(redisPublisher, command.NewBusConfig(watermillLogger))

	postgresSubscriber := outbox.NewPostgresSubscriber(dbConn.DB, watermillLogger)
	eventProcessorConfig := event.NewEventProcessorConfig(redisClient, watermillLogger)
	commandProcessorConfig := command.NewProcessorConfig(redisClient, watermillLogger)

	watermillRouter := message.NewWatermillRouter(
		postgresSubscriber,
		redisPublisher,
		eventProcessorConfig,
		eventsHandler,
		commandProcessorConfig,
		commandHandler,
		watermillLogger,
	)

	echoRouter := libHttp.NewHttpRouter(
		eventBus,
		commandBus,
		spreadsheetsService,
		ticketsRepo,
		showsRepo,
		bookingRepo,
	)

	return Service{
		dbConn,
		watermillRouter,
		echoRouter,
	}
}

func (s Service) Run(ctx context.Context) error {
	if err := db.InitializeDBSchema(s.dbConn); err != nil {
		return fmt.Errorf("failed to initialize database schema: %w", err)
	}
	errgrp, ctx := errgroup.WithContext(ctx)

	errgrp.Go(func() error {
		return s.watermillRouter.Run(ctx)
	})

	errgrp.Go(func() error {
		<-s.watermillRouter.Running()

		err := s.echoRouter.Start(":8080")
		if err != nil && err != http.ErrServerClosed {
			return err
		}

		return nil
	})

	errgrp.Go(func() error {
		<-ctx.Done()
		return s.echoRouter.Shutdown(context.Background())
	})

	return errgrp.Wait()
}
