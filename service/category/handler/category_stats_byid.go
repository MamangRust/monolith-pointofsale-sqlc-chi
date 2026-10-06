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

type categoryStatsByIdHandleGrpc struct {
	pbcategories.UnimplementedCategoryStatsByIdServiceServer
	categoryStatsById service.CategoryStatsByIdService
	logger            logger.LoggerInterface
}

func NewCategoryStatsByIdHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pbcategories.CategoryStatsByIdServiceServer {
	return &categoryStatsByIdHandleGrpc{
		categoryStatsById: service.CategoryStatsById,
		logger:            logger,
	}
}

func (s *categoryStatsByIdHandleGrpc) FindMonthlyTotalPricesById(ctx context.Context, req *pbcategories.FindYearMonthTotalPriceById) (*pbcategories.ApiResponseCategoryMonthlyTotalPrice, error) {
	s.logger.Info("FindMonthlyTotalPricesById categories called", zap.Int32("id", req.GetCategoryId()))

	year := int(req.GetYear())
	month := int(req.GetMonth())
	id := int(req.GetCategoryId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidYear
	}
	if month <= 0 || month > 12 {
		return nil, category_errors.ErrGrpcFailedInvalidMonth
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidId
	}

	reqService := requests.MonthTotalPriceCategory{
		Year:       year,
		Month:      month,
		CategoryID: id,
	}

	methods, err := s.categoryStatsById.FindMonthlyTotalPriceById(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthlyTotalPricesById categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthlyTotalPricesById categories success")

	return &pbcategories.ApiResponseCategoryMonthlyTotalPrice{
		Status:  "success",
		Message: "Monthly sales retrieved successfully",
		Data:    mapResponseCategoryMonthlyTotalPricesById(methods),
	}, nil
}

func (s *categoryStatsByIdHandleGrpc) FindYearlyTotalPricesById(ctx context.Context, req *pbcategories.FindYearTotalPriceById) (*pbcategories.ApiResponseCategoryYearlyTotalPrice, error) {
	s.logger.Info("FindYearlyTotalPricesById categories called", zap.Int32("id", req.GetCategoryId()))

	year := int(req.GetYear())
	id := int(req.GetCategoryId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidYear
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidId
	}

	reqService := requests.YearTotalPriceCategory{
		Year:       year,
		CategoryID: id,
	}

	methods, err := s.categoryStatsById.FindYearlyTotalPriceById(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindYearlyTotalPricesById categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearlyTotalPricesById categories success")

	return &pbcategories.ApiResponseCategoryYearlyTotalPrice{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapResponseCategoryYearlyTotalPricesById(methods),
	}, nil
}

func (s *categoryStatsByIdHandleGrpc) FindMonthPriceById(ctx context.Context, req *pbcategories.FindYearCategoryById) (*pbcategories.ApiResponseCategoryMonthPrice, error) {
	s.logger.Info("FindMonthPriceById categories called", zap.Int32("id", req.GetCategoryId()))

	year := int(req.GetYear())
	id := int(req.GetCategoryId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidYear
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidId
	}

	reqService := requests.MonthPriceId{
		Year:       year,
		CategoryID: id,
	}

	methods, err := s.categoryStatsById.FindMonthPriceById(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthPriceById categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthPriceById categories success")

	return &pbcategories.ApiResponseCategoryMonthPrice{
		Status:  "success",
		Message: "Merchant monthly payment methods retrieved successfully",
		Data:    mapResponsesCategoryMonthlyPricesById(methods),
	}, nil
}

func (s *categoryStatsByIdHandleGrpc) FindYearPriceById(ctx context.Context, req *pbcategories.FindYearCategoryById) (*pbcategories.ApiResponseCategoryYearPrice, error) {
	s.logger.Info("FindYearPriceById categories called", zap.Int32("id", req.GetCategoryId()))

	year := int(req.GetYear())
	id := int(req.GetCategoryId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidYear
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidId
	}

	reqService := requests.YearPriceId{
		Year:       year,
		CategoryID: id,
	}

	methods, err := s.categoryStatsById.FindYearPriceById(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindYearPriceById categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearPriceById categories success")

	return &pbcategories.ApiResponseCategoryYearPrice{
		Status:  "success",
		Message: "Merchant yearly payment methods retrieved successfully",
		Data:    mapResponsesCategoryYearlyPricesById(methods),
	}, nil
}
