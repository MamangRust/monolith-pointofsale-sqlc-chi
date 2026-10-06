package outbox

import (
	"context"
	"errors"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/gormerr"
	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrInvalidInboxKey = errors.New("invalid consumer inbox key")

const (
	inboxStatusProcessing = "processing"
	inboxStatusProcessed  = "processed"
	inboxStatusPending    = "pending"

	// inboxLeaseInterval extends a reservation lease by one minute on the
	// database clock, matching the original sqlc consumer_inbox query
	// (current_timestamp + interval '1 minute').
	inboxLeaseInterval = "current_timestamp + interval '1 minute'"

	// inboxReleaseDefaultError is stored when Release receives a nil error,
	// matching the original outbox.Release behavior.
	inboxReleaseDefaultError = "consumer processing failed"
)

// ConsumerInbox is the durable deduplication contract used by Kafka handlers
// (Phase 3 — Durable Idempotency). It replaces in-memory-only deduplication:
// reservations survive consumer restarts and rebalances, so at-least-once
// redelivery cannot send the same email twice.
type ConsumerInbox interface {
	Reserve(ctx context.Context, consumerName, eventKey, topic string, partition int32, offset int64) (bool, bool, int64, error)
	MarkProcessed(ctx context.Context, consumerName, eventKey string, reservationVersion int64) error
	Release(ctx context.Context, consumerName, eventKey string, reservationVersion int64, processingErr error) error
}

// Reserve claims an event for a consumer. It returns false when the event was
// already processed. An expired processing lease may be reclaimed after a
// consumer crashes.
//
// It renders the original sqlc ReserveConsumerInbox statement with GORM:
// INSERT ... ON CONFLICT (consumer_name, event_key) DO UPDATE ... WHERE the
// existing row is neither processed nor under an active lease. RowsAffected==0
// means the conflicting row rejected the guard, so the caller learns whether
// the event was already processed and which reservation version it carries.
func Reserve(ctx context.Context, db *gorm.DB, consumerName, eventKey, topic string, partition int32, offset int64) (bool, bool, int64, error) {
	if db == nil || consumerName == "" || eventKey == "" {
		return false, false, 0, ErrInvalidInboxKey
	}

	tx := db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return false, false, 0, tx.Error
	}

	result := tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "consumer_name"}, {Name: "event_key"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"status":              inboxStatusProcessing,
			"attempts":            gorm.Expr("consumer_inbox.attempts + 1"),
			"reservation_version": gorm.Expr("consumer_inbox.reservation_version + 1"),
			"lease_until":         gorm.Expr(inboxLeaseInterval),
			"last_error":          "",
			"topic":               gorm.Expr("excluded.topic"),
			"partition_id":        gorm.Expr("excluded.partition_id"),
			"message_offset":      gorm.Expr("excluded.message_offset"),
		}),
		Where: clause.Where{Exprs: []clause.Expression{
			gorm.Expr("consumer_inbox.status <> 'processed' AND consumer_inbox.lease_until <= current_timestamp"),
		}},
	}).Model(&models.ConsumerInbox{}).Create(map[string]interface{}{
		"consumer_name":       consumerName,
		"event_key":           eventKey,
		"topic":               topic,
		"partition_id":        partition,
		"message_offset":      offset,
		"status":              inboxStatusProcessing,
		"attempts":            1,
		"reservation_version": 1,
		"lease_until":         gorm.Expr(inboxLeaseInterval),
		"last_error":          "",
		"processed_at":        nil,
	})
	if result.Error != nil {
		_ = tx.Rollback().Error
		return false, false, 0, result.Error
	}

	// The claim succeeded: either a fresh insert or an expired-lease reclaim.
	// Read back the new reservation version inside the same transaction.
	if result.RowsAffected == 1 {
		var row struct{ ReservationVersion int64 }
		err := tx.Model(&models.ConsumerInbox{}).
			Select("reservation_version").
			Where("consumer_name = ? AND event_key = ?", consumerName, eventKey).
			Take(&row).Error
		if err != nil {
			_ = tx.Rollback().Error
			return false, false, 0, err
		}
		if err := tx.Commit().Error; err != nil {
			return false, false, 0, err
		}
		return true, false, row.ReservationVersion, nil
	}

	// The conflicting row rejected the guard: it is either already processed
	// or still under an active lease. Report its state so callers dedupe or
	// redeliver instead of double-sending.
	var existing models.ConsumerInbox
	err := tx.Select("status", "reservation_version").
		Where("consumer_name = ? AND event_key = ?", consumerName, eventKey).
		Take(&existing).Error
	if err != nil {
		_ = tx.Rollback().Error
		if gormerr.IsRecordNotFound(err) {
			// Nothing to dedupe against yet; the next delivery retries cleanly.
			return false, false, 0, nil
		}
		return false, false, 0, err
	}
	if err := tx.Commit().Error; err != nil {
		return false, false, 0, err
	}
	return false, existing.Status == inboxStatusProcessed, existing.ReservationVersion, nil
}

