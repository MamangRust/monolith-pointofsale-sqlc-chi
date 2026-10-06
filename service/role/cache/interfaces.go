package mencache

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type RoleCommandCache interface {
	DeleteCachedRole(ctx context.Context, id int)
	DeleteCachedRoleAllCache(ctx context.Context)
}

type RoleQueryCache interface {
	SetCachedRoles(ctx context.Context, req *requests.FindAllRoles, data []*models.RoleRow, total *int)
	SetCachedRoleById(ctx context.Context, data *models.Role)
	SetCachedRoleByUserId(ctx context.Context, userId int, data []*models.Role)
	SetCachedRoleActive(ctx context.Context, req *requests.FindAllRoles, data []*models.RoleActiveRow, total *int)
	SetCachedRoleTrashed(ctx context.Context, req *requests.FindAllRoles, data []*models.RoleTrashedRow, total *int)

	GetCachedRoles(ctx context.Context, req *requests.FindAllRoles) ([]*models.RoleRow, *int, bool)
	GetCachedRoleByUserId(ctx context.Context, userId int) ([]*models.Role, bool)
	GetCachedRoleById(ctx context.Context, id int) (*models.Role, bool)
	GetCachedRoleActive(ctx context.Context, req *requests.FindAllRoles) ([]*models.RoleActiveRow, *int, bool)
	GetCachedRoleTrashed(ctx context.Context, req *requests.FindAllRoles) ([]*models.RoleTrashedRow, *int, bool)
}
