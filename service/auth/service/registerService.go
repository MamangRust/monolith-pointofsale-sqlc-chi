package service

import (
	"context"
	"strconv"
	"strings"
	"time"

	mencache "github.com/MamangRust/monolith-point-of-sale-auth/cache"
	"github.com/MamangRust/monolith-point-of-sale-auth/repository"
	"github.com/MamangRust/monolith-point-of-sale-pkg/database/gormtx"
	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-pkg/email"
	"github.com/MamangRust/monolith-point-of-sale-pkg/event"
	"github.com/MamangRust/monolith-point-of-sale-pkg/hash"
	"github.com/MamangRust/monolith-point-of-sale-pkg/kafka"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-pkg/randomstring"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	sharederrorhandler "github.com/MamangRust/monolith-point-of-sale-shared/errorhandler"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/user_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RegisterServiceDeps struct {
	Cache mencache.RegisterCache

	DB *gorm.DB

	User repository.UserRepository

	Role repository.RoleRepository

	UserRole repository.UserRoleRepository

	Outbox repository.OutboxRepository

	Hash hash.HashPassword

	Kafka *kafka.Kafka

	Logger logger.LoggerInterface

	Observability observability.TraceLoggerObservability
}

type registerService struct {
	mencache mencache.RegisterCache

	db *gorm.DB

	user repository.UserRepository

	role repository.RoleRepository

	userRole repository.UserRoleRepository

	outbox repository.OutboxRepository

	hash hash.HashPassword

	kafka *kafka.Kafka

	logger logger.LoggerInterface

	observability observability.TraceLoggerObservability
}

func NewRegisterService(params *RegisterServiceDeps) *registerService {
	return &registerService{
		mencache:      params.Cache,
		db:            params.DB,
		user:          params.User,
		role:          params.Role,
		userRole:      params.UserRole,
		outbox:        params.Outbox,
		hash:          params.Hash,
		kafka:         params.Kafka,
		logger:        params.Logger,
		observability: params.Observability,
	}
}

func (s *registerService) Register(ctx context.Context, request *requests.RegisterRequest) (*models.User, error) {
	const method = "Register"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("email", request.Email))

	defer func() {
		end(status)
	}()

	existingUser, err := s.user.FindByEmail(ctx, request.Email)
	if err == nil && existingUser != nil {
		status = "error"
		return sharederrorhandler.HandleError[*models.User](
			s.logger,
			user_errors.ErrUserEmailAlready,
			method,
			span,
			zap.String("email", request.Email),
		)
	}

	passwordHash, err := s.hash.HashPassword(request.Password)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*models.User](s.logger, err, method, span)
	}
	request.Password = passwordHash

	const defaultRoleName = "ROLE_ADMIN"
	role, err := s.role.FindByName(ctx, defaultRoleName)
	if err != nil || role == nil {
		status = "error"
		return sharederrorhandler.HandleError[*models.User](s.logger, err, method, span, zap.String("role_name", defaultRoleName))
	}

	random, err := randomstring.GenerateRandomString(10)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*models.User](s.logger, err, method, span)
	}
	request.VerifiedCode = random
	// Auto-verify in non-production environments so the full register → verify →
	// login journey is testable locally; production keeps email verification.
	switch strings.ToLower(viper.GetString("APP_ENV")) {
	case "production", "kubernetes":
		request.IsVerified = false
	default:
		request.IsVerified = true
	}

	htmlBody := email.GenerateEmailHTML(map[string]string{
		"Title":   "Welcome to SanEdge",
		"Message": "Your account has been successfully created.",
		"Button":  "Verify Now",
		"Link":    "https://sanedge.example.com/login?verify_code=" + request.VerifiedCode,
	})

	payloadBytes, err := event.MarshalEmail("auth.register", request.Email, "Welcome to SanEdge", htmlBody)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*models.User](s.logger, err, method, span)
	}

	// The user and role writes now live in the User/Role services, so they can
	// no longer share a transaction with the outbox event. The event is still
	// enqueued durably (its own transaction) so the registration email is not
	// lost between the remote write and the publish.
	newUser, err := s.user.CreateUser(ctx, request)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*models.User](s.logger, err, method, span)
	}

	if _, err := s.userRole.AssignRoleToUser(ctx, &requests.CreateUserRoleRequest{
		UserId: int(newUser.UserID),
		RoleId: int(role.RoleID),
	}); err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*models.User](s.logger, err, method, span, zap.Int("user.id", int(newUser.UserID)))
	}

	if s.db != nil {
		if outboxErr := gormtx.RunInTx(ctx, s.db, func(tx *gorm.DB) error {
			return s.outbox.EnqueueInTx(ctx, tx, "email-service-topic-auth-register", strconv.Itoa(int(newUser.UserID)), payloadBytes)
		}); outboxErr != nil {
			s.logger.Error("failed to enqueue registration email", zap.Error(outboxErr), zap.String("email", request.Email))
		}
	} else if s.kafka != nil {
		if sendErr := s.kafka.SendMessage(ctx, "email-service-topic-auth-register", strconv.Itoa(int(newUser.UserID)), payloadBytes); sendErr != nil {
			s.logger.Error("failed to send registration email via kafka", zap.Error(sendErr), zap.String("email", request.Email))
		}
	}

	s.mencache.SetVerificationCodeCache(ctx, request.Email, random, 15*time.Minute)

	logSuccess("User registered successfully",
		zap.String("email", request.Email),
		zap.String("first_name", request.FirstName),
		zap.String("last_name", request.LastName),
	)

	return newUser, nil
}
