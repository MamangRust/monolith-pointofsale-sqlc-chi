package handler

import (
	"context"
	pbcashiers "github.com/MamangRust/monolith-point-of-sale-pb/cashiers"
	"go.uber.org/zap"

	"github.com/MamangRust/monolith-point-of-sale-cashier/service"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/cashier_errors"
)

type cashierStatsHandleGrpc struct {
	pbcashiers.UnimplementedCashierStatsServiceServer
	cashierStats service.CashierStatsService
	logger       logger.LoggerInterface
}

func NewCashierStatsHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pbcashiers.CashierStatsServiceServer {
	return &cashierStatsHandleGrpc{
		cashierStats: service.CashierStats,
		logger:       logger,
	}
}

func (s *cashierStatsHandleGrpc) FindMonthlyTotalSales(ctx context.Context, req *pbcashiers.FindYearMonthTotalSales) (*pbcashiers.ApiResponseCashierMonthlyTotalSales, error) {
	s.logger.Info("FindMonthlyTotalSales cashier called", zap.Int32("year", req.GetYear()), zap.Int32("month", req.GetMonth()))

	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidYear
	}
	if month <= 0 || month >= 12 {
		return nil, cashier_errors.ErrGrpcFailedInvalidMonth
	}

	reqService := requests.MonthTotalSales{
		Year:  year,
		Month: month,
	}

	methods, err := s.cashierStats.FindMonthlyTotalSales(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthlyTotalSales cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthlyTotalSales cashier success")

	return &pbcashiers.ApiResponseCashierMonthlyTotalSales{
		Status:  "success",
		Message: "Monthly sales retrieved successfully",
		Data:    mapResponseCashierMonthlyTotalSales(methods),
	}, nil
}

func (s *cashierStatsHandleGrpc) FindYearlyTotalSales(ctx context.Context, req *pbcashiers.FindYearTotalSales) (*pbcashiers.ApiResponseCashierYearlyTotalSales, error) {
	s.logger.Info("FindYearlyTotalSales cashier called", zap.Int32("year", req.GetYear()))

	year := int(req.GetYear())
	if year <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidYear
	}

	methods, err := s.cashierStats.FindYearlyTotalSales(ctx, year)
	if err != nil {
		s.logger.Error("FindYearlyTotalSales cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearlyTotalSales cashier success")

	return &pbcashiers.ApiResponseCashierYearlyTotalSales{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapResponseCashierYearlyTotalSales(methods),
	}, nil
}

func (s *cashierStatsHandleGrpc) FindMonthSales(ctx context.Context, req *pbcashiers.FindYearCashier) (*pbcashiers.ApiResponseCashierMonthSales, error) {
	s.logger.Info("FindMonthSales cashier called", zap.Int32("year", req.GetYear()))

	year := int(req.GetYear())
	if year <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidYear
	}

	methods, err := s.cashierStats.FindMonthlySales(ctx, year)
	if err != nil {
		s.logger.Error("FindMonthSales cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthSales cashier success")

	return &pbcashiers.ApiResponseCashierMonthSales{
		Status:  "success",
		Message: "Monthly sales retrieved successfully",
		Data:    mapResponsesCashierMonthlySales(methods),
	}, nil
}

func (s *cashierStatsHandleGrpc) FindYearSales(ctx context.Context, req *pbcashiers.FindYearCashier) (*pbcashiers.ApiResponseCashierYearSales, error) {
	s.logger.Info("FindYearSales cashier called", zap.Int32("year", req.GetYear()))

	year := int(req.GetYear())
	if year <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidYear
	}

	methods, err := s.cashierStats.FindYearlySales(ctx, year)
	if err != nil {
		s.logger.Error("FindYearSales cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearSales cashier success")

	return &pbcashiers.ApiResponseCashierYearSales{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapResponsesCashierYearlySales(methods),
	}, nil
}
