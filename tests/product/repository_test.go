package product_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-product/repository"
	"github.com/MamangRust/monolith-point-of-sale-shared/cache"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	tests "github.com/MamangRust/monolith-point-of-sale-test"

	"github.com/stretchr/testify/suite"
)

type ProductRepositoryTestSuite struct {
	suite.Suite
	ts        *tests.TestSuite
	repo      *repository.Repositories
	productID int
}

func (s *ProductRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	gormDB, err := s.ts.GormDB()
	s.Require().NoError(err)

	log, _ := logger.NewLogger("test", nil)
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(ts.RedisClient(), log, cacheMetrics)
	obs, _ := observability.NewObservability("test", log)

	depsClients := tests.StartDependencyClients(gormDB, cacheStore, log, obs)
	s.repo = repository.NewRepositories(gormDB, depsClients.CategoryQuery, depsClients.MerchantQuery)

	// Seed a merchant and category for product tests
	var userID, categoryID int
	err = gormDB.WithContext(s.ts.Ctx).Raw(`INSERT INTO users (firstname, lastname, email, password, verification_code, is_verified) VALUES ($1, $2, $3, $4, 'test-verify', true) RETURNING user_id`,
		"Prod", "Repo", "prod.repo@example.com", "password123",
	).Scan(&userID).Error
	s.Require().NoError(err)

	err = gormDB.WithContext(s.ts.Ctx).Raw(`INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING category_id`,
		"Test Category", "Category for product tests",
	).Scan(&categoryID).Error
	s.Require().NoError(err)

	err = gormDB.WithContext(s.ts.Ctx).Raw(`INSERT INTO merchants (user_id, name, description, address, contact_email, contact_phone, status) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING merchant_id`,
		userID, "Test Merchant", "Desc", "Addr", "pm@example.com", "123", "active",
	).Scan(&userID).Error
	s.Require().NoError(err)
}

func (s *ProductRepositoryTestSuite) TearDownSuite() {
	tests.StopDependencyClients()
	s.ts.Teardown()
}

func (s *ProductRepositoryTestSuite) Test1_CreateProduct() {
	ctx := context.Background()

	req := &requests.CreateProductRequest{
		MerchantID:   1,
		CategoryID:   1,
		Name:         "Test Product",
		Description:  "Product description",
		Price:        100,
		CountInStock: 50,
		Brand:        "Test Brand",
		Weight:       1,
		ImageProduct: "test.jpg",
	}

	res, err := s.repo.ProductCommand.CreateProduct(ctx, req)
	s.NoError(err)
	s.NotNil(res)
	s.Equal(req.Name, res.Name)
	s.productID = int(res.ProductID)
}

func (s *ProductRepositoryTestSuite) Test2_FindById() {
	s.Require().NotZero(s.productID)
	ctx := context.Background()

	found, err := s.repo.ProductQuery.FindById(ctx, s.productID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(s.productID, int(found.ProductID))
}

func (s *ProductRepositoryTestSuite) Test3_UpdateProduct() {
	s.Require().NotZero(s.productID)
	ctx := context.Background()

	req := &requests.UpdateProductRequest{
		ProductID:    &s.productID,
		MerchantID:   1,
		CategoryID:   1,
		Name:         "Updated Product",
		Description:  "Updated description",
		Price:        200,
		CountInStock: 100,
		Brand:        "Updated Brand",
		Weight:       2,
		ImageProduct: "updated.jpg",
	}

	res, err := s.repo.ProductCommand.UpdateProduct(ctx, req)
	s.NoError(err)
	s.NotNil(res)
	s.Equal("Updated Product", res.Name)
}

func (s *ProductRepositoryTestSuite) Test4_TrashAndRestore() {
	s.Require().NotZero(s.productID)
	ctx := context.Background()

	// Trash
	trashed, err := s.repo.ProductCommand.TrashedProduct(ctx, s.productID)
	s.NoError(err)
	s.NotNil(trashed)

	// Restore
	restored, err := s.repo.ProductCommand.RestoreProduct(ctx, s.productID)
	s.NoError(err)
	s.NotNil(restored)

	// Verify restored
	found, err := s.repo.ProductQuery.FindById(ctx, s.productID)
	s.NoError(err)
	s.NotNil(found)
}

func (s *ProductRepositoryTestSuite) Test5_DeletePermanent() {
	s.Require().NotZero(s.productID)
	ctx := context.Background()

	// Must be trashed first for permanent delete
	_, err := s.repo.ProductCommand.TrashedProduct(ctx, s.productID)
	s.NoError(err)

	success, err := s.repo.ProductCommand.DeleteProductPermanent(ctx, s.productID)
	s.NoError(err)
	s.True(success)
}

func TestProductRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ProductRepositoryTestSuite))
}
