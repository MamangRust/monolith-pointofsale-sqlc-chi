package handler

import (
	"context"

	pbroles "github.com/MamangRust/monolith-point-of-sale-pb/roles"
	pbuserrole "github.com/MamangRust/monolith-point-of-sale-pb/user_role"
	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-role/service"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

// userRoleHandleGrpc serves the UserRoleService, the split-out surface that owns
// user-role assignment and per-user role resolution.
type userRoleHandleGrpc struct {
	pbuserrole.UnimplementedUserRoleServiceServer
	roleQuery   service.RoleQueryService
	roleCommand service.RoleCommandService
	logger      logger.LoggerInterface
}

func NewUserRoleHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pbuserrole.UserRoleServiceServer {
	return &userRoleHandleGrpc{
		roleQuery:   service.RoleQuery,
		roleCommand: service.RoleCommand,
		logger:      logger,
	}
}

func (s *userRoleHandleGrpc) FindByUserId(ctx context.Context, req *pbuserrole.FindByIdUserRoleRequest) (*pbroles.ApiResponsesRole, error) {
	s.logger.Info("FindByUserId (user_role) called", zap.Int32("user_id", req.GetUserId()))

	if req.GetUserId() <= 0 {
		return nil, errors.ErrBadRequest.WithMessage("Invalid user ID")
	}

	roles, err := s.roleQuery.FindByUserId(ctx, int(req.GetUserId()))
	if err != nil {
		s.logger.Error("FindByUserId (user_role) failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindByUserId (user_role) success")

	return &pbroles.ApiResponsesRole{
		Status:  "success",
		Message: "Successfully fetched user roles",
		Data:    mapResponsesRoleFromDB(roles),
	}, nil
}

func (s *userRoleHandleGrpc) AssignRoleToUser(ctx context.Context, reqPb *pbuserrole.AssignRoleToUserRequest) (*pbuserrole.ApiResponseUserRole, error) {
	s.logger.Info("AssignRoleToUser (user_role) called", zap.Int32("user_id", reqPb.GetUserId()), zap.Int32("role_id", reqPb.GetRoleId()))

	if reqPb.GetUserId() <= 0 || reqPb.GetRoleId() <= 0 {
		return nil, errors.ErrBadRequest.WithMessage("Invalid user or role ID")
	}

	userRole, err := s.roleCommand.AssignRoleToUser(ctx, &requests.CreateUserRoleRequest{
		UserId: int(reqPb.GetUserId()),
		RoleId: int(reqPb.GetRoleId()),
	})
	if err != nil {
		s.logger.Error("AssignRoleToUser (user_role) failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("AssignRoleToUser (user_role) success")

	return &pbuserrole.ApiResponseUserRole{
		Status:  "success",
		Message: "Successfully assigned role to user",
		Data:    mapUserRoleResponse(userRole),
	}, nil
}

func (s *userRoleHandleGrpc) RemoveRoleFromUser(ctx context.Context, reqPb *pbuserrole.RemoveRoleFromUserRequest) (*emptypb.Empty, error) {
	s.logger.Info("RemoveRoleFromUser (user_role) called", zap.Int32("user_id", reqPb.GetUserId()), zap.Int32("role_id", reqPb.GetRoleId()))

	if reqPb.GetUserId() <= 0 || reqPb.GetRoleId() <= 0 {
		return nil, errors.ErrBadRequest.WithMessage("Invalid user or role ID")
	}

	if err := s.roleCommand.RemoveRoleFromUser(ctx, &requests.RemoveUserRoleRequest{
		UserId: int(reqPb.GetUserId()),
		RoleId: int(reqPb.GetRoleId()),
	}); err != nil {
		s.logger.Error("RemoveRoleFromUser (user_role) failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RemoveRoleFromUser (user_role) success")

	return &emptypb.Empty{}, nil
}

func mapUserRoleResponse(userRole *models.UserRole) *pbuserrole.UserRoleResponse {
	if userRole == nil {
		return nil
	}
	return &pbuserrole.UserRoleResponse{
		UserRoleId: userRole.UserRoleID,
		UserId:     userRole.UserID,
		RoleId:     userRole.RoleID,
	}
}
