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

type cashierStatsByIdHandleGrpc struct {
	pbcashiers.UnimplementedCashierStatsByIdServiceServer
	cashierStatsById service.CashierStatsByIdService
	logger           logger.LoggerInterface
}

func NewCashierStatsByIdHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pbcashiers.CashierStatsByIdServiceServer {
	return &cashierStatsByIdHandleGrpc{
		cashierStatsById: service.CashierStatsById,
		logger:           logger,
	}
}

func (s *cashierStatsByIdHandleGrpc) FindMonthlyTotalSalesById(ctx context.Context, req *pbcashiers.FindYearMonthTotalSalesById) (*pbcashiers.ApiResponseCashierMonthlyTotalSales, error) {
	s.logger.Info("FindMonthlyTotalSalesById cashier called", zap.Int32("id", req.GetCashierId()))

	year := int(req.GetYear())
	month := int(req.GetMonth())
	id := int(req.GetCashierId())

	if year <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidYear
	}
	if month <= 0 || month >= 12 {
		return nil, cashier_errors.ErrGrpcFailedInvalidMonth
	}
	if id <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidId
	}

	reqService := requests.MonthTotalSalesCashier{
		Year:      year,
		Month:     month,
		CashierID: id,
	}

	methods, err := s.cashierStatsById.FindMonthlyTotalSalesById(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthlyTotalSalesById cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthlyTotalSalesById cashier success")

	return &pbcashiers.ApiResponseCashierMonthlyTotalSales{
		Status:  "success",
		Message: "Monthly sales retrieved successfully",
		Data:    mapResponseCashierMonthlyTotalSalesById(methods),
	}, nil
}

func (s *cashierStatsByIdHandleGrpc) FindYearlyTotalSalesById(ctx context.Context, req *pbcashiers.FindYearTotalSalesById) (*pbcashiers.ApiResponseCashierYearlyTotalSales, error) {
	s.logger.Info("FindYearlyTotalSalesById cashier called", zap.Int32("id", req.GetCashierId()))

	year := int(req.GetYear())
	id := int(req.GetCashierId())

	if year <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidYear
	}
	if id <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidId
	}

	reqService := requests.YearTotalSalesCashier{
		Year:      year,
		CashierID: id,
	}

	methods, err := s.cashierStatsById.FindYearlyTotalSalesById(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindYearlyTotalSalesById cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearlyTotalSalesById cashier success")

	return &pbcashiers.ApiResponseCashierYearlyTotalSales{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapResponseCashierYearlyTotalSalesById(methods),
	}, nil
}

func (s *cashierStatsByIdHandleGrpc) FindMonthSalesById(ctx context.Context, req *pbcashiers.FindYearCashierById) (*pbcashiers.ApiResponseCashierMonthSales, error) {
	s.logger.Info("FindMonthSalesById cashier called", zap.Int32("id", req.GetCashierId()))

	year := int(req.GetYear())
	cashierId := int(req.GetCashierId())

	if year <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidYear
	}
	if cashierId <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidId
	}

	reqService := requests.MonthCashierId{
		Year:      year,
		CashierID: cashierId,
	}

	methods, err := s.cashierStatsById.FindMonthlyCashierById(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthSalesById cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthSalesById cashier success")

	return &pbcashiers.ApiResponseCashierMonthSales{
		Status:  "success",
		Message: "Cashier monthly sales retrieved successfully",
		Data:    mapResponsesCashierMonthlySalesById(methods),
	}, nil
}

func (s *cashierStatsByIdHandleGrpc) FindYearSalesById(ctx context.Context, req *pbcashiers.FindYearCashierById) (*pbcashiers.ApiResponseCashierYearSales, error) {
	s.logger.Info("FindYearSalesById cashier called", zap.Int32("id", req.GetCashierId()))

	year := int(req.GetYear())
	cashierId := int(req.GetCashierId())

	if year <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidYear
	}
	if cashierId <= 0 {
		return nil, cashier_errors.ErrGrpcFailedInvalidId
	}

	reqService := requests.YearCashierId{
		Year:      year,
		CashierID: cashierId,
	}

	methods, err := s.cashierStatsById.FindYearlyCashierById(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindYearSalesById cashier failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearSalesById cashier success")

	return &pbcashiers.ApiResponseCashierYearSales{
		Status:  "success",
		Message: "Cashier yearly sales retrieved successfully",
		Data:    mapResponsesCashierYearlySalesById(methods),
	}, nil
}
