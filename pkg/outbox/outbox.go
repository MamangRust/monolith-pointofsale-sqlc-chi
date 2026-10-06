package outbox

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Outbox constants control the durable retry behavior of the outbox relay
// (Phase 6 — Transactional Outbox). They mirror the ecommerce transaction
// service so behavior is uniform across repositories.
const (
	OutboxMaxAttempts    = 5
	OutboxBackoff        = 30 * time.Second
	OutboxRelayInterval  = 5 * time.Second
	OutboxRelayBatchSize = 100
	// OutboxClaimLease is how long a relay worker owns a claimed event. If the
	// worker dies after claiming but before marking the event delivered, the
	// lease expires and another relay instance re-claims and retries it.
	OutboxClaimLease = 1 * time.Minute
	// OutboxRetention is how long delivered/dead events are kept before the
	// relay purges them as part of the retention policy.
	OutboxRetention = 7 * 24 * time.Hour

	// OutboxRetentionEveryTicks runs the retention purge every N relay ticks so
	// the DELETE scan does not run on every relay cycle.
	OutboxRetentionEveryTicks = 60
)

// OutboxPublisher is the minimal Kafka producer surface the relay needs.
// *kafka.Kafka satisfies this interface via SendMessage.
type OutboxPublisher interface {
	SendMessage(ctx context.Context, topic string, key string, value []byte) error
}

// OutboxService persists email events durably and relays them to Kafka with
// retry and dead-letter semantics. Producers enqueue the event inside the same
// database transaction as the business write (EnqueueInTx) so a crash between
// the two cannot lose the event; the relay then guarantees delivery.
type OutboxService struct {
	db        *gorm.DB
	publisher OutboxPublisher
	logger    logger.LoggerInterface
}

// NewOutboxService builds the outbox service. The publisher may be nil (e.g.
// local dev without Kafka) — the relay then drains the queue without sending.
func NewOutboxService(db *gorm.DB, publisher OutboxPublisher, log logger.LoggerInterface) *OutboxService {
	return &OutboxService{db: db, publisher: publisher, logger: log}
}

// EnqueueInTx persists a pending event inside the given database transaction so
// the caller can commit the business write and the event atomically. This is the
// production path: the event survives the commit and is published by the relay.
func (s *OutboxService) EnqueueInTx(ctx context.Context, tx *gorm.DB, topic, key string, payload []byte) error {
	event := &models.OutboxEvent{
		Topic:         topic,
		EventKey:      key,
		Payload:       payload,
		Status:        "pending",
		NextAttemptAt: time.Now(),
	}
	if err := tx.WithContext(ctx).Create(event).Error; err != nil {
		return err
	}
	s.logger.Info("outbox event enqueued", zap.String("topic", topic), zap.String("key", key))
	return nil
}

// Enqueue persists a pending event AFTER the business transaction has already
// committed. It is the NON-ATOMIC fallback path: a crash between the commit and
// this insert silently loses the event, so it must not be used where the
// business write is local. It exists for aggregator services whose business
// write happens in another service over gRPC (best-effort guarantee).
func (s *OutboxService) Enqueue(ctx context.Context, topic, key string, payload []byte) error {
	if s.db == nil {
		return nil
	}
	event := &models.OutboxEvent{
		Topic:         topic,
		EventKey:      key,
		Payload:       payload,
		Status:        "pending",
		NextAttemptAt: time.Now(),
	}
	if err := s.db.WithContext(ctx).Create(event).Error; err != nil {
		return err
	}
	s.logger.Info("outbox event enqueued", zap.String("topic", topic), zap.String("key", key))
	return nil
}

