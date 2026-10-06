package service

import (
	"context"

	mencache "github.com/MamangRust/monolith-point-of-sale-category/cache"
	"github.com/MamangRust/monolith-point-of-sale-category/repository"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
)

type Service struct {
	CategoryQuery           CategoryQueryService
	CategoryCommand         CategoryCommandService
	CategoryStats           CategoryStatsService
	CategoryStatsById       CategoryStatsByIdService
	CategoryStatsByMerchant CategoryStatsByMerchantQueryService
}

type Deps struct {
	Ctx           context.Context
	Mencache      mencache.Mencache
	Repositories  repository.Repositories
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

func NewService(deps *Deps) *Service {
	return &Service{
		CategoryQuery: NewCategoryQueryService(&categoryQueryDeps{
			Cache:         deps.Mencache,
			CategoryQuery: deps.Repositories,
			Logger:        deps.Logger,
			Observability: deps.Observability,
		}),
		CategoryCommand: NewCategoryCommandService(&categoryCommandDeps{
			Cache:           deps.Mencache,
			CategoryQuery:   deps.Repositories,
			CategoryCommand: deps.Repositories,
			Logger:          deps.Logger,
			Observability:   deps.Observability,
		}),
		CategoryStats: NewCategoryStatsService(&categoryStatsDeps{
			Cache:                   deps.Mencache,
			CategoryStatsRepository: deps.Repositories,
			Logger:                  deps.Logger,
			Observability:           deps.Observability,
		}),
		CategoryStatsById: NewCategoryStatsByIdService(&categoryStatsByIdDeps{
			Cache:                       deps.Mencache,
			CategoryStatsByIdRepository: deps.Repositories,
			Logger:                      deps.Logger,
			Observability:               deps.Observability,
		}),
		CategoryStatsByMerchant: NewCategoryStatsByMerchantQueryService(&categoryStatsByMerchantDeps{
			Cache:                             deps.Mencache,
			CategoryStatsByMerchantRepository: deps.Repositories,
			Logger:                            deps.Logger,
			Observability:                     deps.Observability,
		}),
	}
}
