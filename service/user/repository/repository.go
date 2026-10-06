package repository

import (
	pbroles "github.com/MamangRust/monolith-point-of-sale-pb/roles"
	pbuserrole "github.com/MamangRust/monolith-point-of-sale-pb/user_role"
	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter"
	roleadapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/role"
	userroleadapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/user_role"
	"gorm.io/gorm"
)

// GuardOptions carries the resilience guard options for each outbound gRPC
// dependency.
type GuardOptions struct {
	Role     []adapter.GuardOption
	UserRole []adapter.GuardOption
}

// Repositories bundles the user data-access repositories. It uses named fields
// (not embedding) because UserQueryRepository and RoleRepository both declare
// FindById.
type Repositories struct {
	UserCommand UserCommandRepository
	UserQuery   UserQueryRepository
	Role        RoleRepository
	UserRole    UserRoleRepository
}

// Deps holds the DB handle plus the raw gRPC clients for the role and user-role
// services. The adapters are built here with the supplied guard options so
// resilience (timeout/circuit-breaker/bulkhead) is applied uniformly.
type Deps struct {
	Db       *gorm.DB
	Role     pbroles.RoleQueryServiceClient
	UserRole pbuserrole.UserRoleServiceClient
	Guards   GuardOptions
}

// NewRepositories builds the GORM-backed user repositories and wraps the remote
// role/user-role clients in the shared adapters.
func NewRepositories(deps *Deps) Repositories {
	return Repositories{
		UserCommand: NewUserCommandRepository(deps.Db),
		UserQuery:   NewUserQueryRepository(deps.Db),
		Role:        roleadapter.New(deps.Role, nil, deps.Guards.Role...),
		UserRole:    userroleadapter.New(deps.UserRole, deps.Guards.UserRole...),
	}
}
