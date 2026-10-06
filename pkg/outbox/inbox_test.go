package outbox

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// dryRunDB returns a GORM instance whose SQL is built but never executed, so
// the inbox statements can be inspected without a live database.
func dryRunDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(postgres.New(postgres.Config{DSN: "host=localhost dbname=none"}),
		&gorm.Config{DryRun: true, DisableAutomaticPing: true})
	if err != nil {
		t.Fatalf("failed to open dry-run gorm: %v", err)
	}
	return db
}

func TestReserveValidatesKeys(t *testing.T) {
	for name, tc := range map[string]struct {
		db           *gorm.DB
		consumerName string
		eventKey     string
	}{
		"nil db":          {db: nil, consumerName: "email-service-group", eventKey: "topic:evt-1"},
		"empty consumer":  {db: dryRunDB(t), consumerName: "", eventKey: "topic:evt-1"},
		"empty event key": {db: dryRunDB(t), consumerName: "email-service-group", eventKey: ""},
	} {
		t.Run(name, func(t *testing.T) {
			if _, _, _, err := Reserve(context.Background(), tc.db, tc.consumerName, tc.eventKey, "topic", 0, 1); !errors.Is(err, ErrInvalidInboxKey) {
				t.Fatalf("expected ErrInvalidInboxKey, got %v", err)
			}
		})
	}
}

func TestMarkProcessedValidatesKeys(t *testing.T) {
	if err := MarkProcessed(context.Background(), nil, "email-service-group", "topic:evt-1", 1); !errors.Is(err, ErrInvalidInboxKey) {
		t.Fatalf("expected ErrInvalidInboxKey, got %v", err)
	}
	if err := MarkProcessed(context.Background(), dryRunDB(t), "", "topic:evt-1", 1); !errors.Is(err, ErrInvalidInboxKey) {
		t.Fatalf("expected ErrInvalidInboxKey, got %v", err)
	}
}

func TestReleaseValidatesKeys(t *testing.T) {
	if err := Release(context.Background(), nil, "email-service-group", "topic:evt-1", 1, nil); !errors.Is(err, ErrInvalidInboxKey) {
		t.Fatalf("expected ErrInvalidInboxKey, got %v", err)
	}
	if err := Release(context.Background(), dryRunDB(t), "", "topic:evt-1", 1, nil); !errors.Is(err, ErrInvalidInboxKey) {
		t.Fatalf("expected ErrInvalidInboxKey, got %v", err)
	}
}

func TestReserveUsesLeaseFencedUpsert(t *testing.T) {
	db := dryRunDB(t)
	sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Clauses(clause.OnConflict{
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
			"consumer_name":       "email-service-group",
			"event_key":           "topic:evt-1",
			"topic":               "topic",
			"partition_id":        0,
			"message_offset":      1,
			"status":              inboxStatusProcessing,
			"attempts":            1,
			"reservation_version": 1,
			"lease_until":         gorm.Expr(inboxLeaseInterval),
			"last_error":          "",
			"processed_at":        nil,
		})
	})
	for _, needle := range []string{"consumer_inbox", "ON CONFLICT", "lease_until <= current_timestamp"} {
		if !strings.Contains(sql, needle) {
			t.Fatalf("reserve SQL does not contain %q: %s", needle, sql)
		}
	}
}

func TestMarkProcessedFencesOnActiveReservation(t *testing.T) {
	db := dryRunDB(t)
	sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(&models.ConsumerInbox{}).
			Where("consumer_name = ? AND event_key = ? AND status = ? AND reservation_version = ?",
				"email-service-group", "topic:evt-1", inboxStatusProcessing, 1).
			Updates(map[string]interface{}{
				"status":       inboxStatusProcessed,
				"processed_at": gorm.Expr("current_timestamp"),
				"lease_until":  gorm.Expr("current_timestamp"),
				"last_error":   "",
			})
	})
	for _, needle := range []string{"consumer_inbox", "reservation_version", "status = 'processing'"} {
		if !strings.Contains(sql, needle) {
			t.Fatalf("mark processed SQL does not contain %q: %s", needle, sql)
		}
	}
}

func TestReleaseRecordsProcessingError(t *testing.T) {
	db := dryRunDB(t)
	processingErr := errors.New("smtp down")
	sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(&models.ConsumerInbox{}).
			Where("consumer_name = ? AND event_key = ? AND status = ? AND reservation_version = ?",
				"email-service-group", "topic:evt-1", inboxStatusProcessing, 1).
			Updates(map[string]interface{}{
				"status":      inboxStatusPending,
				"lease_until": gorm.Expr("current_timestamp"),
				"last_error":  processingErr.Error(),
			})
	})
	for _, needle := range []string{"consumer_inbox", "last_error", "smtp down"} {
		if !strings.Contains(sql, needle) {
			t.Fatalf("release SQL does not contain %q: %s", needle, sql)
		}
	}
}