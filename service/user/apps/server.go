package apps

import (
	"os"
	"time"

	pbroles "github.com/MamangRust/monolith-point-of-sale-pb/roles"
	pbuserrole "github.com/MamangRust/monolith-point-of-sale-pb/user_role"
	pbusers "github.com/MamangRust/monolith-point-of-sale-pb/users"
	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter"
	"github.com/MamangRust/monolith-point-of-sale-pkg/hash"
	"github.com/MamangRust/monolith-point-of-sale-pkg/resilience"
	"github.com/MamangRust/monolith-point-of-sale-pkg/server"
	"github.com/MamangRust/monolith-point-of-sale-shared/cache"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	mencache "github.com/MamangRust/monolith-point-of-sale-user/cache"
	"github.com/MamangRust/monolith-point-of-sale-user/handler"
	"github.com/MamangRust/monolith-point-of-sale-user/repository"
	"github.com/MamangRust/monolith-point-of-sale-user/service"
	"google.golang.org/grpc"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	// Roles and user-role assignments live behind the role service; dial it and
	// guard the adapter so a slow/unavailable dependency cannot stall the user
	// service.
	roleAddr := os.Getenv("GRPC_ROLE_ADDR")
	if roleAddr == "" {
		roleAddr = "localhost:50052"
	}

	roleConn, err := server.NewGRPCClient(roleAddr)
	if err != nil {
		return nil, err
	}

	go func() {
		<-srv.Ctx.Done()
		roleConn.Close()
	}()

	roleQueryClient := pbroles.NewRoleQueryServiceClient(roleConn)
	userRoleClient := pbuserrole.NewUserRoleServiceClient(roleConn)

	guardRole := resilience.NewDependencyGuard("role", 5, 30, 100, 3*time.Second, srv.Logger)
	guardUserRole := resilience.NewDependencyGuard("user_role", 5, 30, 100, 3*time.Second, srv.Logger)

	repos := repository.NewRepositories(&repository.Deps{
		Db:       srv.GormDB,
		Role:     roleQueryClient,
		UserRole: userRoleClient,
		Guards: repository.GuardOptions{
			Role:     []adapter.GuardOption{adapter.WithDependencyGuard(guardRole)},
			UserRole: []adapter.GuardOption{adapter.WithDependencyGuard(guardUserRole)},
		},
	})
	hash := hash.NewHashingPassword()
	traceLoggerObservability := observability.NewTraceLoggerObservability(srv.Logger)

	cacheMetrics, err := observability.NewCacheMetrics("user")
	if err != nil {
		return nil, err
	}

	cacheStore := cache.NewCacheStore(srv.Redis, srv.Logger, cacheMetrics)
	mencacheObj := mencache.NewMencache(cacheStore)

	services := service.NewService(&service.Deps{
		Mencache:      mencacheObj,
		Repositories:  repos,
		Hash:          hash,
		Logger:        srv.Logger,
		Observability: traceLoggerObservability,
	})

	handlers := handler.NewHandler(&handler.Deps{
		Service: services,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbusers.RegisterUserQueryServiceServer(gs, handlers.User)
		pbusers.RegisterUserCommandServiceServer(gs, handlers.UserCommand)
	}

	return srv, nil
}
