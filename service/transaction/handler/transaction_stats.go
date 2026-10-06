package handler

import (
	"context"
	pbtransactions "github.com/MamangRust/monolith-point-of-sale-pb/transactions"
	"go.uber.org/zap"

	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/transaction_errors"
	"github.com/MamangRust/monolith-point-of-sale-transacton/service"
)

type transactionStatsHandleGrpc struct {
	pbtransactions.UnimplementedTransactionStatsServiceServer
	transactionStats service.TransactionStatsService
	logger           logger.LoggerInterface
}

func NewTransactionStatsHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pbtransactions.TransactionStatsServiceServer {
	return &transactionStatsHandleGrpc{
		transactionStats: service.TransactionStats,
		logger:           logger,
	}
}

func (s *transactionStatsHandleGrpc) FindMonthStatusSuccess(ctx context.Context, request *pbtransactions.FindMonthlyTransactionStatus) (*pbtransactions.ApiResponseTransactionMonthAmountSuccess, error) {
	s.logger.Info("FindMonthStatusSuccess transactions called", zap.Int32("year", request.GetYear()))

	year := int(request.GetYear())
	month := int(request.GetMonth())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}
	if month <= 0 || month >= 12 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	reqService := requests.MonthAmountTransaction{
		Year:  year,
		Month: month,
	}

	res, err := s.transactionStats.FindMonthlyAmountSuccess(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthStatusSuccess transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthStatusSuccess transactions success")

	return &pbtransactions.ApiResponseTransactionMonthAmountSuccess{
		Status:  "success",
		Message: "Monthly success data retrieved successfully",
		Data:    mapResponsesTransactionMonthlyAmountSuccess(res),
	}, nil
}

