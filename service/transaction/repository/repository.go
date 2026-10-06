package repository

import (
	cashieradapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/cashier"
	merchantadapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/merchant"
	orderadapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/order"
	orderitemadapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/orderitem"

	pbcashiers "github.com/MamangRust/monolith-point-of-sale-pb/cashiers"
	pbmerchants "github.com/MamangRust/monolith-point-of-sale-pb/merchants"
	pborder_items "github.com/MamangRust/monolith-point-of-sale-pb/order_items"
	pborders "github.com/MamangRust/monolith-point-of-sale-pb/orders"
	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter"
	"gorm.io/gorm"
)

// Repositories bundles the transaction's own GORM repositories with the shared
// gRPC adapters used to reach cashier, merchant, order and order item.
type Repositories struct {
	CashierQuery               cashieradapter.QueryRepository
	MerchantQuery              merchantadapter.Repository
	OrderItemQuery             orderitemadapter.QueryRepository
	OrderQuery                 orderadapter.QueryRepository
	TransactionCommand         TransactionCommandRepository
	TransactionQuery           TransactionQueryRepository
	TransactionStatsByMerchant TransactionStatsByMerchantRepository
	TransactionStats           TransactionStatsRepository
}

// GuardOptions collects the dependency guards for each remote dependency.
type GuardOptions struct {
	Cashier   []adapter.GuardOption
	Merchant  []adapter.GuardOption
	Order     []adapter.GuardOption
	OrderItem []adapter.GuardOption
}

func NewRepositories(
	db *gorm.DB,
	cashierClient pbcashiers.CashierQueryServiceClient,
	merchantClient pbmerchants.MerchantQueryServiceClient,
	orderClient pborders.OrderQueryServiceClient,
	orderItemClient pborder_items.OrderItemQueryServiceClient,
	guards ...GuardOptions,
) *Repositories {
	var g GuardOptions
	if len(guards) > 0 {
		g = guards[0]
	}

	return &Repositories{
		CashierQuery:               cashieradapter.New(cashierClient, g.Cashier...),
		MerchantQuery:              merchantadapter.New(merchantClient, g.Merchant...),
		OrderItemQuery:             orderitemadapter.New(orderItemClient, nil, g.OrderItem...),
		OrderQuery:                 orderadapter.New(orderClient, g.Order...),
		TransactionCommand:         NewTransactionCommandRepository(db),
		TransactionQuery:           NewTransactionQueryRepository(db),
		TransactionStatsByMerchant: NewTransactionStatsByMerchantRepository(db),
		TransactionStats:           NewTransactionStatsRepository(db),
	}
}
