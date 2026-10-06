package handler

import (
	"github.com/MamangRust/monolith-point-of-sale-order/service"
	pborders "github.com/MamangRust/monolith-point-of-sale-pb/orders"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Order                pborders.OrderQueryServiceServer
	OrderCommand         pborders.OrderCommandServiceServer
	OrderStats           pborders.OrderStatsServiceServer
	OrderStatsById       pborders.OrderStatsByIdServiceServer
	OrderStatsByMerchant pborders.OrderStatsByMerchantServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Order:                NewOrderQueryHandleGrpc(deps.Service, deps.Logger),
		OrderCommand:         NewOrderCommandHandleGrpc(deps.Service, deps.Logger),
		OrderStats:           NewOrderStatsHandleGrpc(deps.Service, deps.Logger),
		OrderStatsById:       NewOrderStatsByIdHandleGrpc(deps.Service, deps.Logger),
		OrderStatsByMerchant: NewOrderStatsByMerchantHandleGrpc(deps.Service, deps.Logger),
	}
}
