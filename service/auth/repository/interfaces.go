package repository

import (
	"context"

	roleadapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/role"
	useradapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/user"
	userroleadapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/user_role"
	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"gorm.io/gorm"
)

// UserRepository is satisfied by the shared user gRPC adapter: reads (login,
// verification, password reset) and writes (registration) both go through the
// User service, which owns the users table.
type UserRepository interface {
	useradapter.QueryRepository
	useradapter.CommandRepository
}

type ResetTokenRepository interface {
	FindByToken(ctx context.Context, token string) (*models.ResetToken, error)
	CreateResetToken(ctx context.Context, req *requests.CreateResetTokenRequest) (*models.ResetToken, error)
	// CreateResetTokenInTx persists the reset token inside the given GORM
	// transaction so the caller can commit the token write and its outbox event
	// atomically (Phase 6).
	CreateResetTokenInTx(ctx context.Context, tx *gorm.DB, req *requests.CreateResetTokenRequest) (*models.ResetToken, error)
	DeleteResetToken(ctx context.Context, user_id int) error
}

type RefreshTokenRepository interface {
	FindByToken(ctx context.Context, token string) (*models.RefreshToken, error)
	FindByUserId(ctx context.Context, user_id int) (*models.RefreshToken, error)
	CreateRefreshToken(ctx context.Context, req *requests.CreateRefreshToken) (*models.RefreshToken, error)
	UpdateRefreshToken(ctx context.Context, req *requests.UpdateRefreshToken) (*models.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	DeleteRefreshTokenByUserId(ctx context.Context, user_id int) error
}

// UserRoleRepository is satisfied by the shared user-role gRPC adapter: the
// user_roles writes live behind the UserRole service.
type UserRoleRepository interface {
	userroleadapter.CommandRepository
}

// RoleRepository is satisfied by the shared role gRPC adapter.
type RoleRepository interface {
	roleadapter.QueryRepository
}

type OutboxRepository interface {
	// EnqueueInTx persists a pending outbox event inside the given GORM
	// transaction so the business write and the event commit atomically. The
	// relay (pkg/outbox) then publishes it to Kafka with retry/DLQ semantics.
	EnqueueInTx(ctx context.Context, tx *gorm.DB, topic, key string, payload []byte) error
}