// PublishPending claims up to limit pending events whose retry window has
// elapsed, publishes each to Kafka, and marks it delivered. Claiming is atomic
// (FOR UPDATE SKIP LOCKED + lease), so concurrent relay instances never publish
// the same event twice. It returns the number of events successfully delivered.
func (s *OutboxService) PublishPending(ctx context.Context, limit int) (int, error) {
	if s.db == nil || s.publisher == nil {
		return 0, nil
	}

	// Claim up to limit due events inside a transaction: the SELECT ... FOR
	// UPDATE SKIP LOCKED makes concurrent relays skip rows another instance is
	// claiming, and extending next_attempt_at to the lease expiry means a crash
	// between claim and delivery is retried once the lease passes.
	var events []models.OutboxEvent
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Where("status = 'pending' AND next_attempt_at <= ?", time.Now()).
			Order("outbox_id ASC").
			Limit(limit).
			Find(&events).Error; err != nil {
			return err
		}
		if len(events) == 0 {
			return nil
		}
		ids := make([]int64, len(events))
		for i, event := range events {
			ids[i] = event.OutboxID
		}
		return tx.Model(&models.OutboxEvent{}).
			Where("outbox_id IN ? AND status = 'pending'", ids).
			Update("next_attempt_at", time.Now().Add(OutboxClaimLease)).Error
	})
	if err != nil {
		return 0, err
	}

	delivered := 0
	for _, event := range events {
		if err := s.publisher.SendMessage(ctx, event.Topic, event.EventKey, event.Payload); err != nil {
			s.logger.Error("failed to publish outbox event, scheduling retry",
				zap.Error(err),
				zap.Int64("outbox_id", event.OutboxID),
				zap.String("topic", event.Topic),
				zap.Int32("attempts", event.Attempts),
			)
			if int(event.Attempts)+1 >= OutboxMaxAttempts {
				if err := s.db.WithContext(ctx).Model(&models.OutboxEvent{}).
					Where("outbox_id = ? AND status = 'pending'", event.OutboxID).
					Update("status", "dead").Error; err != nil {
					s.logger.Error("failed to dead-letter outbox event", zap.Error(err), zap.Int64("outbox_id", event.OutboxID))
				}
				continue
			}
			nextAttempt := time.Now().Add(OutboxBackoff * time.Duration(event.Attempts+1))
			if err := s.db.WithContext(ctx).Model(&models.OutboxEvent{}).
				Where("outbox_id = ? AND status = 'pending'", event.OutboxID).
				Updates(map[string]interface{}{
					"attempts":        event.Attempts + 1,
					"next_attempt_at": nextAttempt,
				}).Error; err != nil {
				s.logger.Error("failed to record outbox retry", zap.Error(err), zap.Int64("outbox_id", event.OutboxID))
			}
			continue
		}
		if err := s.db.WithContext(ctx).Model(&models.OutboxEvent{}).
			Where("outbox_id = ? AND status = 'pending'", event.OutboxID).
			Update("status", "delivered").Error; err != nil {
			s.logger.Error("failed to mark outbox event delivered", zap.Error(err), zap.Int64("outbox_id", event.OutboxID))
			continue
		}
		delivered++
	}
	return delivered, nil
}

// Start runs the relay loop until ctx is cancelled.
func (s *OutboxService) Start(ctx context.Context, interval time.Duration, limit int) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	tickCount := 0
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("outbox relay stopped")
			return
		case <-ticker.C:
			if _, err := s.PublishPending(ctx, limit); err != nil {
				s.logger.Error("outbox relay cycle failed", zap.Error(err))
			}
			tickCount++
			// Retention runs periodically (not every tick) to avoid scanning the
			// outbox table on every relay cycle; it purges delivered/dead events
			// whose terminal state is older than the retention window.
			if tickCount%OutboxRetentionEveryTicks == 0 {
				result := s.db.WithContext(ctx).
					Where("status IN ('delivered', 'dead') AND updated_at < ?", time.Now().Add(-OutboxRetention)).
					Delete(&models.OutboxEvent{})
				if result.Error != nil {
					s.logger.Error("outbox retention cleanup failed", zap.Error(result.Error))
				} else if result.RowsAffected > 0 {
					s.logger.Info("outbox retention cleanup", zap.Int64("removed", result.RowsAffected))
				}
			}
		}
	}
}