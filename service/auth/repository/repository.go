package repository

import (
	pbroles "github.com/MamangRust/monolith-point-of-sale-pb/roles"
	pbuserrole "github.com/MamangRust/monolith-point-of-sale-pb/user_role"
	pbusers "github.com/MamangRust/monolith-point-of-sale-pb/users"
	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter"
	roleadapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/role"
	useradapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/user"
	userroleadapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/user_role"
	"gorm.io/gorm"
)

// Repositories bundles all auth repositories. It uses named fields (not
// embedding) because UserRepository and RoleRepository both declare FindById,
// and RefreshTokenRepository and ResetTokenRepository both declare FindByToken.
type Repositories struct {
	User         UserRepository
	RefreshToken RefreshTokenRepository
	UserRole     UserRoleRepository
	Role         RoleRepository
	ResetToken   ResetTokenRepository
	Outbox       OutboxRepository
}

// GuardOptions collects the dependency guards for each remote dependency.
type GuardOptions struct {
	User     []adapter.GuardOption
	Role     []adapter.GuardOption
	UserRole []adapter.GuardOption
}

// Compile-time assertions: the shared adapters satisfy auth's repository
// contracts directly.
var (
	_ UserRepository     = (*useradapter.Repository)(nil)
	_ RoleRepository     = (*roleadapter.Repository)(nil)
	_ UserRoleRepository = (*userroleadapter.Repository)(nil)
)

func NewRepositories(
	db *gorm.DB,
	userQueryClient pbusers.UserQueryServiceClient,
	userCommandClient pbusers.UserCommandServiceClient,
	roleQueryClient pbroles.RoleQueryServiceClient,
	roleCommandClient pbroles.RoleCommandServiceClient,
	userRoleClient pbuserrole.UserRoleServiceClient,
	guards ...GuardOptions,
) Repositories {
	var g GuardOptions
	if len(guards) > 0 {
		g = guards[0]
	}

	userRepo := useradapter.New(userQueryClient, userCommandClient, g.User...)
	roleRepo := roleadapter.New(roleQueryClient, roleCommandClient, g.Role...)
	userRoleRepo := userroleadapter.New(userRoleClient)

	return Repositories{
		User:         userRepo,
		RefreshToken: NewRefreshTokenRepository(db),
		UserRole:     userRoleRepo,
		Role:         roleRepo,
		ResetToken:   NewResetTokenRepository(db),
		Outbox:       NewOutboxRepository(),
	}
}
