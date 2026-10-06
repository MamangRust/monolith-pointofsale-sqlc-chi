package repository

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	sharedErrors "github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"gorm.io/gorm"
)

type outboxRepository struct{}

// NewOutboxRepository creates a new instance of OutboxRepository. It holds no
// connection of its own: events are always written through the *gorm.DB
// transaction passed to EnqueueInTx so they commit atomically with the
// business write (Phase 6 — transactional outbox).
func NewOutboxRepository() OutboxRepository {
	return &outboxRepository{}
}

// EnqueueInTx persists a pending event inside the given GORM transaction so the
// caller can commit the business write and the event atomically. The payload is
// stored as raw JSON bytes in the jsonb column. The relay (pkg/outbox) then
// publishes it to Kafka with retry/dead-letter semantics.
func (r *outboxRepository) EnqueueInTx(ctx context.Context, tx *gorm.DB, topic, key string, payload []byte) error {
	event := &models.OutboxEvent{
		Topic:         topic,
		EventKey:      key,
		Payload:       payload,
		Status:        "pending",
		NextAttemptAt: time.Now(),
	}

	if err := tx.WithContext(ctx).Create(event).Error; err != nil {
		return sharedErrors.ErrInternal.WithInternal(err)
	}
	return nil
}
