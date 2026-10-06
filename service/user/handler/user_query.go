package handler

import (
	"context"
	pbcommon "github.com/MamangRust/monolith-point-of-sale-pb/common"
	pbusers "github.com/MamangRust/monolith-point-of-sale-pb/users"
	"go.uber.org/zap"
	"math"

	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/user_errors"
	"github.com/MamangRust/monolith-point-of-sale-user/service"
)

type userQueryHandleGrpc struct {
	pbusers.UnimplementedUserQueryServiceServer
	userQueryService service.UserQueryService
	logger           logger.LoggerInterface
}

func NewUserQueryHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pbusers.UserQueryServiceServer {
	return &userQueryHandleGrpc{
		userQueryService: service.UserQuery,
		logger:           logger,
	}
}

func (s *userQueryHandleGrpc) FindByEmail(ctx context.Context, request *pbusers.FindByEmailUserRequest) (*pbusers.ApiResponseUserWithPassword, error) {
	s.logger.Info("FindByEmail user called", zap.String("email", request.GetEmail()))

	if request.GetEmail() == "" {
		return nil, errors.ErrBadRequest.WithMessage("Invalid email")
	}

	user, err := s.userQueryService.FindByEmailWithPassword(ctx, request.GetEmail())
	if err != nil {
		s.logger.Error("FindByEmail user failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindByEmail user success")

	return &pbusers.ApiResponseUserWithPassword{
		Status:  "success",
		Message: "Successfully fetched user",
		Data:    mapUserWithPasswordToProto(user),
	}, nil
}

func (s *userQueryHandleGrpc) FindByEmailAndVerify(ctx context.Context, request *pbusers.FindByEmailUserRequest) (*pbusers.ApiResponseUserWithPassword, error) {
	s.logger.Info("FindByEmailAndVerify user called", zap.String("email", request.GetEmail()))

	if request.GetEmail() == "" {
		return nil, errors.ErrBadRequest.WithMessage("Invalid email")
	}

	user, err := s.userQueryService.FindByEmailAndVerify(ctx, request.GetEmail())
	if err != nil {
		s.logger.Error("FindByEmailAndVerify user failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindByEmailAndVerify user success")

	return &pbusers.ApiResponseUserWithPassword{
		Status:  "success",
		Message: "Successfully fetched verified user",
		Data:    mapUserWithPasswordToProto(user),
	}, nil
}

func (s *userQueryHandleGrpc) FindByVerificationCode(ctx context.Context, request *pbusers.FindByVerificationCodeUserRequest) (*pbusers.ApiResponseUser, error) {
	s.logger.Info("FindByVerificationCode user called")

	if request.GetVerificationCode() == "" {
		return nil, errors.ErrBadRequest.WithMessage("Invalid verification code")
	}

	user, err := s.userQueryService.FindByVerificationCode(ctx, request.GetVerificationCode())
	if err != nil {
		s.logger.Error("FindByVerificationCode user failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindByVerificationCode user success")

	return &pbusers.ApiResponseUser{
		Status:  "success",
		Message: "Successfully fetched user",
		Data:    mapUserToProto(user),
	}, nil
}

func (s *userQueryHandleGrpc) FindAll(ctx context.Context, request *pbusers.FindAllUserRequest) (*pbusers.ApiResponsePaginationUser, error) {
	s.logger.Info("FindAll users called", zap.Int32("page", request.GetPage()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllUsers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.userQueryService.FindAll(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindAll users failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	var totalRecordsVal int
	if totalRecords != nil {
		totalRecordsVal = *totalRecords
	}
	totalPages := int(math.Ceil(float64(totalRecordsVal) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecordsVal),
	}

	s.logger.Info("FindAll users success")

	return &pbusers.ApiResponsePaginationUser{
		Status:     "success",
		Message:    "Successfully fetched users",
		Data:       mapGetUsersRowsToProto(users),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *userQueryHandleGrpc) FindById(ctx context.Context, request *pbusers.FindByIdUserRequest) (*pbusers.ApiResponseUser, error) {
	s.logger.Info("FindById user called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, user_errors.ErrGrpcUserNotFound
	}

	user, err := s.userQueryService.FindByID(ctx, id)
	if err != nil {
		s.logger.Error("FindById user failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindById user success")

	return &pbusers.ApiResponseUser{
		Status:  "success",
		Message: "Successfully fetched user",
		Data:    mapUserToProto(user),
	}, nil
}

func (s *userQueryHandleGrpc) FindByActive(ctx context.Context, request *pbusers.FindAllUserRequest) (*pbusers.ApiResponsePaginationUserDeleteAt, error) {
	s.logger.Info("FindByActive users called", zap.Int32("page", request.GetPage()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllUsers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.userQueryService.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByActive users failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	var totalRecordsVal int
	if totalRecords != nil {
		totalRecordsVal = *totalRecords
	}
	totalPages := int(math.Ceil(float64(totalRecordsVal) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecordsVal),
	}

	s.logger.Info("FindByActive users success")

	return &pbusers.ApiResponsePaginationUserDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active users",
		Data:       mapActiveUsersToProto(users),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *userQueryHandleGrpc) FindByTrashed(ctx context.Context, request *pbusers.FindAllUserRequest) (*pbusers.ApiResponsePaginationUserDeleteAt, error) {
	s.logger.Info("FindByTrashed users called")

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllUsers{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.userQueryService.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByTrashed users failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	var totalRecordsVal int
	if totalRecords != nil {
		totalRecordsVal = *totalRecords
	}
	totalPages := int(math.Ceil(float64(totalRecordsVal) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(totalRecordsVal),
	}

	s.logger.Info("FindByTrashed users success")

	return &pbusers.ApiResponsePaginationUserDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed users",
		Data:       mapTrashedUsersToProto(users),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}
