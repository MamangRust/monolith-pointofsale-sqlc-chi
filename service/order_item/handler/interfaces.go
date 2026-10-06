package handler

import (
	pborder_items "github.com/MamangRust/monolith-point-of-sale-pb/order_items"
)

type OrderItemQueryHandleGrpc interface {
	pborder_items.OrderItemQueryServiceServer
}

type OrderItemCommandHandleGrpc interface {
	pborder_items.OrderItemCommandServiceServer
}
