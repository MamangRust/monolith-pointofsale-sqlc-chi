package service

import (
	"context"

	mencache "github.com/MamangRust/monolith-point-of-sale-cashier/cache"
	"github.com/MamangRust/monolith-point-of-sale-cashier/repository"
	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	sharederrorhandler "github.com/MamangRust/monolith-point-of-sale-shared/errorhandler"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/cashier_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type cashierStatsByMerchantDeps struct {
	Cache         mencache.CashierStatsByMerchantCache
	CashierStats  repository.CashierStatByMerchantRepository
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

type cashierStatsByMerchantQueryService struct {
	mencache      mencache.CashierStatsByMerchantCache
	cashierStats  repository.CashierStatByMerchantRepository
	logger        logger.LoggerInterface
	observability observability.TraceLoggerObservability
}

func NewCashierStatsByMerchantQueryService(params *cashierStatsByMerchantDeps) CashierStatsByMerchant {
	return &cashierStatsByMerchantQueryService{
		mencache:      params.Cache,
		cashierStats:  params.CashierStats,
		logger:        params.Logger,
		observability: params.Observability,
	}
}

func (s *cashierStatsByMerchantQueryService) FindMonthlyTotalSalesByMerchant(ctx context.Context, req *requests.MonthTotalSalesMerchant) ([]*models.CashierMonthlyTotalSalesRow, error) {
	const method = "FindMonthlyTotalSalesByMerchant"
	month := req.Month
	year := req.Year

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.Int("month", month))
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthlyTotalSalesByMerchantCache(ctx, req); found {
		logSuccess("Successfully fetched monthly total sales by merchant from cache", zap.Int("year", year), zap.Int("month", month))
		return data, nil
	}

	res, err := s.cashierStats.GetMonthlyTotalSalesByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*models.CashierMonthlyTotalSalesRow](
			s.logger,
			cashier_errors.ErrFailedFindMonthlyTotalSalesByMerchant,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetMonthlyTotalSalesByMerchantCache(ctx, req, res)
	logSuccess("Successfully fetched monthly total sales by merchant", zap.Int("year", year), zap.Int("month", month))
	return res, nil
}

func (s *cashierStatsByMerchantQueryService) FindYearlyTotalSalesByMerchant(ctx context.Context, req *requests.YearTotalSalesMerchant) ([]*models.CashierYearlyTotalSalesRow, error) {
	const method = "FindYearlyTotalSalesByMerchant"
	year := req.Year
	merchant_id := req.MerchantID

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.Int("merchant_id", merchant_id))
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyTotalSalesByMerchantCache(ctx, req); found {
		logSuccess("Successfully fetched yearly total sales by merchant from cache", zap.Int("year", year), zap.Int("merchant.id", merchant_id))
		return data, nil
	}

	res, err := s.cashierStats.GetYearlyTotalSalesByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*models.CashierYearlyTotalSalesRow](
			s.logger,
			cashier_errors.ErrFailedFindYearlyTotalSalesByMerchant,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetYearlyTotalSalesByMerchantCache(ctx, req, res)
	logSuccess("Successfully fetched yearly total sales by merchant", zap.Int("year", year), zap.Int("merchant.id", merchant_id))
	return res, nil
}

func (s *cashierStatsByMerchantQueryService) FindMonthlyCashierByMerchant(ctx context.Context, req *requests.MonthCashierMerchant) ([]*models.CashierMonthlySaleByMerchantRow, error) {
	const method = "FindMonthlyCashierByMerchant"
	year := req.Year
	merchant_id := req.MerchantID

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.Int("merchant.id", merchant_id))
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetMonthlyCashierByMerchantCache(ctx, req); found {
		logSuccess("Successfully fetched monthly cashier sales by merchant from cache", zap.Int("year", year), zap.Int("merchant.id", merchant_id))
		return data, nil
	}

	res, err := s.cashierStats.GetMonthlyCashierByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*models.CashierMonthlySaleByMerchantRow](
			s.logger,
			cashier_errors.ErrFailedFindMonthlyCashierByMerchant,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetMonthlyCashierByMerchantCache(ctx, req, res)
	logSuccess("Successfully fetched monthly cashier sales by merchant", zap.Int("year", year), zap.Int("merchant.id", merchant_id))
	return res, nil
}

func (s *cashierStatsByMerchantQueryService) FindYearlyCashierByMerchant(ctx context.Context, req *requests.YearCashierMerchant) ([]*models.CashierYearlySaleRow, error) {
	const method = "FindYearlyCashierByMerchant"
	year := req.Year
	merchant_id := req.MerchantID

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("year", year), attribute.Int("merchant.id", merchant_id))
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetYearlyCashierByMerchantCache(ctx, req); found {
		logSuccess("Successfully fetched yearly cashier sales by merchant from cache", zap.Int("year", year), zap.Int("merchant.id", merchant_id))
		return data, nil
	}

	res, err := s.cashierStats.GetYearlyCashierByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*models.CashierYearlySaleRow](
			s.logger,
			cashier_errors.ErrFailedFindYearlyCashierByMerchant,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetYearlyCashierByMerchantCache(ctx, req, res)
	logSuccess("Successfully fetched yearly cashier sales by merchant", zap.Int("year", year), zap.Int("merchant.id", merchant_id))
	return res, nil
}
