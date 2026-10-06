package handler

import (
	"context"
	pbmerchants "github.com/MamangRust/monolith-point-of-sale-pb/merchants"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/MamangRust/monolith-point-of-sale-merchant/service"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/merchant_errors"
)

type merchantCommandHandleGrpc struct {
	pbmerchants.UnimplementedMerchantCommandServiceServer
	merchantCommand service.MerchantCommandService
	logger          logger.LoggerInterface
}

func NewMerchantCommandHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pbmerchants.MerchantCommandServiceServer {
	return &merchantCommandHandleGrpc{
		merchantCommand: service.MerchantCommand,
		logger:          logger,
	}
}

func (s *merchantCommandHandleGrpc) Create(ctx context.Context, request *pbmerchants.CreateMerchantRequest) (*pbmerchants.ApiResponseMerchant, error) {
	s.logger.Info("Create merchant called", zap.String("name", request.GetName()))

	req := &requests.CreateMerchantRequest{
		UserID:       int(request.GetUserId()),
		Name:         request.GetName(),
		Description:  request.GetDescription(),
		Address:      request.GetAddress(),
		ContactEmail: request.GetContactEmail(),
		ContactPhone: request.GetContactPhone(),
		Status:       request.GetStatus(),
	}

	if err := req.Validate(); err != nil {
		return nil, merchant_errors.ErrGrpcValidateCreateMerchant
	}

	merchant, err := s.merchantCommand.CreateMerchant(ctx, req)
	if err != nil {
		s.logger.Error("Create merchant failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Create merchant success")

	return &pbmerchants.ApiResponseMerchant{
		Status:  "success",
		Message: "Successfully created merchant",
		Data:    mapResponseMerchant(merchant),
	}, nil
}

func (s *merchantCommandHandleGrpc) Update(ctx context.Context, request *pbmerchants.UpdateMerchantRequest) (*pbmerchants.ApiResponseMerchant, error) {
	s.logger.Info("Update merchant called", zap.Int32("id", request.GetMerchantId()))

	id := int(request.GetMerchantId())
	if id <= 0 {
		return nil, merchant_errors.ErrGrpcInvalidID
	}

	req := &requests.UpdateMerchantRequest{
		MerchantID:   &id,
		UserID:       int(request.GetUserId()),
		Name:         request.GetName(),
		Description:  request.GetDescription(),
		Address:      request.GetAddress(),
		ContactEmail: request.GetContactEmail(),
		ContactPhone: request.GetContactPhone(),
		Status:       request.GetStatus(),
	}

	if err := req.Validate(); err != nil {
		return nil, merchant_errors.ErrGrpcValidateUpdateMerchant
	}

	merchant, err := s.merchantCommand.UpdateMerchant(ctx, req)
	if err != nil {
		s.logger.Error("Update merchant failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Update merchant success")

	return &pbmerchants.ApiResponseMerchant{
		Status:  "success",
		Message: "Successfully updated merchant",
		Data:    mapResponseMerchant(merchant),
	}, nil
}

func (s *merchantCommandHandleGrpc) UpdateMerchantStatus(ctx context.Context, req *pbmerchants.UpdateMerchantStatusRequest) (*pbmerchants.ApiResponseMerchant, error) {
	s.logger.Info("UpdateMerchantStatus merchant called", zap.Int32("id", req.GetMerchantId()))

	id := int(req.GetMerchantId())
	if id <= 0 {
		return nil, merchant_errors.ErrGrpcInvalidID
	}

	request := requests.UpdateMerchantStatusRequest{
		MerchantID: &id,
		Status:     req.GetStatus(),
	}

	if err := request.Validate(); err != nil {
		return nil, merchant_errors.ErrGrpcValidateUpdateMerchantStatus
	}

	merchant, err := s.merchantCommand.UpdateMerchantStatus(ctx, &request)
	if err != nil {
		s.logger.Error("UpdateMerchantStatus merchant failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("UpdateMerchantStatus merchant success")

	return &pbmerchants.ApiResponseMerchant{
		Status:  "success",
		Message: "Successfully updated merchant status",
		Data:    mapResponseMerchant(merchant),
	}, nil
}

func (s *merchantCommandHandleGrpc) TrashedMerchant(ctx context.Context, request *pbmerchants.FindByIdMerchantRequest) (*pbmerchants.ApiResponseMerchantDeleteAt, error) {
	s.logger.Info("TrashedMerchant called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, merchant_errors.ErrGrpcInvalidID
	}

	merchant, err := s.merchantCommand.TrashedMerchant(ctx, id)
	if err != nil {
		s.logger.Error("TrashedMerchant failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("TrashedMerchant success")

	return &pbmerchants.ApiResponseMerchantDeleteAt{
		Status:  "success",
		Message: "Successfully trashed merchant",
		Data:    mapResponseMerchantDeleteAt(merchant),
	}, nil
}

func (s *merchantCommandHandleGrpc) RestoreMerchant(ctx context.Context, request *pbmerchants.FindByIdMerchantRequest) (*pbmerchants.ApiResponseMerchant, error) {
	s.logger.Info("RestoreMerchant called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, merchant_errors.ErrGrpcInvalidID
	}

	merchant, err := s.merchantCommand.RestoreMerchant(ctx, id)
	if err != nil {
		s.logger.Error("RestoreMerchant failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreMerchant success")

	return &pbmerchants.ApiResponseMerchant{
		Status:  "success",
		Message: "Successfully restored merchant",
		Data:    mapResponseMerchant(merchant),
	}, nil
}

func (s *merchantCommandHandleGrpc) DeleteMerchantPermanent(ctx context.Context, request *pbmerchants.FindByIdMerchantRequest) (*pbmerchants.ApiResponseMerchantDelete, error) {
	s.logger.Info("DeleteMerchantPermanent called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, merchant_errors.ErrGrpcInvalidID
	}

	_, err := s.merchantCommand.DeleteMerchantPermanent(ctx, id)
	if err != nil {
		s.logger.Error("DeleteMerchantPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteMerchantPermanent success")

	return &pbmerchants.ApiResponseMerchantDelete{
		Status:  "success",
		Message: "Successfully deleted merchant permanently",
	}, nil
}

func (s *merchantCommandHandleGrpc) RestoreAllMerchant(ctx context.Context, _ *emptypb.Empty) (*pbmerchants.ApiResponseMerchantAll, error) {
	s.logger.Info("RestoreAllMerchant called")

	_, err := s.merchantCommand.RestoreAllMerchant(ctx)
	if err != nil {
		s.logger.Error("RestoreAllMerchant failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreAllMerchant success")

	return &pbmerchants.ApiResponseMerchantAll{
		Status:  "success",
		Message: "Successfully restore all merchant",
	}, nil
}

func (s *merchantCommandHandleGrpc) DeleteAllMerchantPermanent(ctx context.Context, _ *emptypb.Empty) (*pbmerchants.ApiResponseMerchantAll, error) {
	s.logger.Info("DeleteAllMerchantPermanent called")

	_, err := s.merchantCommand.DeleteAllMerchantPermanent(ctx)
	if err != nil {
		s.logger.Error("DeleteAllMerchantPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteAllMerchantPermanent success")

	return &pbmerchants.ApiResponseMerchantAll{
		Status:  "success",
		Message: "Successfully delete merchant permanen",
	}, nil
}

// Map helpers
