package handler

import (
	"github.com/MamangRust/monolith-point-of-sale-category/service"
	pbcategories "github.com/MamangRust/monolith-point-of-sale-pb/categories"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
)

type Deps struct {
	Service *service.Service
	Logger  logger.LoggerInterface
}

type Handler struct {
	Category                pbcategories.CategoryQueryServiceServer
	CategoryCommand         pbcategories.CategoryCommandServiceServer
	CategoryStats           pbcategories.CategoryStatsServiceServer
	CategoryStatsById       pbcategories.CategoryStatsByIdServiceServer
	CategoryStatsByMerchant pbcategories.CategoryStatsByMerchantServiceServer
}

func NewHandler(deps *Deps) *Handler {
	return &Handler{
		Category:                NewCategoryQueryHandleGrpc(deps.Service, deps.Logger),
		CategoryCommand:         NewCategoryCommandHandleGrpc(deps.Service, deps.Logger),
		CategoryStats:           NewCategoryStatsHandleGrpc(deps.Service, deps.Logger),
		CategoryStatsById:       NewCategoryStatsByIdHandleGrpc(deps.Service, deps.Logger),
		CategoryStatsByMerchant: NewCategoryStatsByMerchantHandleGrpc(deps.Service, deps.Logger),
	}
}
