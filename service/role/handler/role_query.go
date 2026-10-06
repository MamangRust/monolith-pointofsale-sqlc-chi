package handler

import (
	"context"
	pbcommon "github.com/MamangRust/monolith-point-of-sale-pb/common"
	pbroles "github.com/MamangRust/monolith-point-of-sale-pb/roles"
	"go.uber.org/zap"
	"math"

	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-role/service"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/role_errors"
)

type roleQueryHandleGrpc struct {
	pbroles.UnimplementedRoleQueryServiceServer
	roleQuery service.RoleQueryService
	logger    logger.LoggerInterface
}

func NewRoleQueryHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pbroles.RoleQueryServiceServer {
	return &roleQueryHandleGrpc{
		roleQuery: service.RoleQuery,
		logger:    logger,
	}
}

func (s *roleQueryHandleGrpc) FindByNameRole(ctx context.Context, req *pbroles.FindByNameRoleRequest) (*pbroles.ApiResponseRole, error) {
	s.logger.Info("FindByNameRole called", zap.String("name", req.GetName()))

	if req.GetName() == "" {
		return nil, errors.ErrBadRequest.WithMessage("Invalid role name")
	}

	role, err := s.roleQuery.FindByName(ctx, req.GetName())
	if err != nil {
		s.logger.Error("FindByNameRole failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindByNameRole success")

	return &pbroles.ApiResponseRole{
		Status:  "success",
		Message: "Successfully fetched role",
		Data:    mapResponseRole(role),
	}, nil
}

func (s *roleQueryHandleGrpc) FindAllRole(ctx context.Context, req *pbroles.FindAllRoleRequest) (*pbroles.ApiResponsePaginationRole, error) {
	s.logger.Info("FindAllRole called", zap.Int32("page", req.GetPage()))

	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllRoles{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	roles, totalRecords, err := s.roleQuery.FindAll(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindAllRoles failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindAllRoles success")

	return &pbroles.ApiResponsePaginationRole{
		Status:     "success",
		Message:    "Successfully fetched role records",
		Data:       mapResponsesRole(roles),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *roleQueryHandleGrpc) FindByIdRole(ctx context.Context, req *pbroles.FindByIdRoleRequest) (*pbroles.ApiResponseRole, error) {
	s.logger.Info("FindByIdRole called", zap.Int32("roleId", req.GetRoleId()))

	roleID := int(req.GetRoleId())
	if roleID <= 0 {
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	role, err := s.roleQuery.FindById(ctx, roleID)
	if err != nil {
		s.logger.Error("FindByIdRole failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindByIdRole success")

	return &pbroles.ApiResponseRole{
		Status:  "success",
		Message: "Successfully fetched role",
		Data:    mapResponseRole(role),
	}, nil
}

func (s *roleQueryHandleGrpc) FindByUserId(ctx context.Context, req *pbroles.FindByIdUserRoleRequest) (*pbroles.ApiResponsesRole, error) {
	s.logger.Info("FindByUserId called", zap.Int32("userId", req.GetUserId()))

	userID := int(req.GetUserId())
	if userID <= 0 {
		return nil, role_errors.ErrGrpcRoleInvalidId
	}

	roles, err := s.roleQuery.FindByUserId(ctx, userID)
	if err != nil {
		s.logger.Error("FindByUserId failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindByUserId success")

	return &pbroles.ApiResponsesRole{
		Status:  "success",
		Message: "Successfully fetched role by user ID",
		Data:    mapResponsesRoleFromDB(roles),
	}, nil
}

func (s *roleQueryHandleGrpc) FindByActive(ctx context.Context, req *pbroles.FindAllRoleRequest) (*pbroles.ApiResponsePaginationRoleDeleteAt, error) {
	s.logger.Info("FindByActive roles called", zap.Int32("page", req.GetPage()))

	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllRoles{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	roles, totalRecords, err := s.roleQuery.FindByActiveRole(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByActive roles failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByActive roles success")

	return &pbroles.ApiResponsePaginationRoleDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active roles",
		Data:       mapResponsesRoleFromActive(roles),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *roleQueryHandleGrpc) FindByTrashed(ctx context.Context, req *pbroles.FindAllRoleRequest) (*pbroles.ApiResponsePaginationRoleDeleteAt, error) {
	s.logger.Info("FindByTrashed roles called")

	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllRoles{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	roles, totalRecords, err := s.roleQuery.FindByTrashedRole(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByTrashed roles failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByTrashed roles success")

	return &pbroles.ApiResponsePaginationRoleDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed roles",
		Data:       mapResponsesRoleFromTrashed(roles),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}
