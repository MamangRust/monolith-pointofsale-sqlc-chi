package handler

import (
	pbtransactions "github.com/MamangRust/monolith-point-of-sale-pb/transactions"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-transacton/service"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Transaction                pbtransactions.TransactionQueryServiceServer
	TransactionCommand         pbtransactions.TransactionCommandServiceServer
	TransactionStats           pbtransactions.TransactionStatsServiceServer
	TransactionStatsByMerchant pbtransactions.TransactionStatsByMerchantServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Transaction:                NewTransactionQueryHandleGrpc(deps.Service, deps.Logger),
		TransactionCommand:         NewTransactionCommandHandleGrpc(deps.Service, deps.Logger),
		TransactionStats:           NewTransactionStatsHandleGrpc(deps.Service, deps.Logger),
		TransactionStatsByMerchant: NewTransactionStatsByMerchantHandleGrpc(deps.Service, deps.Logger),
	}
}
