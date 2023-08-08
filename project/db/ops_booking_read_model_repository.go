package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"tickets/entities"
	"time"

	"github.com/ThreeDotsLabs/go-event-driven/common/log"
	"github.com/jmoiron/sqlx"
)

type OpsBookingReadModel struct {
	db *sqlx.DB
}

func NewOpsBookingReadModel(db *sqlx.DB) *OpsBookingReadModel {
	if db == nil {
		panic("NewOpsBookingReadModel db is nil")
	}

	return &OpsBookingReadModel{db}
}

func (r OpsBookingReadModel) OnBookingMade(ctx context.Context, booking *entities.BookingMade) error {
	log.FromContext(ctx).Info("Creating ops booking")

	err := r.createReadModel(ctx, entities.OpsBooking{
		BookingID:  booking.BookingID,
		BookedAt:   booking.Header.PublishedAt,
		Tickets:    nil,
		LastUpdate: time.Now(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (r OpsBookingReadModel) createReadModel(ctx context.Context, booking entities.OpsBooking) error {
	payload, err := json.Marshal(booking)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(
		ctx,
		`
		INSERT INTO
			read_model_ops_booking (payload, booking_id)
		VALUES
			($1, $2)
		ON CONFLICT 
			(booking_id) DO NOTHING; 
		`,
		payload, booking.BookingID,
	)
	if err != nil {
		return fmt.Errorf("could not create read model: %w", err)
	}

	return nil
}

func (r OpsBookingReadModel) findReadModelByBookingID(ctx context.Context, bookingID string, tx *sqlx.Tx) (entities.OpsBooking, error) {
	var payload []byte
	err := tx.QueryRowContext(ctx, `
		SELECT payload from read_model_ops_booking WHERE booking_id = $1;
	`, bookingID).Scan(payload)
	if err != nil {
		return entities.OpsBooking{}, err
	}

	return r.unmarshalReadModelFromDB(payload)
}

func (r OpsBookingReadModel) findReadModelByTicketID(ctx context.Context, ticketID string, tx *sqlx.Tx) (entities.OpsBooking, error) {
	var payload []byte

	err := r.db.QueryRowContext(ctx, `
		SELECT payload FROM read_model_ops_booking WHERE payload::jsonb -> 'tickets' ? $1
	`, ticketID).Scan(&payload)
	if err != nil {
		return entities.OpsBooking{}, err
	}

	return r.unmarshalReadModelFromDB(payload)
}

func (r OpsBookingReadModel) updateBookingReadModel(
	ctx context.Context,
	bookingID string,
	updateFunc func(booking entities.OpsBooking) (entities.OpsBooking, error),
) error {
	return updateInTx(
		ctx,
		r.db,
		sql.LevelRepeatableRead,
		func(ctx context.Context, tx *sqlx.Tx) error {
			rm, err := r.findReadModelByBookingID(ctx, bookingID, tx)
			if err == sql.ErrNoRows {
				return fmt.Errorf("read model for booking %s does not exists yet", bookingID)
			}
			if err != nil {
				return fmt.Errorf("could not find read model: %w", err)
			}

			updatedRm, err := updateFunc(rm)
			if err != nil {
				return err
			}

			return r.updateReadModel(ctx, tx, updatedRm)
		},
	)
}

func (r OpsBookingReadModel) updateReadModel(
	ctx context.Context,
	tx *sqlx.Tx,
	rm entities.OpsBooking,
) error {
	rm.LastUpdate = time.Now()

	payload, err := json.Marshal(rm)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO
			read_model_ops_booking (payload, booking_id)
		VALUES
			($1, $2)
		ON CONFLICT (booking_id) DO UPDATE SET payload = excluded.payload;
	`, payload, rm.BookingID)
	if err != nil {
		return fmt.Errorf("could not update read model: %w", err)
	}

	return nil
}

func (r OpsBookingReadModel) unmarshalReadModelFromDB(payload []byte) (entities.OpsBooking, error) {
	var opsBooking entities.OpsBooking
	if err := json.Unmarshal(payload, &opsBooking); err != nil {
		return entities.OpsBooking{}, err
	}

	if opsBooking.Tickets == nil {
		opsBooking.Tickets = make(map[string]entities.OpsTicket)
	}

	return opsBooking, nil
}

func updateInTx(
	ctx context.Context,
	db *sqlx.DB,
	isolation sql.IsolationLevel,
	fn func(ctx context.Context, tx *sqlx.Tx) error,
) (err error) {
	tx, err := db.BeginTxx(ctx, &sql.TxOptions{Isolation: isolation})
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = errors.Join(err, rollbackErr)
			}
			return
		}

		err = tx.Commit()
	}()

	return fn(ctx, tx)
}
