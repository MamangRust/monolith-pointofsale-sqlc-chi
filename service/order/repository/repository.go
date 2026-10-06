package repository

import (
	cashieradapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/cashier"
	merchantadapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/merchant"
	orderitemadapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/orderitem"
	productadapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/product"

	pbcashiers "github.com/MamangRust/monolith-point-of-sale-pb/cashiers"
	pbmerchants "github.com/MamangRust/monolith-point-of-sale-pb/merchants"
	pborder_items "github.com/MamangRust/monolith-point-of-sale-pb/order_items"
	pbproducts "github.com/MamangRust/monolith-point-of-sale-pb/products"
	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter"
	"gorm.io/gorm"
)

// Repositories bundles the order's own GORM repositories with the shared gRPC
// adapters used to reach the other services (cashier, merchant, product and
// order item).
type Repositories struct {
	OrderQuery          OrderQueryRepository
	OrderCommand        OrderCommandRepository
	OrderStats          OrderStatsRepository
	OrderStatByMerchant OrderStatByMerchantRepository

	CashierQuery     cashieradapter.QueryRepository
	MerchantQuery    merchantadapter.Repository
	ProductQuery     productadapter.QueryRepository
	ProductCommand   productadapter.CommandRepository
	OrderItemQuery   orderitemadapter.QueryRepository
	OrderItemCommand orderitemadapter.CommandRepository
}

// GuardOptions collects the dependency guards for each remote dependency.
type GuardOptions struct {
	Cashier   []adapter.GuardOption
	Merchant  []adapter.GuardOption
	Product   []adapter.GuardOption
	OrderItem []adapter.GuardOption
}

func NewRepositories(
	db *gorm.DB,
	cashierClient pbcashiers.CashierQueryServiceClient,
	merchantClient pbmerchants.MerchantQueryServiceClient,
	productQueryClient pbproducts.ProductQueryServiceClient,
	productCommandClient pbproducts.ProductCommandServiceClient,
	orderItemQueryClient pborder_items.OrderItemQueryServiceClient,
	orderItemCommandClient pborder_items.OrderItemCommandServiceClient,
	guards ...GuardOptions,
) *Repositories {
	var g GuardOptions
	if len(guards) > 0 {
		g = guards[0]
	}

	// One product adapter instance backs both the query and command roles.
	productRepo := productadapter.New(productQueryClient, productCommandClient, g.Product...)
	// One order-item adapter instance backs both the query and command roles.
	orderItemRepo := orderitemadapter.New(orderItemQueryClient, orderItemCommandClient, g.OrderItem...)

	return &Repositories{
		OrderQuery:          NewOrderQueryRepository(db),
		OrderCommand:        NewOrderCommandRepository(db),
		OrderStats:          NewOrderStatsRepository(db),
		OrderStatByMerchant: NewOrderStatsByMerchantRepository(db),

		CashierQuery:     cashieradapter.New(cashierClient, g.Cashier...),
		MerchantQuery:    merchantadapter.New(merchantClient, g.Merchant...),
		ProductQuery:     productRepo,
		ProductCommand:   productRepo,
		OrderItemQuery:   orderItemRepo,
		OrderItemCommand: orderItemRepo,
	}
}
