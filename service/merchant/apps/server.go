package apps

import (
	"context"
	"os"
	"time"

	pbmerchants "github.com/MamangRust/monolith-point-of-sale-pb/merchants"
	pbmerchants_document "github.com/MamangRust/monolith-point-of-sale-pb/merchants_document"
	pbusers "github.com/MamangRust/monolith-point-of-sale-pb/users"

	mencache "github.com/MamangRust/monolith-point-of-sale-merchant/cache"
	"github.com/MamangRust/monolith-point-of-sale-merchant/handler"
	"github.com/MamangRust/monolith-point-of-sale-merchant/repository"
	"github.com/MamangRust/monolith-point-of-sale-merchant/service"
	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter"
	"github.com/MamangRust/monolith-point-of-sale-pkg/kafka"
	"github.com/MamangRust/monolith-point-of-sale-pkg/outbox"
	"github.com/MamangRust/monolith-point-of-sale-pkg/resilience"
	"github.com/MamangRust/monolith-point-of-sale-pkg/server"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewServer(cfg *server.Config) (*server.GRPCServer, error) {
	srv, err := server.New(cfg)
	if err != nil {
		return nil, err
	}

	userAddr := os.Getenv("GRPC_USER_ADDR")
	if userAddr == "" {
		userAddr = "localhost:50053"
	}

	srv.Logger.Info("Connecting to User service via gRPC", zap.String("addr", userAddr))
	userConn, err := server.NewGRPCClient(userAddr)
	if err != nil {
		return nil, err
	}

	go func() {
		<-srv.Ctx.Done()
		srv.Logger.Info("Closing merchant service remote gRPC connections")
		userConn.Close()
	}()

	userClient := pbusers.NewUserQueryServiceClient(userConn)

	guardUser := resilience.NewDependencyGuard("user", 5, 30, 100, 3*time.Second, srv.Logger)

	repos := repository.NewRepositories(srv.GormDB, userClient,
		repository.GuardOptions{User: []adapter.GuardOption{adapter.WithDependencyGuard(guardUser)}},
	)

	var myKafka *kafka.Kafka
	if brokers := os.Getenv("KAFKA_BROKERS"); brokers != "" {
		myKafka = kafka.NewKafka(srv.Logger, []string{brokers})
	}
	traceLoggerObservability := observability.NewTraceLoggerObservability(srv.Logger)

	mencacheObj := mencache.NewMencache(srv.CacheStore)

	outboxService := outbox.NewOutboxService(srv.GormDB, myKafka, srv.Logger)

	services := service.NewService(&service.Deps{
		Ctx:           context.Background(),
		Mencache:      mencacheObj,
		Kafka:         myKafka,
		Repositories:  repos,
		DB:            srv.GormDB,
		Outbox:        outboxService,
		Logger:        srv.Logger,
		Observability: traceLoggerObservability,
	})

	handlers := handler.NewHandler(&handler.Deps{
		Service: services,
		Logger:  srv.Logger,
	})

	srv.RegisterServices = func(gs *grpc.Server) {
		pbmerchants.RegisterMerchantQueryServiceServer(gs, handlers.Merchant)
		pbmerchants_document.RegisterMerchantDocumentQueryServiceServer(gs, handlers.MerchantDocument)
	}

	// Start the outbox relay so events committed with the business writes are
	// published to Kafka with durable retry and dead-letter semantics.
	go outboxService.Start(srv.Ctx, outbox.OutboxRelayInterval, outbox.OutboxRelayBatchSize)

	return srv, nil
}
