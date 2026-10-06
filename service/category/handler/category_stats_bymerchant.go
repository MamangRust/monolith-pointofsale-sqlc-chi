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

type categoryStatsByMerchantHandleGrpc struct {
	pbcategories.UnimplementedCategoryStatsByMerchantServiceServer
	categoryStatsByMerchant service.CategoryStatsByMerchantQueryService
	logger                  logger.LoggerInterface
}

func NewCategoryStatsByMerchantHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pbcategories.CategoryStatsByMerchantServiceServer {
	return &categoryStatsByMerchantHandleGrpc{
		categoryStatsByMerchant: service.CategoryStatsByMerchant,
		logger:                  logger,
	}
}

func (s *categoryStatsByMerchantHandleGrpc) FindMonthlyTotalPricesByMerchant(ctx context.Context, req *pbcategories.FindYearMonthTotalPriceByMerchant) (*pbcategories.ApiResponseCategoryMonthlyTotalPrice, error) {
	s.logger.Info("FindMonthlyTotalPricesByMerchant categories called", zap.Int32("merchantId", req.GetMerchantId()))

	year := int(req.GetYear())
	month := int(req.GetMonth())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidYear
	}
	if month <= 0 || month > 12 {
		return nil, category_errors.ErrGrpcFailedInvalidMonth
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidMerchantId
	}

	reqService := requests.MonthTotalPriceMerchant{
		Year:       year,
		Month:      month,
		MerchantID: id,
	}

	methods, err := s.categoryStatsByMerchant.FindMonthlyTotalPriceByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthlyTotalPricesByMerchant categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthlyTotalPricesByMerchant categories success")

	return &pbcategories.ApiResponseCategoryMonthlyTotalPrice{
		Status:  "success",
		Message: "Monthly sales retrieved successfully",
		Data:    mapResponseCategoryMonthlyTotalPricesByMerchant(methods),
	}, nil
}

func (s *categoryStatsByMerchantHandleGrpc) FindYearlyTotalPricesByMerchant(ctx context.Context, req *pbcategories.FindYearTotalPriceByMerchant) (*pbcategories.ApiResponseCategoryYearlyTotalPrice, error) {
	s.logger.Info("FindYearlyTotalPricesByMerchant categories called", zap.Int32("merchantId", req.GetMerchantId()))

	year := int(req.GetYear())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidYear
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidMerchantId
	}

	reqService := requests.YearTotalPriceMerchant{
		Year:       year,
		MerchantID: id,
	}

	methods, err := s.categoryStatsByMerchant.FindYearlyTotalPriceByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindYearlyTotalPricesByMerchant categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearlyTotalPricesByMerchant categories success")

	return &pbcategories.ApiResponseCategoryYearlyTotalPrice{
		Status:  "success",
		Message: "Yearly payment methods retrieved successfully",
		Data:    mapResponseCategoryYearlyTotalPricesByMerchant(methods),
	}, nil
}

func (s *categoryStatsByMerchantHandleGrpc) FindMonthPriceByMerchant(ctx context.Context, req *pbcategories.FindYearCategoryByMerchant) (*pbcategories.ApiResponseCategoryMonthPrice, error) {
	s.logger.Info("FindMonthPriceByMerchant categories called", zap.Int32("merchantId", req.GetMerchantId()))

	year := int(req.GetYear())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidYear
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidMerchantId
	}

	reqService := requests.MonthPriceMerchant{
		Year:       year,
		MerchantID: id,
	}

	methods, err := s.categoryStatsByMerchant.FindMonthPriceByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindMonthPriceByMerchant categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindMonthPriceByMerchant categories success")

	return &pbcategories.ApiResponseCategoryMonthPrice{
		Status:  "success",
		Message: "Merchant monthly payment methods retrieved successfully",
		Data:    mapResponsesCategoryMonthlyPricesByMerchant(methods),
	}, nil
}

func (s *categoryStatsByMerchantHandleGrpc) FindYearPriceByMerchant(ctx context.Context, req *pbcategories.FindYearCategoryByMerchant) (*pbcategories.ApiResponseCategoryYearPrice, error) {
	s.logger.Info("FindYearPriceByMerchant categories called", zap.Int32("merchantId", req.GetMerchantId()))

	year := int(req.GetYear())
	id := int(req.GetMerchantId())

	if year <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidYear
	}
	if id <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidMerchantId
	}

	reqService := requests.YearPriceMerchant{
		Year:       year,
		MerchantID: id,
	}

	methods, err := s.categoryStatsByMerchant.FindYearPriceByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindYearPriceByMerchant categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindYearPriceByMerchant categories success")

	return &pbcategories.ApiResponseCategoryYearPrice{
		Status:  "success",
		Message: "Merchant yearly payment methods retrieved successfully",
		Data:    mapResponsesCategoryYearlyPricesByMerchant(methods),
	}, nil
}
