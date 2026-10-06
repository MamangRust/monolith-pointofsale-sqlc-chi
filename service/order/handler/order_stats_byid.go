package handler

import (
	"github.com/MamangRust/monolith-point-of-sale-order/service"
	pborders "github.com/MamangRust/monolith-point-of-sale-pb/orders"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
)

type orderStatsByIdHandleGrpc struct {
	pborders.UnimplementedOrderStatsByIdServiceServer
	orderStats service.OrderStatsService
	logger     logger.LoggerInterface
}

func NewOrderStatsByIdHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pborders.OrderStatsByIdServiceServer {
	return &orderStatsByIdHandleGrpc{
		orderStats: service.OrderStats,
		logger:     logger,
	}
}
