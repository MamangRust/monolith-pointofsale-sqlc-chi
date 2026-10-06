package handler

import (
	"context"
	pbusers "github.com/MamangRust/monolith-point-of-sale-pb/users"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/user_errors"
	"github.com/MamangRust/monolith-point-of-sale-user/service"
)

type userCommandHandleGrpc struct {
	pbusers.UnimplementedUserCommandServiceServer
	userCommandService service.UserCommandService
	logger             logger.LoggerInterface
}

func NewUserCommandHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pbusers.UserCommandServiceServer {
	return &userCommandHandleGrpc{
		userCommandService: service.UserCommand,
		logger:             logger,
	}
}

func (s *userCommandHandleGrpc) CreateUserRecord(ctx context.Context, request *pbusers.CreateUserRecordRequest) (*pbusers.ApiResponseUser, error) {
	s.logger.Info("CreateUserRecord called", zap.String("email", request.GetEmail()))

	req := &requests.RegisterRequest{
		FirstName:    request.GetFirstname(),
		LastName:     request.GetLastname(),
		Email:        request.GetEmail(),
		Password:     request.GetPassword(),
		VerifiedCode: request.GetVerificationCode(),
		IsVerified:   request.GetIsVerified(),
	}

	user, err := s.userCommandService.CreateUserRecord(ctx, req)
	if err != nil {
		s.logger.Error("CreateUserRecord failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("CreateUserRecord success")

	return &pbusers.ApiResponseUser{
		Status:  "success",
		Message: "Successfully created user",
		Data:    mapUserToProto(user),
	}, nil
}

func (s *userCommandHandleGrpc) UpdateUserIsVerified(ctx context.Context, request *pbusers.UpdateUserIsVerifiedUserRequest) (*pbusers.ApiResponseUser, error) {
	s.logger.Info("UpdateUserIsVerified called", zap.Int32("id", request.GetId()))

	if request.GetId() <= 0 {
		return nil, user_errors.ErrGrpcUserInvalidId
	}

	user, err := s.userCommandService.UpdateUserIsVerified(ctx, int(request.GetId()), request.GetIsVerified())
	if err != nil {
		s.logger.Error("UpdateUserIsVerified failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("UpdateUserIsVerified success")

	return &pbusers.ApiResponseUser{
		Status:  "success",
		Message: "Successfully updated user verification",
		Data:    mapUserToProto(user),
	}, nil
}

func (s *userCommandHandleGrpc) UpdateUserPassword(ctx context.Context, request *pbusers.UpdateUserPasswordUserRequest) (*pbusers.ApiResponseUser, error) {
	s.logger.Info("UpdateUserPassword called", zap.Int32("id", request.GetId()))

	if request.GetId() <= 0 {
		return nil, user_errors.ErrGrpcUserInvalidId
	}
	if request.GetPassword() == "" {
		return nil, errors.ErrBadRequest.WithMessage("Invalid password")
	}

	user, err := s.userCommandService.UpdateUserPassword(ctx, int(request.GetId()), request.GetPassword())
	if err != nil {
		s.logger.Error("UpdateUserPassword failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("UpdateUserPassword success")

	return &pbusers.ApiResponseUser{
		Status:  "success",
		Message: "Successfully updated user password",
		Data:    mapUserToProto(user),
	}, nil
}

func (s *userCommandHandleGrpc) Create(ctx context.Context, request *pbusers.CreateUserRequest) (*pbusers.ApiResponseUser, error) {
	s.logger.Info("Create user called", zap.String("email", request.GetEmail()))

	req := &requests.CreateUserRequest{
		FirstName:       request.GetFirstname(),
		LastName:        request.GetLastname(),
		Email:           request.GetEmail(),
		Password:        request.GetPassword(),
		ConfirmPassword: request.GetConfirmPassword(),
	}

	if err := req.Validate(); err != nil {
		return nil, user_errors.ErrGrpcValidateCreateUser
	}

	user, err := s.userCommandService.CreateUser(ctx, req)
	if err != nil {
		s.logger.Error("Create user failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Create user success")

	return &pbusers.ApiResponseUser{
		Status:  "success",
		Message: "Successfully created user",
		Data:    mapUserToProto(user),
	}, nil
}

func (s *userCommandHandleGrpc) Update(ctx context.Context, request *pbusers.UpdateUserRequest) (*pbusers.ApiResponseUser, error) {
	s.logger.Info("Update user called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, user_errors.ErrGrpcUserInvalidId
	}

	req := &requests.UpdateUserRequest{
		UserID:          &id,
		FirstName:       request.GetFirstname(),
		LastName:        request.GetLastname(),
		Email:           request.GetEmail(),
		Password:        request.GetPassword(),
		ConfirmPassword: request.GetConfirmPassword(),
	}

	if err := req.Validate(); err != nil {
		return nil, user_errors.ErrGrpcValidateCreateUser
	}

	user, err := s.userCommandService.UpdateUser(ctx, req)
	if err != nil {
		s.logger.Error("Update user failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Update user success")

	return &pbusers.ApiResponseUser{
		Status:  "success",
		Message: "Successfully updated user",
		Data:    mapUserToProto(user),
	}, nil
}

func (s *userCommandHandleGrpc) TrashedUser(ctx context.Context, request *pbusers.FindByIdUserRequest) (*pbusers.ApiResponseUserDeleteAt, error) {
	s.logger.Info("TrashedUser called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, user_errors.ErrGrpcUserInvalidId
	}

	user, err := s.userCommandService.TrashedUser(ctx, id)
	if err != nil {
		s.logger.Error("TrashedUser failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("TrashedUser success")

	return &pbusers.ApiResponseUserDeleteAt{
		Status:  "success",
		Message: "Successfully trashed user",
		Data:    mapUserDeleteAtToProto(user),
	}, nil
}

func (s *userCommandHandleGrpc) RestoreUser(ctx context.Context, request *pbusers.FindByIdUserRequest) (*pbusers.ApiResponseUserDeleteAt, error) {
	s.logger.Info("RestoreUser called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, user_errors.ErrGrpcUserInvalidId
	}

	user, err := s.userCommandService.RestoreUser(ctx, id)
	if err != nil {
		s.logger.Error("RestoreUser failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreUser success")

	return &pbusers.ApiResponseUserDeleteAt{
		Status:  "success",
		Message: "Successfully restored user",
		Data:    mapUserDeleteAtToProto(user),
	}, nil
}

func (s *userCommandHandleGrpc) DeleteUserPermanent(ctx context.Context, request *pbusers.FindByIdUserRequest) (*pbusers.ApiResponseUserDelete, error) {
	s.logger.Info("DeleteUserPermanent called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, user_errors.ErrGrpcUserInvalidId
	}

	_, err := s.userCommandService.DeleteUserPermanent(ctx, id)
	if err != nil {
		s.logger.Error("DeleteUserPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteUserPermanent success")

	return &pbusers.ApiResponseUserDelete{
		Status:  "success",
		Message: "Successfully deleted user permanently",
	}, nil
}

func (s *userCommandHandleGrpc) RestoreAllUser(ctx context.Context, _ *emptypb.Empty) (*pbusers.ApiResponseUserAll, error) {
	s.logger.Info("RestoreAllUser called")

	_, err := s.userCommandService.RestoreAllUser(ctx)
	if err != nil {
		s.logger.Error("RestoreAllUser failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreAllUser success")

	return &pbusers.ApiResponseUserAll{
		Status:  "success",
		Message: "Successfully restore all user",
	}, nil
}

func (s *userCommandHandleGrpc) DeleteAllUserPermanent(ctx context.Context, _ *emptypb.Empty) (*pbusers.ApiResponseUserAll, error) {
	s.logger.Info("DeleteAllUserPermanent called")

	_, err := s.userCommandService.DeleteAllUserPermanent(ctx)
	if err != nil {
		s.logger.Error("DeleteAllUserPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteAllUserPermanent success")

	return &pbusers.ApiResponseUserAll{
		Status:  "success",
		Message: "Successfully delete user permanen",
	}, nil
}

// Map helpers
