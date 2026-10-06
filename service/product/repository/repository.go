package repository

import (
	pbcategories "github.com/MamangRust/monolith-point-of-sale-pb/categories"
	pbmerchants "github.com/MamangRust/monolith-point-of-sale-pb/merchants"
	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter"
	categoryadapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/category"
	merchantadapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/merchant"
	"gorm.io/gorm"
)

// Repositories groups the product's own GORM repositories with the shared gRPC
// adapters for category and merchant.
type Repositories struct {
	ProductQuery   ProductQueryRepository
	ProductCommand ProductCommandRepository
	CategoryQuery  categoryadapter.Repository
	MerchantQuery  merchantadapter.Repository
}

// GuardOptions collects the dependency guards for each remote dependency.
type GuardOptions struct {
	Category []adapter.GuardOption
	Merchant []adapter.GuardOption
}

// NewRepositories creates a new instance of Repositories with the provided GORM
// database and the category/merchant gRPC clients.
func NewRepositories(
	db *gorm.DB,
	categoryClient pbcategories.CategoryQueryServiceClient,
	merchantClient pbmerchants.MerchantQueryServiceClient,
	guards ...GuardOptions,
) *Repositories {
	var g GuardOptions
	if len(guards) > 0 {
		g = guards[0]
	}

	return &Repositories{
		ProductQuery:   NewProductQueryRepository(db),
		ProductCommand: NewProductCommandRepository(db),
		CategoryQuery:  categoryadapter.New(categoryClient, g.Category...),
		MerchantQuery:  merchantadapter.New(merchantClient, g.Merchant...),
	}
}
