package repository

import (
	"context"

	roleadapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/role"
	userroleadapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/user_role"
	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type UserQueryRepository interface {
	FindAllUsers(ctx context.Context, req *requests.FindAllUsers) ([]*models.UserRow, error)
	FindById(ctx context.Context, user_id int) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByEmailWithPassword(ctx context.Context, email string) (*models.User, error)
	FindByEmailAndVerify(ctx context.Context, email string) (*models.User, error)
	FindByVerificationCode(ctx context.Context, verification_code string) (*models.User, error)
	FindByActive(ctx context.Context, req *requests.FindAllUsers) ([]*models.UserActiveRow, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*models.UserTrashedRow, error)
}

type UserCommandRepository interface {
	CreateUser(ctx context.Context, request *requests.CreateUserRequest) (*models.User, error)
	CreateUserRecord(ctx context.Context, request *requests.RegisterRequest) (*models.User, error)
	UpdateUser(ctx context.Context, request *requests.UpdateUserRequest) (*models.User, error)
	UpdateUserIsVerified(ctx context.Context, user_id int, is_verified bool) (*models.User, error)
	UpdateUserPassword(ctx context.Context, user_id int, password string) (*models.User, error)
	TrashedUser(ctx context.Context, user_id int) (*models.User, error)
	RestoreUser(ctx context.Context, user_id int) (*models.User, error)
	DeleteUserPermanent(ctx context.Context, user_id int) (bool, error)
	RestoreAllUser(ctx context.Context) (bool, error)
	DeleteAllUserPermanent(ctx context.Context) (bool, error)
}

// RoleRepository is the role read contract, backed by the shared role adapter:
// roles live behind the role service, not the users table.
type RoleRepository = roleadapter.QueryRepository

// UserRoleRepository is the role-assignment contract, backed by the shared
// user-role adapter, so CreateUser can attach the default role.
type UserRoleRepository = userroleadapter.CommandRepository

// Compile-time assertions: the shared adapters satisfy the repository contracts
// directly.
var (
	_ RoleRepository     = (*roleadapter.Repository)(nil)
	_ UserRoleRepository = (*userroleadapter.Repository)(nil)
)
