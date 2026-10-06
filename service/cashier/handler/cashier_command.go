package handler

import (
	"context"
	pbcashiers "github.com/MamangRust/monolith-point-of-sale-pb/cashiers"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/MamangRust/monolith-point-of-sale-cashier/service"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/cashier_errors"
)

type cashierCommandHandleGrpc struct {
	pbcashiers.UnimplementedCashierCommandServiceServer
	cashierCommand service.CashierCommandService
	logger         logger.LoggerInterface
}

func NewCashierCommandHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pbcashiers.CashierCommandServiceServer {
	return &cashierCommandHandleGrpc{
		cashierCommand: service.CashierCommand,
		logger:         logger,
	}
}

func (s *cashierCommandHandleGrpc) CreateCashier(ctx context.Context, request *pbcashiers.CreateCashierRequest) (*pbcashiers.ApiResponseCashier, error) {
	s.logger.Info("CreateCashier called", zap.String("name", request.GetName()))

	req := &requests.CreateCashierRequest{
		Name:       request.GetName(),
		MerchantID: int(request.GetMerchantId()),
		UserID:     int(request.GetUserId()),
	}

	if err := req.Validate(); err != nil {
		return nil, cashier_errors.ErrGrpcValidateCreateCashier
	}

	cashier, err := s.cashierCommand.CreateCashier(ctx, req)
	if err != nil {
		s.logger.Error("CreateCashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("CreateCashier success")

	return &pbcashiers.ApiResponseCashier{
		Status:  "success",
		Message: "Successfully created cashier",
		Data:    mapResponseCashier(cashier),
	}, nil
}

func (s *cashierCommandHandleGrpc) UpdateCashier(ctx context.Context, request *pbcashiers.UpdateCashierRequest) (*pbcashiers.ApiResponseCashier, error) {
	s.logger.Info("UpdateCashier called", zap.Int32("id", request.GetCashierId()))

	id := int(request.GetCashierId())
	if id <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidId
	}

	req := &requests.UpdateCashierRequest{
		CashierID: &id,
		Name:      request.GetName(),
	}

	if err := req.Validate(); err != nil {
		return nil, cashier_errors.ErrGrpcValidateUpdateCashier
	}

	cashier, err := s.cashierCommand.UpdateCashier(ctx, req)
	if err != nil {
		s.logger.Error("UpdateCashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("UpdateCashier success")

	return &pbcashiers.ApiResponseCashier{
		Status:  "success",
		Message: "Successfully updated cashier",
		Data:    mapResponseCashier(cashier),
	}, nil
}

func (s *cashierCommandHandleGrpc) TrashedCashier(ctx context.Context, request *pbcashiers.FindByIdCashierRequest) (*pbcashiers.ApiResponseCashierDeleteAt, error) {
	s.logger.Info("TrashedCashier called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidId
	}

	cashier, err := s.cashierCommand.TrashedCashier(ctx, id)
	if err != nil {
		s.logger.Error("TrashedCashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("TrashedCashier success")

	return &pbcashiers.ApiResponseCashierDeleteAt{
		Status:  "success",
		Message: "Successfully trashed cashier",
		Data:    mapResponseCashierDeleteAt(cashier),
	}, nil
}

func (s *cashierCommandHandleGrpc) RestoreCashier(ctx context.Context, request *pbcashiers.FindByIdCashierRequest) (*pbcashiers.ApiResponseCashierDeleteAt, error) {
	s.logger.Info("RestoreCashier called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidId
	}

	cashier, err := s.cashierCommand.RestoreCashier(ctx, id)
	if err != nil {
		s.logger.Error("RestoreCashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreCashier success")

	return &pbcashiers.ApiResponseCashierDeleteAt{
		Status:  "success",
		Message: "Successfully restored cashier",
		Data:    mapResponseCashierDeleteAt(cashier),
	}, nil
}

func (s *cashierCommandHandleGrpc) DeleteCashierPermanent(ctx context.Context, request *pbcashiers.FindByIdCashierRequest) (*pbcashiers.ApiResponseCashierDelete, error) {
	s.logger.Info("DeleteCashierPermanent called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidId
	}

	_, err := s.cashierCommand.DeleteCashierPermanent(ctx, id)
	if err != nil {
		s.logger.Error("DeleteCashierPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteCashierPermanent success")

	return &pbcashiers.ApiResponseCashierDelete{
		Status:  "success",
		Message: "Successfully deleted cashier permanently",
	}, nil
}

func (s *cashierCommandHandleGrpc) RestoreAllCashier(ctx context.Context, _ *emptypb.Empty) (*pbcashiers.ApiResponseCashierAll, error) {
	s.logger.Info("RestoreAllCashier called")

	_, err := s.cashierCommand.RestoreAllCashier(ctx)
	if err != nil {
		s.logger.Error("RestoreAllCashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreAllCashier success")

	return &pbcashiers.ApiResponseCashierAll{
		Status:  "success",
		Message: "Successfully restore all cashier",
	}, nil
}

func (s *cashierCommandHandleGrpc) DeleteAllCashierPermanent(ctx context.Context, _ *emptypb.Empty) (*pbcashiers.ApiResponseCashierAll, error) {
	s.logger.Info("DeleteAllCashierPermanent called")

	_, err := s.cashierCommand.DeleteAllCashierPermanent(ctx)
	if err != nil {
		s.logger.Error("DeleteAllCashierPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteAllCashierPermanent success")

	return &pbcashiers.ApiResponseCashierAll{
		Status:  "success",
		Message: "Successfully delete cashier permanen",
	}, nil
}

// Map helpers
