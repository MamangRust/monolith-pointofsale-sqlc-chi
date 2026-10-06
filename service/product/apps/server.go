package apps

import (
	"context"
	"os"
	"time"

	pbcategories "github.com/MamangRust/monolith-point-of-sale-pb/categories"
	pbmerchants "github.com/MamangRust/monolith-point-of-sale-pb/merchants"
	pbproducts "github.com/MamangRust/monolith-point-of-sale-pb/products"
	"google.golang.org/grpc"

	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter"
	"github.com/MamangRust/monolith-point-of-sale-pkg/resilience"
	"github.com/MamangRust/monolith-point-of-sale-pkg/server"
	mencache "github.com/MamangRust/monolith-point-of-sale-product/cache"
	"github.com/MamangRust/monolith-point-of-sale-product/handler"
	"github.com/MamangRust/monolith-point-of-sale-product/repository"
	"github.com/MamangRust/monolith-point-of-sale-product/service"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	categoryAddr := getEnv("GRPC_CATEGORY_ADDR", "localhost:50054")
	merchantAddr := getEnv("GRPC_MERCHANT_ADDR", "localhost:50056")

	categoryConn, err := server.NewGRPCClient(categoryAddr)
	if err != nil {
		return nil, err
	}

	merchantConn, err := server.NewGRPCClient(merchantAddr)
	if err != nil {
		categoryConn.Close()
		return nil, err
	}

	go func() {
		<-srv.Ctx.Done()
		categoryConn.Close()
		merchantConn.Close()
	}()

	categoryClient := pbcategories.NewCategoryQueryServiceClient(categoryConn)
	merchantClient := pbmerchants.NewMerchantQueryServiceClient(merchantConn)

	guardCategory := resilience.NewDependencyGuard("category", 5, 30, 100, 3*time.Second, srv.Logger)
	guardMerchant := resilience.NewDependencyGuard("merchant", 5, 30, 100, 3*time.Second, srv.Logger)

	repos := repository.NewRepositories(srv.GormDB, categoryClient, merchantClient,
		repository.GuardOptions{
			Category: []adapter.GuardOption{adapter.WithDependencyGuard(guardCategory)},
			Merchant: []adapter.GuardOption{adapter.WithDependencyGuard(guardMerchant)},
		},
	)

	mencacheObj := mencache.NewMencache(srv.CacheStore)

	obs := observability.NewTraceLoggerObservability(srv.Logger)

	services := service.NewService(&service.Deps{
		Mencache:      mencacheObj,
		Ctx:           context.Background(),
		Repositories:  repos,
		Logger:        srv.Logger,
		Observability: obs,
	})

	handlers := handler.NewHandler(&handler.Deps{
		Service: services,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbproducts.RegisterProductQueryServiceServer(gs, handlers.Product)
	}

	return srv, nil
}
