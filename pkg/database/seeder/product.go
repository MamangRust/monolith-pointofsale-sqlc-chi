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

type productSeeder struct {
	db     *gorm.DB
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewProductSeeder(db *gorm.DB, ctx context.Context, logger logger.LoggerInterface) *productSeeder {
	return &productSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *productSeeder) Seed() error {
	var merchants []models.Merchant
	if err := r.db.WithContext(r.ctx).
		Order("created_at DESC").
		Limit(20).
		Offset(0).
		Find(&merchants).Error; err != nil {
		r.logger.Error("Failed to get merchants:", zap.Any("error", err))
		return err
	}

	var categories []models.Category
	if err := r.db.WithContext(r.ctx).
		Order("created_at DESC").
		Limit(20).
		Offset(0).
		Find(&categories).Error; err != nil {
		r.logger.Error("Failed to get categories:", zap.Any("error", err))
		return err
	}

	if len(merchants) == 0 || len(categories) == 0 {
		r.logger.Error("No merchants or categories found, skipping seeding")
		return nil
	}

	productNames := []string{
		"Smartphone", "Laptop", "Wireless Earbuds", "Gaming Mouse", "Mechanical Keyboard",
		"Smartwatch", "Power Bank", "Bluetooth Speaker", "External Hard Drive", "Monitor",
	}
	brands := []string{"Samsung", "Apple", "Sony", "Logitech", "Razer", "Xiaomi", "HP", "Dell", "Acer", "Asus"}
	images := []string{
		"image1.jpg", "image2.jpg", "image3.jpg", "image4.jpg", "image5.jpg",
		"image6.jpg", "image7.jpg", "image8.jpg", "image9.jpg", "image10.jpg",
	}

	for i := 0; i < 10; i++ {
		merchant := merchants[rand.Intn(len(merchants))]
		category := categories[rand.Intn(len(categories))]
		name := productNames[rand.Intn(len(productNames))]
		brand := brands[rand.Intn(len(brands))]
		price := int32(rand.Intn(5000000) + 50000)
		countInStock := int32(rand.Intn(100) + 1)
		weight := int32(rand.Intn(5000) + 100)
		slug := fmt.Sprintf("%s-%d", name, rand.Intn(1000))
		image := images[rand.Intn(len(images))]
		barcode := fmt.Sprintf("BC-%d", rand.Intn(9999999))

		product := &models.Product{
			MerchantID:   merchant.MerchantID,
			CategoryID:   category.CategoryID,
			Name:         name,
			Description:  ptrString(fmt.Sprintf("Description for %s", name)),
			Price:        price,
			CountInStock: countInStock,
			Brand:        ptrString(brand),
			Weight:       ptrInt32(weight),
			SlugProduct:  ptrString(slug),
			ImageProduct: ptrString(image),
			Barcode:      ptrString(barcode),
		}

		if err := r.db.WithContext(r.ctx).Create(product).Error; err != nil {
			r.logger.Error("Failed to create product:", zap.Error(err))
			return err
		}

	}

	r.logger.Info("Product seeding completed successfully.", zap.Int("count", 10))
	return nil
}
