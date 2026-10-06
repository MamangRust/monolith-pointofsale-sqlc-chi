package handler

import (
	pbtransactions "github.com/MamangRust/monolith-point-of-sale-pb/transactions"
)

type TransactionQueryHandleGrpc interface {
	pbtransactions.TransactionQueryServiceServer
}

type TransactionCommandHandleGrpc interface {
	pbtransactions.TransactionCommandServiceServer
}

type TransactionStatsHandleGrpc interface {
	pbtransactions.TransactionStatsServiceServer
}

type TransactionStatsByMerchantHandleGrpc interface {
	pbtransactions.TransactionStatsByMerchantServiceServer
}
