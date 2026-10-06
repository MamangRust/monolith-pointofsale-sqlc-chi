package service

import (
	"context"

	mencache "github.com/MamangRust/monolith-point-of-sale-category/cache"
	"github.com/MamangRust/monolith-point-of-sale-category/repository"
	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	sharederrorhandler "github.com/MamangRust/monolith-point-of-sale-shared/errorhandler"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type categoryStatsByMerchantDeps struct {
	Cache                             mencache.CategoryStatsByMerchantCache
	CategoryStatsByMerchantRepository repository.CategoryStatsByMerchantRepository
	Logger                            logger.LoggerInterface
	Observability                     observability.TraceLoggerObservability
}

type categoryStatsByMerchantQueryService struct {
	mencache                          mencache.CategoryStatsByMerchantCache
	categoryStatsByMerchantRepository repository.CategoryStatsByMerchantRepository
	logger                            logger.LoggerInterface
	observability                     observability.TraceLoggerObservability
}

func NewCategoryStatsByMerchantQueryService(params *categoryStatsByMerchantDeps) CategoryStatsByMerchantQueryService {
	return &categoryStatsByMerchantQueryService{
		mencache:                          params.Cache,
		categoryStatsByMerchantRepository: params.CategoryStatsByMerchantRepository,
		logger:                            params.Logger,
		observability:                     params.Observability,
	}
}

func (s *categoryStatsByMerchantQueryService) FindMonthlyTotalPriceByMerchant(ctx context.Context, req *requests.MonthTotalPriceMerchant) ([]*models.CategoryMonthlyTotalPriceByMerchantRow, error) {
	const method = "FindMonthlyTotalPriceByMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("month", req.Month),
		attribute.Int("merchant.id", req.MerchantID),
	)
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMonthTotalPriceByMerchantCache(ctx, req); found {
		logSuccess("Successfully fetched monthly total price by merchant from cache", zap.Int("year", req.Year), zap.Int("month", req.Month))
		return data, nil
	}

	res, err := s.categoryStatsByMerchantRepository.GetMonthlyTotalPriceByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*models.CategoryMonthlyTotalPriceByMerchantRow](
			s.logger,
			err,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedMonthTotalPriceByMerchantCache(ctx, req, res)
	logSuccess("Successfully fetched monthly total price by merchant", zap.Int("year", req.Year), zap.Int("month", req.Month))
	return res, nil
}

func (s *categoryStatsByMerchantQueryService) FindYearlyTotalPriceByMerchant(ctx context.Context, req *requests.YearTotalPriceMerchant) ([]*models.CategoryYearlyTotalPriceByMerchantRow, error) {
	const method = "FindYearlyTotalPriceByMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("merchant.id", req.MerchantID),
	)
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedYearTotalPriceByMerchantCache(ctx, req); found {
		logSuccess("Successfully fetched yearly total price by merchant from cache", zap.Int("year", req.Year), zap.Int("merchant.id", req.MerchantID))
		return data, nil
	}

	res, err := s.categoryStatsByMerchantRepository.GetYearlyTotalPricesByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*models.CategoryYearlyTotalPriceByMerchantRow](
			s.logger,
			err,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedYearTotalPriceByMerchantCache(ctx, req, res)
	logSuccess("Successfully fetched yearly total price by merchant", zap.Int("year", req.Year), zap.Int("merchant.id", req.MerchantID))
	return res, nil
}

func (s *categoryStatsByMerchantQueryService) FindMonthPriceByMerchant(ctx context.Context, req *requests.MonthPriceMerchant) ([]*models.CategoryMonthlyByMerchantRow, error) {
	const method = "FindMonthPriceByMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("merchant.id", req.MerchantID),
	)
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMonthPriceByMerchantCache(ctx, req); found {
		logSuccess("Successfully fetched monthly category prices by merchant from cache", zap.Int("year", req.Year), zap.Int("merchant.id", req.MerchantID))
		return data, nil
	}

	res, err := s.categoryStatsByMerchantRepository.GetMonthPriceByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*models.CategoryMonthlyByMerchantRow](
			s.logger,
			err,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedMonthPriceByMerchantCache(ctx, req, res)
	logSuccess("Successfully fetched monthly category prices by merchant", zap.Int("year", req.Year), zap.Int("merchant.id", req.MerchantID))
	return res, nil
}

func (s *categoryStatsByMerchantQueryService) FindYearPriceByMerchant(ctx context.Context, req *requests.YearPriceMerchant) ([]*models.CategoryYearlyByMerchantRow, error) {
	const method = "FindYearPriceByMerchant"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", req.Year),
		attribute.Int("merchant.id", req.MerchantID),
	)
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedYearPriceByMerchantCache(ctx, req); found {
		logSuccess("Successfully fetched yearly category prices by merchant from cache", zap.Int("year", req.Year), zap.Int("merchant.id", req.MerchantID))
		return data, nil
	}

	res, err := s.categoryStatsByMerchantRepository.GetYearPriceByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*models.CategoryYearlyByMerchantRow](
			s.logger,
			err,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedYearPriceByMerchantCache(ctx, req, res)
	logSuccess("Successfully fetched yearly category prices by merchant", zap.Int("year", req.Year), zap.Int("merchant.id", req.MerchantID))
	return res, nil
}
