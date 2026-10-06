package handler

import (
	pbcashiers "github.com/MamangRust/monolith-point-of-sale-pb/cashiers"
)

type CashierQueryHandleGrpc interface {
	pbcashiers.CashierQueryServiceServer
}

type CashierCommandHandleGrpc interface {
	pbcashiers.CashierCommandServiceServer
}

type CashierStatsHandleGrpc interface {
	pbcashiers.CashierStatsServiceServer
}

type CashierStatsByIdHandleGrpc interface {
	pbcashiers.CashierStatsByIdServiceServer
}

type CashierStatsByMerchantHandleGrpc interface {
	pbcashiers.CashierStatsByMerchantServiceServer
}
