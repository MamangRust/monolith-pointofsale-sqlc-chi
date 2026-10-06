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

type orderStatsByMerchantHandleGrpc struct {
	pborders.UnimplementedOrderStatsByMerchantServiceServer
	orderStatsByMerchant service.OrderStatByMerchantQueryService
	logger               logger.LoggerInterface
}

func NewOrderStatsByMerchantHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pborders.OrderStatsByMerchantServiceServer {
	return &orderStatsByMerchantHandleGrpc{
		orderStatsByMerchant: service.OrderStatsByMerchant,
		logger:               logger,
	}
}

func (s *orderStatsByMerchantHandleGrpc) FindMonthlyTotalRevenueByMerchant(ctx context.Context, req *pborders.FindYearMonthTotalRevenueByMerchant) (*pborders.ApiResponseOrderMonthlyTotalRevenue, error) {
	s.logger.Info("FindMonthlyTotalRevenueByMerchant orders called", zap.Int32("merchantId", req.GetMerchantId()))

	year := int(req.GetYear())
	month := int(req.GetMonth())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, order_errors.ErrGrpcInvalidYear
	}
	if month <= 0 || month >= 12 {
		return nil, order_errors.ErrGrpcInvalidMonth
	}
	if id <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidMerchantId
	}

	reqService := requests.MonthTotalRevenueMerchant{
		Year:       year,
		Month:      month,
		MerchantID: id,
	}

	methods, err := s.orderStatsByMerchant.FindMonthlyTotalRevenueByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthlyTotalRevenueByMerchant orders failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthlyTotalRevenueByMerchant orders success")

	return &pborders.ApiResponseOrderMonthlyTotalRevenue{
		Status:  "success",
		Message: "Monthly sales retrieved successfully",
		Data:    mapResponseOrderMonthlyTotalRevenuesByMerchant(methods),
	}, nil
}

func (s *orderStatsByMerchantHandleGrpc) FindYearlyTotalRevenueByMerchant(ctx context.Context, req *pborders.FindYearTotalRevenueByMerchant) (*pborders.ApiResponseOrderYearlyTotalRevenue, error) {
	s.logger.Info("FindYearlyTotalRevenueByMerchant orders called", zap.Int32("merchantId", req.GetMerchantId()))

	year := int(req.GetYear())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, order_errors.ErrGrpcInvalidYear
	}
	if id <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidMerchantId
	}

	reqService := requests.YearTotalRevenueMerchant{
		Year:       year,
		MerchantID: id,
	}

	methods, err := s.orderStatsByMerchant.FindYearlyTotalRevenueByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindYearlyTotalRevenueByMerchant orders failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearlyTotalRevenueByMerchant orders success")

	return &pborders.ApiResponseOrderYearlyTotalRevenue{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapResponseOrderYearlyTotalRevenuesByMerchant(methods),
	}, nil
}

func (s *orderStatsByMerchantHandleGrpc) FindMonthlyRevenueByMerchant(ctx context.Context, request *pborders.FindYearOrderByMerchant) (*pborders.ApiResponseOrderMonthly, error) {
	s.logger.Info("FindMonthlyRevenueByMerchant orders called", zap.Int32("merchantId", request.GetMerchantId()))

	year := int(request.GetYear())
	id := int(request.GetMerchantId())

	if year <= 0 {
		return nil, order_errors.ErrGrpcInvalidYear
	}
	if id <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidMerchantId
	}

	reqService := requests.MonthOrderMerchant{
		Year:       year,
		MerchantID: id,
	}

	res, err := s.orderStatsByMerchant.FindMonthlyOrderByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthlyRevenueByMerchant orders failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthlyRevenueByMerchant orders success")

	return &pborders.ApiResponseOrderMonthly{
		Status:  "success",
		Message: "Monthly revenue by merchant data retrieved",
		Data:    mapResponsesOrderMonthlyPricesByMerchant(res),
	}, nil
}

func (s *orderStatsByMerchantHandleGrpc) FindYearlyRevenueByMerchant(ctx context.Context, request *pborders.FindYearOrderByMerchant) (*pborders.ApiResponseOrderYearly, error) {
	s.logger.Info("FindYearlyRevenueByMerchant orders called", zap.Int32("merchantId", request.GetMerchantId()))

	year := int(request.GetYear())
	id := int(request.GetMerchantId())

	if year <= 0 {
		return nil, order_errors.ErrGrpcInvalidYear
	}
	if id <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	reqService := requests.YearOrderMerchant{
		Year:       year,
		MerchantID: id,
	}

	res, err := s.orderStatsByMerchant.FindYearlyOrderByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindYearlyRevenueByMerchant orders failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearlyRevenueByMerchant orders success")

	return &pborders.ApiResponseOrderYearly{
		Status:  "success",
		Message: "Yearly revenue by merchant data retrieved",
		Data:    mapResponsesOrderYearlyPricesByMerchant(res),
	}, nil
}
