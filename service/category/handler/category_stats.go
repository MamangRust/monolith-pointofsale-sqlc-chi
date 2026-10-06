package handler

import (
	"context"
	pbcategories "github.com/MamangRust/monolith-point-of-sale-pb/categories"
	"go.uber.org/zap"

	"github.com/MamangRust/monolith-point-of-sale-category/service"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/category_errors"
)

type categoryStatsHandleGrpc struct {
	pbcategories.UnimplementedCategoryStatsServiceServer
	categoryStats service.CategoryStatsService
	logger        logger.LoggerInterface
}

func NewCategoryStatsHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pbcategories.CategoryStatsServiceServer {
	return &categoryStatsHandleGrpc{
		categoryStats: service.CategoryStats,
		logger:        logger,
	}
}

func (s *categoryStatsHandleGrpc) FindMonthlyTotalPrices(ctx context.Context, req *pbcategories.FindYearMonthTotalPrices) (*pbcategories.ApiResponseCategoryMonthlyTotalPrice, error) {
	s.logger.Info("FindMonthlyTotalPrices categories called", zap.Int32("year", req.GetYear()))

	year := int(req.GetYear())
	month := int(req.GetMonth())

	if year <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidYear
	}
	if month <= 0 || month > 12 {
		return nil, category_errors.ErrGrpcFailedInvalidMonth
	}

	reqService := requests.MonthTotalPrice{
		Year:  year,
		Month: month,
	}

	methods, err := s.categoryStats.FindMonthlyTotalPrice(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthlyTotalPrices categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthlyTotalPrices categories success")

	return &pbcategories.ApiResponseCategoryMonthlyTotalPrice{
		Status:  "success",
		Message: "Monthly sales retrieved successfully",
		Data:    mapResponseCategoryMonthlyTotalPrices(methods),
	}, nil
}

func (s *categoryStatsHandleGrpc) FindYearlyTotalPrices(ctx context.Context, req *pbcategories.FindYearTotalPrices) (*pbcategories.ApiResponseCategoryYearlyTotalPrice, error) {
	s.logger.Info("FindYearlyTotalPrices categories called", zap.Int32("year", req.GetYear()))

	year := int(req.GetYear())
	if year <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidYear
	}

	methods, err := s.categoryStats.FindYearlyTotalPrice(ctx, year)
	if err != nil {
		s.logger.Error("FindYearlyTotalPrices categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearlyTotalPrices categories success")

	return &pbcategories.ApiResponseCategoryYearlyTotalPrice{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapResponseCategoryYearlyTotalPrices(methods),
	}, nil
}

func (s *categoryStatsHandleGrpc) FindMonthPrice(ctx context.Context, req *pbcategories.FindYearCategory) (*pbcategories.ApiResponseCategoryMonthPrice, error) {
	s.logger.Info("FindMonthPrice categories called", zap.Int32("year", req.GetYear()))

	year := int(req.GetYear())
	if year <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidYear
	}

	methods, err := s.categoryStats.FindMonthPrice(ctx, year)
	if err != nil {
		s.logger.Error("FindMonthPrice categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthPrice categories success")

	return &pbcategories.ApiResponseCategoryMonthPrice{
		Status:  "success",
		Message: "Monthly payment methods retrieved successfully",
		Data:    mapResponsesCategoryMonthlyPrices(methods),
	}, nil
}

func (s *categoryStatsHandleGrpc) FindYearPrice(ctx context.Context, req *pbcategories.FindYearCategory) (*pbcategories.ApiResponseCategoryYearPrice, error) {
	s.logger.Info("FindYearPrice categories called", zap.Int32("year", req.GetYear()))

	year := int(req.GetYear())
	if year <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidYear
	}

	methods, err := s.categoryStats.FindYearPrice(ctx, year)
	if err != nil {
		s.logger.Error("FindYearPrice categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearPrice categories success")

	return &pbcategories.ApiResponseCategoryYearPrice{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapResponsesCategoryYearlyPrices(methods),
	}, nil
}
