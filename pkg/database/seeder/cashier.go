package seeder

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"

	"go.uber.org/zap"
	"golang.org/x/exp/rand"
	"gorm.io/gorm"
)

type cashierSeeder struct {
	db     *gorm.DB
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewCashierSeeder(db *gorm.DB, ctx context.Context, logger logger.LoggerInterface) *cashierSeeder {
	return &cashierSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *cashierSeeder) Seed() error {
	var merchants []models.Merchant
	if err := r.db.WithContext(r.ctx).
		Order("created_at DESC").
		Limit(20).
		Offset(0).
		Find(&merchants).Error; err != nil {
		r.logger.Error("Failed to fetch merchants:", zap.Any("err", err))
		return err
	}

	var users []models.User
	if err := r.db.WithContext(r.ctx).
		Order("created_at DESC").
		Limit(20).
		Offset(0).
		Find(&users).Error; err != nil {
		r.logger.Error("Failed to fetch users:", zap.Any("error", err))
		return err
	}

	if len(merchants) == 0 || len(users) == 0 {
		r.logger.Error("Merchants or Users not found. Seed operation aborted.")
		return fmt.Errorf("no merchants or users found")
	}

	for i := 1; i <= 10; i++ {
		merchant := merchants[rand.Intn(len(merchants))]
		user := users[rand.Intn(len(users))]

		cashierName := fmt.Sprintf("Cashier %d", i)

		cashier := &models.Cashier{
			MerchantID: merchant.MerchantID,
			UserID:     user.UserID,
			Name:       cashierName,
		}

		if err := r.db.WithContext(r.ctx).Create(cashier).Error; err != nil {
			r.logger.Error("Failed to create cashier:", zap.Any("error", err))
			return err
		}

	}

	r.logger.Info("Cashier seeding completed successfully.", zap.Int("count", 10))
	return nil
}
