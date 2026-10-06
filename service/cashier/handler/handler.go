package handler

import (
	"github.com/MamangRust/monolith-point-of-sale-cashier/service"
	pbcashiers "github.com/MamangRust/monolith-point-of-sale-pb/cashiers"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Cashier                pbcashiers.CashierQueryServiceServer
	CashierCommand         pbcashiers.CashierCommandServiceServer
	CashierStats           pbcashiers.CashierStatsServiceServer
	CashierStatsById       pbcashiers.CashierStatsByIdServiceServer
	CashierStatsByMerchant pbcashiers.CashierStatsByMerchantServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Cashier:                NewCashierQueryHandleGrpc(deps.Service, deps.Logger),
		CashierCommand:         NewCashierCommandHandleGrpc(deps.Service, deps.Logger),
		CashierStats:           NewCashierStatsHandleGrpc(deps.Service, deps.Logger),
		CashierStatsById:       NewCashierStatsByIdHandleGrpc(deps.Service, deps.Logger),
		CashierStatsByMerchant: NewCashierStatsByMerchantHandleGrpc(deps.Service, deps.Logger),
	}
}
