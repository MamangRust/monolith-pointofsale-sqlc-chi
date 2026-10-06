package repository

import (
	"context"
	"errors"
)

// ErrInvalidInboxKey is returned when a consumer inbox operation is called
// with an empty consumer name or event key.
var ErrInvalidInboxKey = errors.New("invalid consumer inbox key")

// ConsumerInbox is the durable deduplication contract used by the Kafka email
// handlers. Implementations persist reservations in PostgreSQL so
// idempotency survives consumer restarts and rebalances; at-least-once
// redelivery cannot send the same email twice.
type ConsumerInbox interface {
	// Reserve claims an event for a consumer. It returns false when the event
	// was already processed. An expired processing lease may be reclaimed
	// after a consumer crashes. The third return value is the reservation
	// version that must be passed to MarkProcessed/Release.
	Reserve(ctx context.Context, consumerName, eventKey, topic string, partition int32, offset int64) (bool, bool, int64, error)
	// MarkProcessed completes only the active reservation.
	MarkProcessed(ctx context.Context, consumerName, eventKey string, reservationVersion int64) error
	// Release releases only the active reservation so the event can be retried.
	Release(ctx context.Context, consumerName, eventKey string, reservationVersion int64, processingErr error) error
}
