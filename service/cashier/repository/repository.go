package repository

import (
	pbmerchants "github.com/MamangRust/monolith-point-of-sale-pb/merchants"
	pbusers "github.com/MamangRust/monolith-point-of-sale-pb/users"
	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter"
	merchantadapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/merchant"
	useradapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/user"
	"gorm.io/gorm"
)

// Repositories is a combined struct covering the cashier's own GORM
// repositories plus the shared gRPC adapters for user and merchant. It uses
// named fields (not embedding) because multiple sub-interfaces declare
// FindById.
type Repositories struct {
	UserQuery              useradapter.QueryRepository
	MerchantQuery          merchantadapter.Repository
	CashierQuery           CashierQueryRepository
	CashierCommand         CashierCommandRepository
	CashierStats           CashierStatsRepository
	CashierStatsByMerchant CashierStatByMerchantRepository
	CashierStatsById       CashierStatByIdRepository
}

// GuardOptions collects the dependency guards for each remote dependency.
type GuardOptions struct {
	User     []adapter.GuardOption
	Merchant []adapter.GuardOption
}

func NewRepositories(
	db *gorm.DB,
	userClient pbusers.UserQueryServiceClient,
	merchantClient pbmerchants.MerchantQueryServiceClient,
	guards ...GuardOptions,
) *Repositories {
	var g GuardOptions
	if len(guards) > 0 {
		g = guards[0]
	}

	return &Repositories{
		UserQuery:              useradapter.New(userClient, nil, g.User...),
		MerchantQuery:          merchantadapter.New(merchantClient, g.Merchant...),
		CashierQuery:           NewCashierQueryRepository(db),
		CashierCommand:         NewCashierCommandRepository(db),
		CashierStats:           NewCashierStatsRepository(db),
		CashierStatsByMerchant: NewCashierStatsByMerchantRepository(db),
		CashierStatsById:       NewCashierStatsByIdRepository(db),
	}
}
