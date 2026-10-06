package seeder

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"

	"go.uber.org/zap"
	"golang.org/x/exp/rand"
	"gorm.io/gorm"
)

type orderSeeder struct {
	db     *gorm.DB
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewOrderSeeder(db *gorm.DB, ctx context.Context, logger logger.LoggerInterface) *orderSeeder {
	return &orderSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *orderSeeder) Seed() error {
	var merchants []models.Merchant
	if err := r.db.WithContext(r.ctx).
		Order("created_at DESC").
		Limit(20).
		Offset(0).
		Find(&merchants).Error; err != nil {
		r.logger.Error("Failed to get merchants", zap.Error(err))
		return err
	}

	var cashiers []models.Cashier
	if err := r.db.WithContext(r.ctx).
		Order("created_at DESC").
		Limit(20).
		Offset(0).
		Find(&cashiers).Error; err != nil {
		r.logger.Error("Failed to get cashiers", zap.Error(err))
		return err
	}

	if len(merchants) == 0 || len(cashiers) == 0 {
		r.logger.Error("No merchants or cashiers found, skipping order seeding")
		return nil
	}

	for i := 0; i < 10; i++ {
		merchant := merchants[rand.Intn(len(merchants))]
		cashier := cashiers[rand.Intn(len(cashiers))]
		totalPrice := int32(rand.Intn(500000) + 50000)

		order := &models.Order{
			MerchantID: merchant.MerchantID,
			CashierID:  cashier.CashierID,
			TotalPrice: int64(totalPrice),
		}

		if err := r.db.WithContext(r.ctx).Create(order).Error; err != nil {
			r.logger.Error("Failed to create order", zap.Error(err))
			return err
		}

		orderID := order.OrderID

		var products []models.Product
		if err := r.db.WithContext(r.ctx).
			Where("merchant_id = ?", merchant.MerchantID).
			Order("created_at DESC").
			Limit(10).
			Offset(0).
			Find(&products).Error; err != nil {
			r.logger.Error("Failed to get products", zap.Error(err))
			return err
		}

		if len(products) == 0 {
			r.logger.Debug("No products found for merchant", zap.Int32("merchant_id", merchant.MerchantID))
			continue
		}

		for j := 0; j < rand.Intn(5)+1; j++ {
			product := products[rand.Intn(len(products))]
			quantity := int32(rand.Intn(5) + 1)
			price := product.Price * quantity

			orderItem := &models.OrderItem{
				OrderID:   orderID,
				ProductID: product.ProductID,
				Quantity:  quantity,
				Price:     price,
			}

			if err := r.db.WithContext(r.ctx).Create(orderItem).Error; err != nil {
				r.logger.Error("Failed to create order item", zap.Error(err))
				return err
			}
		}
	}

	r.logger.Info("Order seeding completed successfully.", zap.Int("count", 10))
	return nil
}
