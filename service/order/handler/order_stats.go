package handler

import (
	"context"
	pborders "github.com/MamangRust/monolith-point-of-sale-pb/orders"
	"go.uber.org/zap"

	"github.com/MamangRust/monolith-point-of-sale-order/service"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/order_errors"
)

type orderStatsHandleGrpc struct {
	pborders.UnimplementedOrderStatsServiceServer
	orderStats service.OrderStatsService
	logger     logger.LoggerInterface
}

func NewOrderStatsHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pborders.OrderStatsServiceServer {
	return &orderStatsHandleGrpc{
		orderStats: service.OrderStats,
		logger:     logger,
	}
}

func (s *orderStatsHandleGrpc) FindMonthlyTotalRevenue(ctx context.Context, req *pborders.FindYearMonthTotalRevenue) (*pborders.ApiResponseOrderMonthlyTotalRevenue, error) {
	s.logger.Info("FindMonthlyTotalRevenue orders called", zap.Int32("year", req.GetYear()))

	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, order_errors.ErrGrpcInvalidYear
	}
	if month <= 0 || month >= 12 {
		return nil, order_errors.ErrGrpcInvalidMonth
	}

	reqService := requests.MonthTotalRevenue{
		Year:  year,
		Month: month,
	}

	methods, err := s.orderStats.FindMonthlyTotalRevenue(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthlyTotalRevenue orders failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthlyTotalRevenue orders success")

	return &pborders.ApiResponseOrderMonthlyTotalRevenue{
		Status:  "success",
		Message: "Monthly sales retrieved successfully",
		Data:    mapResponseOrderMonthlyTotalRevenues(methods),
	}, nil
}

func (s *orderStatsHandleGrpc) FindYearlyTotalRevenue(ctx context.Context, req *pborders.FindYearTotalRevenue) (*pborders.ApiResponseOrderYearlyTotalRevenue, error) {
	s.logger.Info("FindYearlyTotalRevenue orders called", zap.Int32("year", req.GetYear()))

	year := int(req.GetYear())
	if year <= 0 {
		return nil, order_errors.ErrGrpcInvalidYear
	}

	methods, err := s.orderStats.FindYearlyTotalRevenue(ctx, year)
	if err != nil {
		s.logger.Error("FindYearlyTotalRevenue orders failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearlyTotalRevenue orders success")

	return &pborders.ApiResponseOrderYearlyTotalRevenue{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapResponseOrderYearlyTotalRevenues(methods),
	}, nil
}

func (s *orderStatsHandleGrpc) FindMonthlyRevenue(ctx context.Context, request *pborders.FindYearOrder) (*pborders.ApiResponseOrderMonthly, error) {
	s.logger.Info("FindMonthlyRevenue orders called", zap.Int32("year", request.GetYear()))

	year := int(request.GetYear())
	if year <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	res, err := s.orderStats.FindMonthlyOrder(ctx, year)
	if err != nil {
		s.logger.Error("FindMonthlyRevenue orders failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthlyRevenue orders success")

	return &pborders.ApiResponseOrderMonthly{
		Status:  "success",
		Message: "Monthly revenue data retrieved",
		Data:    mapResponsesOrderMonthlyPrices(res),
	}, nil
}

func (s *orderStatsHandleGrpc) FindYearlyRevenue(ctx context.Context, request *pborders.FindYearOrder) (*pborders.ApiResponseOrderYearly, error) {
	s.logger.Info("FindYearlyRevenue orders called", zap.Int32("year", request.GetYear()))

	year := int(request.GetYear())
	if year <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	res, err := s.orderStats.FindYearlyOrder(ctx, year)
	if err != nil {
		s.logger.Error("FindYearlyRevenue orders failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearlyRevenue orders success")

	return &pborders.ApiResponseOrderYearly{
		Status:  "success",
		Message: "Yearly revenue data retrieved",
		Data:    mapResponsesOrderYearlyPrices(res),
	}, nil
}
