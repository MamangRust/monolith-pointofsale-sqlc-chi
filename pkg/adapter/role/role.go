// Package role adapts the Role service gRPC API into the shared domain model.
package role

import (
	"context"

	pbrole "github.com/MamangRust/monolith-point-of-sale-pb/roles"
	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter"
	"github.com/MamangRust/monolith-point-of-sale-pkg/adapter/convert"
	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-pkg/resilience"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/role_errors"
	userrole_errors "github.com/MamangRust/monolith-point-of-sale-shared/errors/user_role_errors"
)

// QueryRepository is the contract consumers depend on for role reads.
type QueryRepository interface {
	FindById(ctx context.Context, roleID int) (*models.Role, error)
	FindByName(ctx context.Context, name string) (*models.Role, error)
}

// CommandRepository is the contract consumers depend on for role assignment.
type CommandRepository interface {
	AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*models.UserRole, error)
	RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error
}

// Repository implements QueryRepository and CommandRepository on top of the
// generated role query and command service clients.
type Repository struct {
	query   pbrole.RoleQueryServiceClient
	command pbrole.RoleCommandServiceClient
	guard   *resilience.DependencyGuard
}

// New builds a role adapter. Either client may be nil when the consumer only
// needs the other half, but calling a method on a nil client panics. Passing
// zero options leaves the guard nil, which makes DependencyGuard.Call a plain
// passthrough.
func New(query pbrole.RoleQueryServiceClient, command pbrole.RoleCommandServiceClient, opts ...adapter.GuardOption) *Repository {
	r := &Repository{query: query, command: command}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r *Repository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

func (r *Repository) FindById(ctx context.Context, roleID int) (*models.Role, error) {
	var resp *pbrole.ApiResponseRole
	err := r.guard.Call(ctx, func(ctx context.Context) error {
		var callErr error
		resp, callErr = r.query.FindByIdRole(ctx, &pbrole.FindByIdRoleRequest{RoleId: int32(roleID)})
		return callErr
	})
	if err != nil || resp == nil || resp.Data == nil {
		return nil, role_errors.ErrRoleNotFound.WithInternal(err)
	}
	return roleToModel(resp.Data), nil
}

func (r *Repository) FindByName(ctx context.Context, name string) (*models.Role, error) {
	var resp *pbrole.ApiResponseRole
	err := r.guard.Call(ctx, func(ctx context.Context) error {
		var callErr error
		resp, callErr = r.query.FindByNameRole(ctx, &pbrole.FindByNameRoleRequest{Name: name})
		return callErr
	})
	if err != nil || resp == nil || resp.Data == nil {
		return nil, role_errors.ErrRoleNotFound.WithInternal(err)
	}
	return roleToModel(resp.Data), nil
}

func (r *Repository) AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*models.UserRole, error) {
	var resp *pbrole.ApiResponseUserRole
	err := r.guard.Call(ctx, func(ctx context.Context) error {
		var callErr error
		resp, callErr = r.command.AssignRoleToUser(ctx, &pbrole.AssignRoleToUserRequest{
			UserId: int32(req.UserId),
			RoleId: int32(req.RoleId),
		})
		return callErr
	})
	if err != nil || resp == nil || resp.Data == nil {
		return nil, userrole_errors.ErrAssignRoleToUser.WithInternal(err)
	}
	return &models.UserRole{
		UserRoleID: resp.Data.UserRoleId,
		UserID:     resp.Data.UserId,
		RoleID:     resp.Data.RoleId,
	}, nil
}

func (r *Repository) RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error {
	err := r.guard.Call(ctx, func(ctx context.Context) error {
		_, callErr := r.command.RemoveRoleFromUser(ctx, &pbrole.RemoveRoleFromUserRequest{
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

func roleToModel(role *pbrole.RoleResponse) *models.Role {
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
