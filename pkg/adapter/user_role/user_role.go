// Package user_role adapts the UserRole service gRPC API into the shared domain
// model. It owns user-role assignment and per-user role resolution; role CRUD
// and role lookups live in the role adapter.
package user_role

import (
	"context"

	pbroles "github.com/MamangRust/monolith-point-of-sale-pb/roles"
	pbuserrole "github.com/MamangRust/monolith-point-of-sale-pb/user_role"
	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter"
	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter/convert"
	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-pkg/resilience"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	userrole_errors "github.com/MamangRust/monolith-point-of-sale-shared/errors/user_role_errors"
)

// QueryRepository is the read path consumers use to resolve a user's roles.
type QueryRepository interface {
	FindByUserId(ctx context.Context, userID int) ([]*models.Role, error)
}

// CommandRepository is the write path consumers use to (un)assign roles.
type CommandRepository interface {
	AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*models.UserRole, error)
	RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error
}

// Repository implements QueryRepository and CommandRepository on top of the
// generated user-role service client.
type Repository struct {
	client pbuserrole.UserRoleServiceClient
	guard  *resilience.DependencyGuard
}

// SetGuard implements adapter.GuardSetter.
func (r *Repository) SetGuard(g *resilience.DependencyGuard) { r.guard = g }

// New wraps a generated user-role service client into a Repository. Resilience
// (timeout/circuit-breaker/bulkhead) is applied per call through the optional
// dependency guard; a nil guard is a passthrough.
func New(client pbuserrole.UserRoleServiceClient, opts ...adapter.GuardOption) *Repository {
	r := &Repository{client: client}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// FindByUserId implements QueryRepository. The user-role service returns the
// roles resolved for the user, which we map straight into the domain model.
func (r *Repository) FindByUserId(ctx context.Context, userID int) ([]*models.Role, error) {
	var res *pbroles.ApiResponsesRole
	err := r.guard.Call(ctx, func(callCtx context.Context) error {
		var callErr error
		res, callErr = r.client.FindByUserId(callCtx, &pbuserrole.FindByIdUserRoleRequest{UserId: int32(userID)})
		return callErr
	})
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}

	roles := make([]*models.Role, 0, len(res.Data))
	for _, item := range res.Data {
		if item == nil {
			continue
		}
		roles = append(roles, roleToModel(item))
	}
	return roles, nil
}

// AssignRoleToUser implements CommandRepository.
func (r *Repository) AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*models.UserRole, error) {
	var res *pbuserrole.ApiResponseUserRole
	err := r.guard.Call(ctx, func(callCtx context.Context) error {
		var callErr error
		res, callErr = r.client.AssignRoleToUser(callCtx, &pbuserrole.AssignRoleToUserRequest{
			UserId: int32(req.UserId),
			RoleId: int32(req.RoleId),
		})
		return callErr
	})
	if err != nil || res == nil || res.Data == nil {
		return nil, userrole_errors.ErrAssignRoleToUser.WithInternal(err)
	}
	return userRoleToModel(res.Data), nil
}

// RemoveRoleFromUser implements CommandRepository.
func (r *Repository) RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error {
	err := r.guard.Call(ctx, func(callCtx context.Context) error {
		_, callErr := r.client.RemoveRoleFromUser(callCtx, &pbuserrole.RemoveRoleFromUserRequest{
			UserId: int32(req.UserId),
			RoleId: int32(req.RoleId),
		})
		return callErr
	})
	if err != nil {
		return userrole_errors.ErrRemoveRole.WithInternal(err)
	}
	return nil
}

func roleToModel(role *pbroles.RoleResponse) *models.Role {
	if role == nil {
		return nil
	}
	return &models.Role{
		RoleID:    role.Id,
		RoleName:  role.Name,
		CreatedAt: convert.Time(role.CreatedAt),
		UpdatedAt: convert.Time(role.UpdatedAt),
	}
}

func userRoleToModel(ur *pbuserrole.UserRoleResponse) *models.UserRole {
	if ur == nil {
		return nil
	}
	return &models.UserRole{
		UserRoleID: ur.UserRoleId,
		UserID:     ur.UserId,
		RoleID:     ur.RoleId,
	}
}
