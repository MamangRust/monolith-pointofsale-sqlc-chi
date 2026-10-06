package handler

import (
	pbcategories "github.com/MamangRust/monolith-point-of-sale-pb/categories"
)

type CategoryQueryHandleGrpc interface {
	pbcategories.CategoryQueryServiceServer
}

type CategoryCommandHandleGrpc interface {
	pbcategories.CategoryCommandServiceServer
}

type CategoryStatsHandleGrpc interface {
	pbcategories.CategoryStatsServiceServer
}

type CategoryStatsByIdHandleGrpc interface {
	pbcategories.CategoryStatsByIdServiceServer
}

type CategoryStatsByMerchantHandleGrpc interface {
	pbcategories.CategoryStatsByMerchantServiceServer
}
