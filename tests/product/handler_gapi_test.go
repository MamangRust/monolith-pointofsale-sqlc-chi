package product_test

import (
	"context"
	pbcategories "github.com/MamangRust/monolith-point-of-sale-pb/categories"
	pbmerchants "github.com/MamangRust/monolith-point-of-sale-pb/merchants"
	pbproducts "github.com/MamangRust/monolith-point-of-sale-pb/products"
	"testing"

	prod_cache "github.com/MamangRust/monolith-point-of-sale-product/cache"
	prod_handler "github.com/MamangRust/monolith-point-of-sale-product/handler"
	prod_repo "github.com/MamangRust/monolith-point-of-sale-product/repository"
	prod_service "github.com/MamangRust/monolith-point-of-sale-product/service"
	"github.com/MamangRust/monolith-point-of-sale-shared/cache"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	tests "github.com/MamangRust/monolith-point-of-sale-test"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ProductGapiTestSuite struct {
	tests.BaseTestSuite
	client        pbproducts.ProductQueryServiceClient
	clientCommand pbproducts.ProductCommandServiceClient
}

func (s *ProductGapiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	// Setup dependencies
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()

	// Infrastructure
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)
	gormDB, err := s.GormDB()
	s.Require().NoError(err)

	// Product dependencies
	mencache := prod_cache.NewMencache(cacheStore)
	repos := prod_repo.NewRepositories(
		gormDB,
		pbcategories.NewCategoryQueryServiceClient(s.Conns["category"]),
		pbmerchants.NewMerchantQueryServiceClient(s.Conns["merchant"]),
	)
	svc := prod_service.NewService(&prod_service.Deps{
		Mencache:      mencache,
		Repositories:  repos,
		Logger:        s.Log,
		Observability: s.Obs,
		Ctx:           context.Background(),
	})

	// Handler
	handler := prod_handler.NewHandler(&prod_handler.Deps{
		Service: svc,
		Logger:  s.Log,
	})

	// Server
	server := grpc.NewServer()
	pbproducts.RegisterProductQueryServiceServer(server, handler.Product)
	pbproducts.RegisterProductCommandServiceServer(server, handler.ProductCommand)

	addr := s.RegisterServer(server)
	conn := s.GetConnection(addr)

	s.client = pbproducts.NewProductQueryServiceClient(conn)
	s.clientCommand = pbproducts.NewProductCommandServiceClient(conn)
}

func (s *ProductGapiTestSuite) TestProductGapiLifecycle() {
	ctx := context.Background()

	// 1. Seed dependencies
	userID := s.SeedUser(ctx)
	catID := s.SeedCategory(ctx)
	merchID := s.SeedMerchant(ctx, userID)

	// 2. Create
	createRes, err := s.clientCommand.Create(ctx, &pbproducts.CreateProductRequest{
		MerchantId:   int32(merchID),
		CategoryId:   int32(catID),
		Name:         "GAPI Item",
		Description:  "GAPI Description",
		Price:        1000,
		CountInStock: 10,
		Brand:        "GAPI Brand",
		Weight:       100,
		ImageProduct: "gapi.jpg",
	})
	s.Require().NoError(err)
	s.Require().NotNil(createRes)
	prodID := createRes.Data.Id

	// 3. FindById
	getRes, err := s.client.FindById(ctx, &pbproducts.FindByIdProductRequest{Id: prodID})
	s.Require().NoError(err)
	s.Equal("GAPI Item", getRes.Data.Name)

	// 4. FindAll
	allRes, err := s.client.FindAll(ctx, &pbproducts.FindAllProductRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(allRes.Data)

	// 5. FindByActive
	activeRes, err := s.client.FindByActive(ctx, &pbproducts.FindAllProductRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(activeRes.Data)

	// 6. Update
	updateRes, err := s.clientCommand.Update(ctx, &pbproducts.UpdateProductRequest{
		ProductId:    prodID,
		MerchantId:   int32(merchID),
		CategoryId:   int32(catID),
		Name:         "GAPI Item Updated",
		Description:  "Updated Description",
		Price:        2000,
		CountInStock: 20,
		Brand:        "Updated Brand",
		Weight:       200,
		ImageProduct: "gapi-updated.jpg",
	})
	s.Require().NoError(err)
	s.Equal("GAPI Item Updated", updateRes.Data.Name)

	// 7. Trash
	_, err = s.clientCommand.TrashedProduct(ctx, &pbproducts.FindByIdProductRequest{Id: prodID})
	s.Require().NoError(err)

	// 8. FindByTrashed
	trashedRes, err := s.client.FindByTrashed(ctx, &pbproducts.FindAllProductRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(trashedRes.Data)

	// 9. Restore
	_, err = s.clientCommand.RestoreProduct(ctx, &pbproducts.FindByIdProductRequest{Id: prodID})
	s.Require().NoError(err)

	// 10. DeletePermanent
	_, _ = s.clientCommand.TrashedProduct(ctx, &pbproducts.FindByIdProductRequest{Id: prodID})
	_, err = s.clientCommand.DeleteProductPermanent(ctx, &pbproducts.FindByIdProductRequest{Id: prodID})
	s.Require().NoError(err)

	// 11. RestoreAll
	_, err = s.clientCommand.RestoreAllProduct(ctx, &emptypb.Empty{})
	s.Require().NoError(err)

	// 12. DeleteAll
	_, err = s.clientCommand.DeleteAllProductPermanent(ctx, &emptypb.Empty{})
	s.Require().NoError(err)
}

func TestProductGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ProductGapiTestSuite))
}
