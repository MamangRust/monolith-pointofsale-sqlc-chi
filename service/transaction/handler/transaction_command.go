package handler

import (
	"context"
	pbtransactions "github.com/MamangRust/monolith-point-of-sale-pb/transactions"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/transaction_errors"
	"github.com/MamangRust/monolith-point-of-sale-transacton/service"
)

type transactionCommandHandleGrpc struct {
	pbtransactions.UnimplementedTransactionCommandServiceServer
	transactionCommand service.TransactionCommandService
	logger             logger.LoggerInterface
}

func NewTransactionCommandHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pbtransactions.TransactionCommandServiceServer {
	return &transactionCommandHandleGrpc{
		transactionCommand: service.TransactionCommand,
		logger:             logger,
	}
}

func (s *transactionCommandHandleGrpc) Create(ctx context.Context, request *pbtransactions.CreateTransactionRequest) (*pbtransactions.ApiResponseTransaction, error) {
	s.logger.Info("Create transaction called", zap.Int32("orderId", request.GetOrderId()))

	req := &requests.CreateTransactionRequest{
		CashierID:     int(request.GetCashierId()),
		OrderID:       int(request.GetOrderId()),
		PaymentMethod: request.GetPaymentMethod(),
		Amount:        int(request.GetAmount()),
	}

	if err := req.Validate(); err != nil {
		return nil, transaction_errors.ErrGrpcValidateCreateTransaction
	}

	transaction, err := s.transactionCommand.CreateTransaction(ctx, req)
	if err != nil {
		s.logger.Error("Create transaction failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Create transaction success")

	return &pbtransactions.ApiResponseTransaction{
		Status:  "success",
		Message: "Successfully created transaction",
		Data:    mapResponseTransaction(transaction),
	}, nil
}

func (s *transactionCommandHandleGrpc) Update(ctx context.Context, request *pbtransactions.UpdateTransactionRequest) (*pbtransactions.ApiResponseTransaction, error) {
	s.logger.Info("Update transaction called", zap.Int32("id", request.GetTransactionId()))

	id := int(request.GetTransactionId())
	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidID
	}

	req := &requests.UpdateTransactionRequest{
		TransactionID: &id,
		OrderID:       int(request.GetOrderId()),
		CashierID:     int(request.GetCashierId()),
		PaymentMethod: request.GetPaymentMethod(),
		Amount:        int(request.GetAmount()),
	}

	if err := req.Validate(); err != nil {
		return nil, transaction_errors.ErrGrpcValidateUpdateTransaction
	}

	transaction, err := s.transactionCommand.UpdateTransaction(ctx, req)
	if err != nil {
		s.logger.Error("Update transaction failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Update transaction success")

	return &pbtransactions.ApiResponseTransaction{
		Status:  "success",
		Message: "Successfully updated transaction",
		Data:    mapResponseTransaction(transaction),
	}, nil
}

func (s *transactionCommandHandleGrpc) TrashedTransaction(ctx context.Context, request *pbtransactions.FindByIdTransactionRequest) (*pbtransactions.ApiResponseTransactionDeleteAt, error) {
	s.logger.Info("TrashedTransaction called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidID
	}

	transaction, err := s.transactionCommand.TrashedTransaction(ctx, id)
	if err != nil {
		s.logger.Error("TrashedTransaction failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("TrashedTransaction success")

	return &pbtransactions.ApiResponseTransactionDeleteAt{
		Status:  "success",
		Message: "Successfully trashed transaction",
		Data:    mapResponseTransactionDeleteAt(transaction),
	}, nil
}

func (s *transactionCommandHandleGrpc) RestoreTransaction(ctx context.Context, request *pbtransactions.FindByIdTransactionRequest) (*pbtransactions.ApiResponseTransactionDeleteAt, error) {
	s.logger.Info("RestoreTransaction called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidID
	}

	transaction, err := s.transactionCommand.RestoreTransaction(ctx, id)
	if err != nil {
		s.logger.Error("RestoreTransaction failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreTransaction success")

	return &pbtransactions.ApiResponseTransactionDeleteAt{
		Status:  "success",
		Message: "Successfully restored transaction",
		Data:    mapResponseTransactionDeleteAt(transaction),
	}, nil
}

func (s *transactionCommandHandleGrpc) DeleteTransactionPermanent(ctx context.Context, request *pbtransactions.FindByIdTransactionRequest) (*pbtransactions.ApiResponseTransactionDelete, error) {
	s.logger.Info("DeleteTransactionPermanent called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidID
	}

	_, err := s.transactionCommand.DeleteTransactionPermanently(ctx, id)
	if err != nil {
		s.logger.Error("DeleteTransactionPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteTransactionPermanent success")

	return &pbtransactions.ApiResponseTransactionDelete{
		Status:  "success",
		Message: "Successfully deleted Transaction permanently",
	}, nil
}

func (s *transactionCommandHandleGrpc) RestoreAllTransaction(ctx context.Context, _ *emptypb.Empty) (*pbtransactions.ApiResponseTransactionAll, error) {
	s.logger.Info("RestoreAllTransaction called")

	_, err := s.transactionCommand.RestoreAllTransactions(ctx)
	if err != nil {
		s.logger.Error("RestoreAllTransaction failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreAllTransaction success")

	return &pbtransactions.ApiResponseTransactionAll{
		Status:  "success",
		Message: "Successfully restore all Transaction",
	}, nil
}

func (s *transactionCommandHandleGrpc) DeleteAllTransactionPermanent(ctx context.Context, _ *emptypb.Empty) (*pbtransactions.ApiResponseTransactionAll, error) {
	s.logger.Info("DeleteAllTransactionPermanent called")

	_, err := s.transactionCommand.DeleteAllTransactionPermanent(ctx)
	if err != nil {
		s.logger.Error("DeleteAllTransactionPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteAllTransactionPermanent success")

	return &pbtransactions.ApiResponseTransactionAll{
		Status:  "success",
		Message: "Successfully delete Transaction permanen",
	}, nil
}

// Map helpers
