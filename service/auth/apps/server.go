package apps

import (
	"os"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-pb"
	pbroles "github.com/MamangRust/monolith-point-of-sale-pb/roles"
	pbuserrole "github.com/MamangRust/monolith-point-of-sale-pb/user_role"
	pbusers "github.com/MamangRust/monolith-point-of-sale-pb/users"
	"google.golang.org/grpc"

	mencache "github.com/MamangRust/monolith-point-of-sale-auth/cache"
	"github.com/MamangRust/monolith-point-of-sale-auth/handler"
	"github.com/MamangRust/monolith-point-of-sale-auth/repository"
	"github.com/MamangRust/monolith-point-of-sale-auth/service"
	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter"
	"github.com/MamangRust/monolith-point-of-sale-pkg/auth"
	"github.com/MamangRust/monolith-point-of-sale-pkg/hash"
	"github.com/MamangRust/monolith-point-of-sale-pkg/kafka"
	"github.com/MamangRust/monolith-point-of-sale-pkg/outbox"
	"github.com/MamangRust/monolith-point-of-sale-pkg/resilience"
	"github.com/MamangRust/monolith-point-of-sale-pkg/server"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	"github.com/spf13/viper"
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

	tokenManager, err := auth.NewManager(viper.GetString("SECRET_KEY"))
	if err != nil {
		return nil, err
	}

	userAddr := getEnv("GRPC_USER_ADDR", "localhost:50053")
	roleAddr := getEnv("GRPC_ROLE_ADDR", "localhost:50052")

	userConn, err := server.NewGRPCClient(userAddr)
	if err != nil {
		return nil, err
	}

	roleConn, err := server.NewGRPCClient(roleAddr)
	if err != nil {
		userConn.Close()
		return nil, err
	}

	go func() {
		<-srv.Ctx.Done()
		userConn.Close()
		roleConn.Close()
	}()

	userQueryClient := pbusers.NewUserQueryServiceClient(userConn)
	userCommandClient := pbusers.NewUserCommandServiceClient(userConn)
	roleQueryClient := pbroles.NewRoleQueryServiceClient(roleConn)
	roleCommandClient := pbroles.NewRoleCommandServiceClient(roleConn)
	userRoleClient := pbuserrole.NewUserRoleServiceClient(roleConn)

	guardUser := resilience.NewDependencyGuard("user", 5, 30, 100, 3*time.Second, srv.Logger)
	guardRole := resilience.NewDependencyGuard("role", 5, 30, 100, 3*time.Second, srv.Logger)
	guardUserRole := resilience.NewDependencyGuard("user_role", 5, 30, 100, 3*time.Second, srv.Logger)

	repos := repository.NewRepositories(
		srv.GormDB,
		userQueryClient,
		userCommandClient,
		roleQueryClient,
		roleCommandClient,
		userRoleClient,
		repository.GuardOptions{
			User:     []adapter.GuardOption{adapter.WithDependencyGuard(guardUser)},
			Role:     []adapter.GuardOption{adapter.WithDependencyGuard(guardRole)},
			UserRole: []adapter.GuardOption{adapter.WithDependencyGuard(guardUserRole)},
		},
	)

	hash := hash.NewHashingPassword()

	var myKafka *kafka.Kafka
	if brokers := os.Getenv("KAFKA_BROKERS"); brokers != "" {
		myKafka = kafka.NewKafka(srv.Logger, []string{brokers})
	}

	mencacheObj := mencache.NewMencache(srv.CacheStore)

	outboxService := outbox.NewOutboxService(srv.GormDB, myKafka, srv.Logger)

	services := service.NewService(&service.Deps{
		Mencache:      mencacheObj,
		Repositories:  repos,
		Hash:          hash,
		Token:         tokenManager,
		Logger:        srv.Logger,
		Kafka:         myKafka,
		GormDB:        srv.GormDB,
		Observability: observability.NewTraceLoggerObservability(srv.Logger),
	})

	handlers := handler.NewHandler(&handler.Deps{
		Service: services,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pb.RegisterAuthServiceServer(gs, handlers.Auth)
	}

	go outboxService.Start(srv.Ctx, outbox.OutboxRelayInterval, outbox.OutboxRelayBatchSize)

	return srv, nil
}
