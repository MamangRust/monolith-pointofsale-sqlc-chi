package apps

import (
	pbroles "github.com/MamangRust/monolith-point-of-sale-pb/roles"
	pbuserrole "github.com/MamangRust/monolith-point-of-sale-pb/user_role"
	"github.com/MamangRust/monolith-point-of-sale-pkg/server"
	mencache "github.com/MamangRust/monolith-point-of-sale-role/cache"
	"github.com/MamangRust/monolith-point-of-sale-role/handler"
	"github.com/MamangRust/monolith-point-of-sale-role/repository"
	"github.com/MamangRust/monolith-point-of-sale-role/service"
	"github.com/MamangRust/monolith-point-of-sale-shared/cache"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	"google.golang.org/grpc"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	repos := repository.NewRepositories(srv.GormDB)
	traceLoggerObservability := observability.NewTraceLoggerObservability(srv.Logger)

	cacheMetrics, err := observability.NewCacheMetrics("role")
	if err != nil {
		return nil, err
	}

	cacheStore := cache.NewCacheStore(srv.Redis, srv.Logger, cacheMetrics)
	mencacheObj := mencache.NewMencache(cacheStore)

	services := service.NewService(&service.Deps{
		Mencache:      mencacheObj,
		Repositories:  repos,
		Logger:        srv.Logger,
		Observability: traceLoggerObservability,
	})

	handlers := handler.NewHandler(&handler.Deps{
		Service: services,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbroles.RegisterRoleQueryServiceServer(gs, handlers.Role)
		pbroles.RegisterRoleCommandServiceServer(gs, handlers.RoleCommand)
		pbuserrole.RegisterUserRoleServiceServer(gs, handlers.UserRole)
	}

	return srv, nil
}
