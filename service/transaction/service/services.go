package service

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/kafka"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-pkg/outbox"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	mencache "github.com/MamangRust/monolith-point-of-sale-transacton/cache"
	"github.com/MamangRust/monolith-point-of-sale-transacton/repository"
	"gorm.io/gorm"
)

type Service struct {
	TransactionQuery           TransactionQueryService
	TransactionCommand         TransactionCommandService
	TransactionStats           TransactionStatsService
	TransactionStatsByMerchant TransactionStatsByMerchantQueryService
}

type Deps struct {
	Ctx           context.Context
	Kafka         *kafka.Kafka
	Mencache      mencache.Mencache
	Repositories  *repository.Repositories
	DB            *gorm.DB
	Outbox        *outbox.OutboxService
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

func NewService(deps *Deps) *Service {
	// Normalize a typed-nil *kafka.Kafka into a nil interface so the
	// graceful-degradation guard (s.kafka == nil) in the command service works.
	var kafkaPublisher EmailEventPublisher
	if deps.Kafka != nil {
		kafkaPublisher = deps.Kafka
	}

	return &Service{
		TransactionQuery:           NewTransactionQueryService(deps.Mencache, deps.Repositories.TransactionQuery, deps.Logger, deps.Observability),
		TransactionCommand:         NewTransactionCommandService(kafkaPublisher, deps.Mencache, deps.Repositories.CashierQuery, deps.Repositories.MerchantQuery, deps.Repositories.TransactionQuery, deps.Repositories.TransactionCommand, deps.Repositories.OrderQuery, deps.Repositories.OrderItemQuery, deps.DB, deps.Outbox, deps.Logger, deps.Observability),
		TransactionStats:           NewTransactionStatsService(deps.Mencache, deps.Repositories.TransactionStats, deps.Logger, deps.Observability),
		TransactionStatsByMerchant: NewTransactionStatsByMerchantQueryService(deps.Mencache, deps.Repositories.TransactionStatsByMerchant, deps.Logger, deps.Observability),
	}
}
