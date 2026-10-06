package handler

import (
	"github.com/MamangRust/monolith-point-of-sale-order-item/service"
	pborder_items "github.com/MamangRust/monolith-point-of-sale-pb/order_items"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	OrderItem        pborder_items.OrderItemQueryServiceServer
	OrderItemCommand pborder_items.OrderItemCommandServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		OrderItem:        NewOrderItemQueryHandleGrpc(deps.Service, deps.Logger),
		OrderItemCommand: NewOrderItemCommandHandleGrpc(deps.Service, deps.Logger),
	}
}
