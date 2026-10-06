package seeder

import (
	"context"
	"math/rand"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type transactionSeeder struct {
	db     *gorm.DB
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewTransactionSeeder(db *gorm.DB, ctx context.Context, logger logger.LoggerInterface) *transactionSeeder {
	return &transactionSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *transactionSeeder) Seed() error {
	var orders []models.Order
	if err := r.db.WithContext(r.ctx).
		Order("created_at DESC").
		Limit(20).
		Offset(0).
		Find(&orders).Error; err != nil {
		r.logger.Error("Failed to get transactions:", zap.Any("error", err))
		return err
	}

	var merchants []models.Merchant
	if err := r.db.WithContext(r.ctx).
		Order("created_at DESC").
		Limit(20).
		Offset(0).
		Find(&merchants).Error; err != nil {
		r.logger.Error("Failed to get transactions:", zap.Any("error", err))
		return err
	}

	for i := 0; i < 10; i++ {
		selectedMerchantId := merchants[rand.Intn(len(merchants))]
		selectedOrderId := orders[rand.Intn(len(orders))]

		var paymentMethod string
		var amount, changeAmount float64
		var paymentStatus string

		paymentMethod = "Credit Card"
		amount = float64(100 + i)
		changeAmount = float64(5 + i)
		paymentStatus = "Completed"

		transaction := &models.Transaction{
			OrderID:       selectedOrderId.OrderID,
			PaymentMethod: paymentMethod,
			Amount:        int32(amount),
			ChangeAmount:  ptrInt32(int32(changeAmount)),
			PaymentStatus: paymentStatus,
			MerchantID:    selectedMerchantId.MerchantID,
		}

		if err := r.db.WithContext(r.ctx).Create(transaction).Error; err != nil {
			r.logger.Error("Failed to create transaction:", zap.Any("error", err))
			return err
		}
	}

	r.logger.Info("Successfully seeded 10 transactions.", zap.Int("count", 10))
	return nil
}
