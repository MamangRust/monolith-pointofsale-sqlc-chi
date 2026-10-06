package handler

import (
	pborders "github.com/MamangRust/monolith-point-of-sale-pb/orders"
)

type OrderQueryHandleGrpc interface {
	pborders.OrderQueryServiceServer
}

type OrderCommandHandleGrpc interface {
	pborders.OrderCommandServiceServer
}

type OrderStatsHandleGrpc interface {
	pborders.OrderStatsServiceServer
}

type OrderStatsByIdHandleGrpc interface {
	pborders.OrderStatsByIdServiceServer
}

type OrderStatsByMerchantHandleGrpc interface {
	pborders.OrderStatsByMerchantServiceServer
}