func (s *transactionStatsHandleGrpc) FindYearStatusSuccess(ctx context.Context, request *pbtransactions.FindYearlyTransactionStatus) (*pbtransactions.ApiResponseTransactionYearAmountSuccess, error) {
	s.logger.Info("FindYearStatusSuccess transactions called", zap.Int32("year", request.GetYear()))

	year := int(request.GetYear())
	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	res, err := s.transactionStats.FindYearlyAmountSuccess(ctx, year)
	if err != nil {
		s.logger.Error("FindYearStatusSuccess transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearStatusSuccess transactions success")

	return &pbtransactions.ApiResponseTransactionYearAmountSuccess{
		Status:  "success",
		Message: "Yearly success data retrieved successfully",
		Data:    mapResponsesTransactionYearlyAmountSuccess(res),
	}, nil
}

func (s *transactionStatsHandleGrpc) FindMonthStatusFailed(ctx context.Context, request *pbtransactions.FindMonthlyTransactionStatus) (*pbtransactions.ApiResponseTransactionMonthAmountFailed, error) {
	s.logger.Info("FindMonthStatusFailed transactions called", zap.Int32("year", request.GetYear()))

	year := int(request.GetYear())
	month := int(request.GetMonth())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}
	if month <= 0 || month >= 12 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	reqService := requests.MonthAmountTransaction{
		Year:  year,
		Month: month,
	}

	res, err := s.transactionStats.FindMonthlyAmountFailed(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthStatusFailed transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthStatusFailed transactions success")

	return &pbtransactions.ApiResponseTransactionMonthAmountFailed{
		Status:  "success",
		Message: "Monthly failed data retrieved successfully",
		Data:    mapResponsesTransactionMonthlyAmountFailed(res),
	}, nil
}

func (s *transactionStatsHandleGrpc) FindYearStatusFailed(ctx context.Context, request *pbtransactions.FindYearlyTransactionStatus) (*pbtransactions.ApiResponseTransactionYearAmountFailed, error) {
	s.logger.Info("FindYearStatusFailed transactions called", zap.Int32("year", request.GetYear()))

	year := int(request.GetYear())
	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	res, err := s.transactionStats.FindYearlyAmountFailed(ctx, year)
	if err != nil {
		s.logger.Error("FindYearStatusFailed transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearStatusFailed transactions success")

	return &pbtransactions.ApiResponseTransactionYearAmountFailed{
		Status:  "success",
		Message: "Yearly failed data retrieved successfully",
		Data:    mapResponsesTransactionYearlyAmountFailed(res),
	}, nil
}

func (s *transactionStatsHandleGrpc) FindMonthMethodSuccess(ctx context.Context, req *pbtransactions.MonthTransactionMethod) (*pbtransactions.ApiResponseTransactionMonthPaymentMethod, error) {
	s.logger.Info("FindMonthMethodSuccess transactions called", zap.Int32("year", req.GetYear()))

	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}
	if month <= 0 || month >= 12 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	methods, err := s.transactionStats.FindMonthlyMethodSuccess(ctx, &requests.MonthMethodTransaction{
		Year:  year,
		Month: month,
	})
	if err != nil {
		s.logger.Error("FindMonthMethodSuccess transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthMethodSuccess transactions success")

	return &pbtransactions.ApiResponseTransactionMonthPaymentMethod{
		Status:  "success",
		Message: "Monthly payment methods retrieved successfully",
		Data:    mapResponsesTransactionMonthlyMethodSuccess(methods),
	}, nil
}

func (s *transactionStatsHandleGrpc) FindYearMethodSuccess(ctx context.Context, req *pbtransactions.YearTransactionMethod) (*pbtransactions.ApiResponseTransactionYearPaymentmethod, error) {
	s.logger.Info("FindYearMethodSuccess transactions called", zap.Int32("year", req.GetYear()))

	year := int(req.GetYear())
	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	methods, err := s.transactionStats.FindYearlyMethodSuccess(ctx, year)
	if err != nil {
		s.logger.Error("FindYearMethodSuccess transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearMethodSuccess transactions success")

	return &pbtransactions.ApiResponseTransactionYearPaymentmethod{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapResponsesTransactionYearlyMethodSuccess(methods),
	}, nil
}

func (s *transactionStatsHandleGrpc) FindMonthMethodFailed(ctx context.Context, req *pbtransactions.MonthTransactionMethod) (*pbtransactions.ApiResponseTransactionMonthPaymentMethod, error) {
	s.logger.Info("FindMonthMethodFailed transactions called", zap.Int32("year", req.GetYear()))

	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}
	if month <= 0 || month >= 12 {
		return nil, transaction_errors.ErrGrpcInvalidMonth
	}

	methods, err := s.transactionStats.FindMonthlyMethodFailed(ctx, &requests.MonthMethodTransaction{
		Year:  year,
		Month: month,
	})
	if err != nil {
		s.logger.Error("FindMonthMethodFailed transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthMethodFailed transactions success")

	return &pbtransactions.ApiResponseTransactionMonthPaymentMethod{
		Status:  "Failed",
		Message: "Monthly payment methods retrieved Failedfully",
		Data:    mapResponsesTransactionMonthlyMethodFailed(methods),
	}, nil
}

func (s *transactionStatsHandleGrpc) FindYearMethodFailed(ctx context.Context, req *pbtransactions.YearTransactionMethod) (*pbtransactions.ApiResponseTransactionYearPaymentmethod, error) {
	s.logger.Info("FindYearMethodFailed transactions called", zap.Int32("year", req.GetYear()))

	year := int(req.GetYear())
	if year <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidYear
	}

	methods, err := s.transactionStats.FindYearlyMethodFailed(ctx, year)
	if err != nil {
		s.logger.Error("FindYearMethodFailed transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearMethodFailed transactions success")

	return &pbtransactions.ApiResponseTransactionYearPaymentmethod{
		Status:  "Failed",
		Message: "Yearly payment methods retrieved Failedfully",
		Data:    mapResponsesTransactionYearlyMethodFailed(methods),
	}, nil
}
