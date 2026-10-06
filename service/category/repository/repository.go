package repository

import "gorm.io/gorm"

// Repositories is an interface containing category query, command and stats repositories.
type Repositories interface {
	CategoryQueryRepository
	CategoryCommandRepository
	CategoryStatsRepository
	CategoryStatsByIdRepository
	CategoryStatsByMerchantRepository
}

type repositories struct {
	CategoryQueryRepository
	CategoryCommandRepository
	CategoryStatsRepository
	CategoryStatsByIdRepository
	CategoryStatsByMerchantRepository
}

// NewRepositories creates a new instance of Repositories with the provided GORM database.
func NewRepositories(db *gorm.DB) Repositories {
	return &repositories{
		CategoryQueryRepository:           NewCategoryQueryRepository(db),
		CategoryCommandRepository:         NewCategoryCommandRepository(db),
		CategoryStatsRepository:           NewCategoryStatsRepository(db),
		CategoryStatsByIdRepository:       NewCategoryStatsByIdRepository(db),
		CategoryStatsByMerchantRepository: NewCategoryStatsByMerchantRepository(db),
	}
}
