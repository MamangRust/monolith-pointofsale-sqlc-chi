package repository

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type RoleQueryRepository interface {
	FindAllRoles(ctx context.Context, req *requests.FindAllRoles) ([]*models.RoleRow, error)
	FindById(ctx context.Context, roleID int) (*models.Role, error)
	FindByName(ctx context.Context, name string) (*models.Role, error)
	FindByUserId(ctx context.Context, userID int) ([]*models.Role, error)
	FindByActiveRole(ctx context.Context, req *requests.FindAllRoles) ([]*models.RoleActiveRow, error)
	FindByTrashedRole(ctx context.Context, req *requests.FindAllRoles) ([]*models.RoleTrashedRow, error)
}

type UserRoleCommandRepository interface {
	AssignRoleToUser(ctx context.Context, request *requests.CreateUserRoleRequest) (*models.UserRole, error)
	RemoveRoleFromUser(ctx context.Context, request *requests.RemoveUserRoleRequest) error
}

type RoleCommandRepository interface {
	CreateRole(ctx context.Context, request *requests.CreateRoleRequest) (*models.Role, error)
	UpdateRole(ctx context.Context, request *requests.UpdateRoleRequest) (*models.Role, error)
	TrashedRole(ctx context.Context, roleID int) (*models.Role, error)
	RestoreRole(ctx context.Context, roleID int) (*models.Role, error)
	DeleteRolePermanent(ctx context.Context, roleID int) (bool, error)
	RestoreAllRole(ctx context.Context) (bool, error)
	DeleteAllRolePermanent(ctx context.Context) (bool, error)
}