// MarkProcessed completes only the active reservation. Like the original sqlc
// :exec query, a version mismatch is a silent no-op rather than an error.
func MarkProcessed(ctx context.Context, db *gorm.DB, consumerName, eventKey string, reservationVersion int64) error {
	if db == nil || consumerName == "" || eventKey == "" {
		return ErrInvalidInboxKey
	}
	return db.WithContext(ctx).Model(&models.ConsumerInbox{}).
		Where("consumer_name = ? AND event_key = ? AND status = ? AND reservation_version = ?",
			consumerName, eventKey, inboxStatusProcessing, reservationVersion).
		Updates(map[string]interface{}{
			"status":       inboxStatusProcessed,
			"processed_at": gorm.Expr("current_timestamp"),
			"lease_until":  gorm.Expr("current_timestamp"),
			"last_error":   "",
		}).Error
}

// Release releases only the active reservation so the retry path can re-claim
// the event. A stale version is a no-op, matching the original query.
func Release(ctx context.Context, db *gorm.DB, consumerName, eventKey string, reservationVersion int64, processingErr error) error {
	if db == nil || consumerName == "" || eventKey == "" {
		return ErrInvalidInboxKey
	}
	lastError := inboxReleaseDefaultError
	if processingErr != nil {
		lastError = processingErr.Error()
	}
	return db.WithContext(ctx).Model(&models.ConsumerInbox{}).
		Where("consumer_name = ? AND event_key = ? AND status = ? AND reservation_version = ?",
			consumerName, eventKey, inboxStatusProcessing, reservationVersion).
		Updates(map[string]interface{}{
			"status":      inboxStatusPending,
			"lease_until": gorm.Expr("current_timestamp"),
			"last_error":  lastError,
		}).Error
}

// PostgresInbox adapts a GORM connection to ConsumerInbox. Reservation and
// completion are committed independently because an external side effect cannot
// share a PostgreSQL transaction with the Kafka consumer.
type PostgresInbox struct {
	db *gorm.DB
}

func NewPostgresInbox(db *gorm.DB) (*PostgresInbox, error) {
	if db == nil {
		return nil, errors.New("inbox db is nil")
	}
	return &PostgresInbox{db: db}, nil
}

func (i *PostgresInbox) Reserve(ctx context.Context, consumerName, eventKey, topic string, partition int32, offset int64) (bool, bool, int64, error) {
	return Reserve(ctx, i.db, consumerName, eventKey, topic, partition, offset)
}

func (i *PostgresInbox) MarkProcessed(ctx context.Context, consumerName, eventKey string, reservationVersion int64) error {
	return MarkProcessed(ctx, i.db, consumerName, eventKey, reservationVersion)
}

func (i *PostgresInbox) Release(ctx context.Context, consumerName, eventKey string, reservationVersion int64, processingErr error) error {
	return Release(ctx, i.db, consumerName, eventKey, reservationVersion, processingErr)
}