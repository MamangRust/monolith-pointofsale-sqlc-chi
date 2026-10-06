package repository

import (
	pbusers "github.com/MamangRust/monolith-point-of-sale-pb/users"
	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter"
	useradapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/user"
	"gorm.io/gorm"
)

// Repositories is a struct of all merchant repositories. Uses named fields
// (not embedding) because MerchantQueryRepository and the user adapter both
// declare FindById.
type Repositories struct {
	MerchantQuery           MerchantQueryRepository
	MerchantCommand         MerchantCommandRepository
	MerchantDocumentQuery   MerchantDocumentQueryRepository
	MerchantDocumentCommand MerchantDocumentCommandRepository
	UserQuery               useradapter.QueryRepository
}

// GuardOptions collects the dependency guards for each remote dependency.
type GuardOptions struct {
	User []adapter.GuardOption
}

func NewRepositories(db *gorm.DB, userClient pbusers.UserQueryServiceClient, guards ...GuardOptions) *Repositories {
	var g GuardOptions
	if len(guards) > 0 {
		g = guards[0]
	}

	return &Repositories{
		MerchantQuery:           NewMerchantQueryRepository(db),
		MerchantCommand:         NewMerchantCommandRepository(db),
		MerchantDocumentQuery:   NewMerchantDocumentQueryRepository(db),
		MerchantDocumentCommand: NewMerchantDocumentCommandRepository(db),
		UserQuery:               useradapter.New(userClient, nil, g.User...),
	}
}
