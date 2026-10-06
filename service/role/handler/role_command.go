package handler

import (
	"context"
	pbroles "github.com/MamangRust/monolith-point-of-sale-pb/roles"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-role/service"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/role_errors"
)

type roleCommandHandleGrpc struct {
	pbroles.UnimplementedRoleCommandServiceServer
	roleCommand service.RoleCommandService
	logger      logger.LoggerInterface
}

func NewRoleCommandHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pbroles.RoleCommandServiceServer {
	return &roleCommandHandleGrpc{
		roleCommand: service.RoleCommand,
		logger:      logger,
	}
}

func (s *roleCommandHandleGrpc) AssignRoleToUser(ctx context.Context, reqPb *pbroles.AssignRoleToUserRequest) (*pbroles.ApiResponseUserRole, error) {
	s.logger.Info("AssignRoleToUser called", zap.Int32("user_id", reqPb.GetUserId()), zap.Int32("role_id", reqPb.GetRoleId()))

	if reqPb.GetUserId() <= 0 || reqPb.GetRoleId() <= 0 {
		return nil, errors.ErrBadRequest.WithMessage("Invalid user or role ID")
	}

	req := &requests.CreateUserRoleRequest{
		UserId: int(reqPb.GetUserId()),
		RoleId: int(reqPb.GetRoleId()),
	}

	userRole, err := s.roleCommand.AssignRoleToUser(ctx, req)
	if err != nil {
		s.logger.Error("AssignRoleToUser failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("AssignRoleToUser success")

	return &pbroles.ApiResponseUserRole{
		Status:  "success",
		Message: "Successfully assigned role to user",
		Data:    mapUserRoleToProto(userRole),
	}, nil
}

func (s *roleCommandHandleGrpc) RemoveRoleFromUser(ctx context.Context, reqPb *pbroles.RemoveRoleFromUserRequest) (*pbroles.ApiResponseUserRole, error) {
	s.logger.Info("RemoveRoleFromUser called", zap.Int32("user_id", reqPb.GetUserId()), zap.Int32("role_id", reqPb.GetRoleId()))

	if reqPb.GetUserId() <= 0 || reqPb.GetRoleId() <= 0 {
		return nil, errors.ErrBadRequest.WithMessage("Invalid user or role ID")
	}

	req := &requests.RemoveUserRoleRequest{
		UserId: int(reqPb.GetUserId()),
		RoleId: int(reqPb.GetRoleId()),
	}

	if err := s.roleCommand.RemoveRoleFromUser(ctx, req); err != nil {
		s.logger.Error("RemoveRoleFromUser failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RemoveRoleFromUser success")

	return &pbroles.ApiResponseUserRole{
		Status:  "success",
		Message: "Successfully removed role from user",
	}, nil
}

func (s *roleCommandHandleGrpc) CreateRole(ctx context.Context, reqPb *pbroles.CreateRoleRequest) (*pbroles.ApiResponseRole, error) {
	s.logger.Info("CreateRole called", zap.String("name", reqPb.Name))

	req := &requests.CreateRoleRequest{
		Name: reqPb.Name,
	}

	if err := req.Validate(); err != nil {
		return nil, role_errors.ErrGrpcValidateCreateRole
	}

	role, err := s.roleCommand.CreateRole(ctx, req)
	if err != nil {
		s.logger.Error("CreateRole failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("CreateRole success")

	return &pbroles.ApiResponseRole{
		Status:  "success",
		Message: "Successfully created role",
		Data:    mapResponseRole(role),
	}, nil
}

func (s *roleCommandHandleGrpc) UpdateRole(ctx context.Context, reqPb *pbroles.UpdateRoleRequest) (*pbroles.ApiResponseRole, error) {
	s.logger.Info("UpdateRole called", zap.Int32("id", reqPb.GetId()))

	roleID := int(reqPb.GetId())
	if roleID <= 0 {
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	name := reqPb.GetName()
	req := &requests.UpdateRoleRequest{
		ID:   &roleID,
		Name: name,
	}

	if err := req.Validate(); err != nil {
		return nil, role_errors.ErrGrpcValidateUpdateRole
	}

	role, err := s.roleCommand.UpdateRole(ctx, req)
	if err != nil {
		s.logger.Error("UpdateRole failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("UpdateRole success")

	return &pbroles.ApiResponseRole{
		Status:  "success",
		Message: "Successfully updated role",
		Data:    mapResponseRole(role),
	}, nil
}

func (s *roleCommandHandleGrpc) TrashedRole(ctx context.Context, req *pbroles.FindByIdRoleRequest) (*pbroles.ApiResponseRole, error) {
	s.logger.Info("TrashedRole called", zap.Int32("id", req.GetRoleId()))

	roleID := int(req.GetRoleId())
	if roleID <= 0 {
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	role, err := s.roleCommand.TrashedRole(ctx, roleID)
	if err != nil {
		s.logger.Error("TrashedRole failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("TrashedRole success")

	return &pbroles.ApiResponseRole{
		Status:  "success",
		Message: "Successfully trashed role",
		Data:    mapResponseRole(role),
	}, nil
}

func (s *roleCommandHandleGrpc) RestoreRole(ctx context.Context, req *pbroles.FindByIdRoleRequest) (*pbroles.ApiResponseRole, error) {
	s.logger.Info("RestoreRole called", zap.Int32("id", req.GetRoleId()))

	roleID := int(req.GetRoleId())
	if roleID <= 0 {
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	role, err := s.roleCommand.RestoreRole(ctx, roleID)
	if err != nil {
		s.logger.Error("RestoreRole failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreRole success")

	return &pbroles.ApiResponseRole{
		Status:  "success",
		Message: "Successfully restored role",
		Data:    mapResponseRole(role),
	}, nil
}

func (s *roleCommandHandleGrpc) DeleteRolePermanent(ctx context.Context, req *pbroles.FindByIdRoleRequest) (*pbroles.ApiResponseRoleDelete, error) {
	s.logger.Info("DeleteRolePermanent called", zap.Int32("id", req.GetRoleId()))

	id := int(req.GetRoleId())
	if id <= 0 {
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	_, err := s.roleCommand.DeleteRolePermanent(ctx, id)
	if err != nil {
		s.logger.Error("DeleteRolePermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteRolePermanent success")

	return &pbroles.ApiResponseRoleDelete{
		Status:  "success",
		Message: "Successfully deleted role permanently",
	}, nil
}

func (s *roleCommandHandleGrpc) RestoreAllRole(ctx context.Context, req *emptypb.Empty) (*pbroles.ApiResponseRoleAll, error) {
	s.logger.Info("RestoreAllRole called")

	_, err := s.roleCommand.RestoreAllRole(ctx)
	if err != nil {
		s.logger.Error("RestoreAllRole failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreAllRole success")

	return &pbroles.ApiResponseRoleAll{
		Status:  "success",
		Message: "Successfully restored all roles",
	}, nil
}

func (s *roleCommandHandleGrpc) DeleteAllRolePermanent(ctx context.Context, req *emptypb.Empty) (*pbroles.ApiResponseRoleAll, error) {
	s.logger.Info("DeleteAllRolePermanent called")

	_, err := s.roleCommand.DeleteAllRolePermanent(ctx)
	if err != nil {
		s.logger.Error("DeleteAllRolePermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteAllRolePermanent success")

	return &pbroles.ApiResponseRoleAll{
		Status:  "success",
		Message: "Successfully deleted all roles",
	}, nil
}

// Map helpers
