package service

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	sharederrorhandler "github.com/MamangRust/monolith-point-of-sale-shared/errorhandler"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/transaction_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	mencache "github.com/MamangRust/monolith-point-of-sale-transacton/cache"
	"github.com/MamangRust/monolith-point-of-sale-transacton/repository"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type transactionStatsByMerchantQueryService struct {
	mencache                             mencache.TransactionStatsByMerchantCache
	transactionStatsByMerchantRepository repository.TransactionStatsByMerchantRepository
	logger                               logger.LoggerInterface
	observability                        observability.TraceLoggerObservability
}

func NewTransactionStatsByMerchantQueryService(
	mencache mencache.TransactionStatsByMerchantCache,
	transactionStatsByMerchantRepository repository.TransactionStatsByMerchantRepository,
	logger logger.LoggerInterface,
	obs observability.TraceLoggerObservability,
) *transactionStatsByMerchantQueryService {
	return &transactionStatsByMerchantQueryService{
		mencache:                             mencache,
		transactionStatsByMerchantRepository: transactionStatsByMerchantRepository,
		logger:                               logger,
		observability:                        obs,
	}
}

func (s *transactionStatsByMerchantQueryService) FindMonthlyAmountSuccessByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*models.TransactionMonthlyAmountSuccessByMerchantRow, error) {
	const method = "FindMonthlyAmountSuccessByMerchant"

	year := req.Year
	month := req.Month
	merchantID := req.MerchantID

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year), attribute.Int("month", month), attribute.Int("merchant.id", merchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMonthAmountSuccessByMerchantCached(ctx, req); found {
		logSuccess("Successfully fetched monthly successful transactions by merchant from cache", zap.Int("year", year), zap.Int("month", month), zap.Int("merchant.id", merchantID))
		return data, nil
	}

	res, err := s.transactionStatsByMerchantRepository.GetMonthlyAmountSuccessByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*models.TransactionMonthlyAmountSuccessByMerchantRow](
			s.logger,
			transaction_errors.ErrFailedFindMonthlyAmountSuccessByMerchant.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedMonthAmountSuccessByMerchantCached(ctx, req, res)

	logSuccess("Successfully fetched monthly successful transactions by merchant", zap.Int("year", year), zap.Int("month", month), zap.Int("merchant.id", merchantID))

	return res, nil
}

func (s *transactionStatsByMerchantQueryService) FindYearlyAmountSuccessByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*models.TransactionYearlyAmountSuccessByMerchantRow, error) {
	const method = "FindYearlyAmountSuccessByMerchant"

	year := req.Year
	merchantID := req.MerchantID

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year), attribute.Int("merchant.id", merchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedYearAmountSuccessByMerchantCached(ctx, req); found {
		logSuccess("Successfully fetched yearly successful transactions by merchant from cache", zap.Int("year", year), zap.Int("merchant.id", merchantID))
		return data, nil
	}

	res, err := s.transactionStatsByMerchantRepository.GetYearlyAmountSuccessByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*models.TransactionYearlyAmountSuccessByMerchantRow](
			s.logger,
			transaction_errors.ErrFailedFindYearlyAmountSuccessByMerchant.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedYearAmountSuccessByMerchantCached(ctx, req, res)

	logSuccess("Successfully fetched yearly successful transactions by merchant", zap.Int("year", year), zap.Int("merchant.id", merchantID))

	return res, nil
}

func (s *transactionStatsByMerchantQueryService) FindMonthlyAmountFailedByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*models.TransactionMonthlyAmountFailedByMerchantRow, error) {
	const method = "FindMonthlyAmountFailedByMerchant"

	year := req.Year
	month := req.Month
	merchantID := req.MerchantID

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year), attribute.Int("month", month), attribute.Int("merchant.id", merchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMonthAmountFailedByMerchantCached(ctx, req); found {
		logSuccess("Successfully fetched monthly failed transactions by merchant from cache", zap.Int("year", year), zap.Int("month", month), zap.Int("merchant.id", merchantID))
		return data, nil
	}

	res, err := s.transactionStatsByMerchantRepository.GetMonthlyAmountFailedByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*models.TransactionMonthlyAmountFailedByMerchantRow](
			s.logger,
			transaction_errors.ErrFailedFindMonthlyAmountFailedByMerchant.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedMonthAmountFailedByMerchantCached(ctx, req, res)

	logSuccess("Successfully fetched monthly failed transactions by merchant", zap.Int("year", year), zap.Int("month", month), zap.Int("merchant.id", merchantID))

	return res, nil
}

func (s *transactionStatsByMerchantQueryService) FindYearlyAmountFailedByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*models.TransactionYearlyAmountFailedByMerchantRow, error) {
	const method = "FindYearlyAmountFailedByMerchant"

	year := req.Year
	merchantID := req.MerchantID

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year), attribute.Int("merchant.id", merchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedYearAmountFailedByMerchantCached(ctx, req); found {
		logSuccess("Successfully fetched yearly failed transactions by merchant from cache", zap.Int("year", year), zap.Int("merchant.id", merchantID))
		return data, nil
	}

	res, err := s.transactionStatsByMerchantRepository.GetYearlyAmountFailedByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*models.TransactionYearlyAmountFailedByMerchantRow](
			s.logger,
			transaction_errors.ErrFailedFindYearlyAmountFailedByMerchant.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedYearAmountFailedByMerchantCached(ctx, req, res)

	logSuccess("Successfully fetched yearly failed transactions by merchant", zap.Int("year", year), zap.Int("merchant.id", merchantID))

	return res, nil
}

func (s *transactionStatsByMerchantQueryService) FindMonthlyMethodByMerchantSuccess(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*models.TransactionMonthlyMethodByMerchantSuccessRow, error) {
	const method = "FindMonthlyMethodByMerchantSuccess"

	year := req.Year
	month := req.Month
	merchantID := req.MerchantID

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year), attribute.Int("month", month), attribute.Int("merchant.id", merchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMonthMethodSuccessByMerchantCached(ctx, req); found {
		logSuccess("Successfully found monthly successful transaction methods by merchant from cache", zap.Int("year", year), zap.Int("month", month), zap.Int("merchant.id", merchantID))
		return data, nil
	}

	res, err := s.transactionStatsByMerchantRepository.GetMonthlyTransactionMethodByMerchantSuccess(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*models.TransactionMonthlyMethodByMerchantSuccessRow](
			s.logger,
			transaction_errors.ErrFailedFindMonthlyMethodByMerchant.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedMonthMethodSuccessByMerchantCached(ctx, req, res)

	logSuccess("Successfully found monthly successful transaction methods by merchant", zap.Int("year", year), zap.Int("month", month), zap.Int("merchant.id", merchantID))

	return res, nil
}

func (s *transactionStatsByMerchantQueryService) FindYearlyMethodByMerchantSuccess(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*models.TransactionYearlyMethodByMerchantSuccessRow, error) {
	const method = "FindYearlyMethodByMerchantSuccess"

	year := req.Year
	merchantID := req.MerchantID

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year), attribute.Int("merchant.id", merchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedYearMethodSuccessByMerchantCached(ctx, req); found {
		logSuccess("Successfully found yearly successful transaction methods by merchant from cache", zap.Int("year", year), zap.Int("merchant.id", merchantID))
		return data, nil
	}

	res, err := s.transactionStatsByMerchantRepository.GetYearlyTransactionMethodByMerchantSuccess(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*models.TransactionYearlyMethodByMerchantSuccessRow](
			s.logger,
			transaction_errors.ErrFailedFindYearlyMethodByMerchant.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedYearMethodSuccessByMerchantCached(ctx, req, res)

	logSuccess("Successfully found yearly successful transaction methods by merchant", zap.Int("year", year), zap.Int("merchant.id", merchantID))

	return res, nil
}

func (s *transactionStatsByMerchantQueryService) FindMonthlyMethodByMerchantFailed(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*models.TransactionMonthlyMethodByMerchantFailedRow, error) {
	const method = "FindMonthlyMethodByMerchantFailed"

	year := req.Year
	month := req.Month
	merchantID := req.MerchantID

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year), attribute.Int("month", month), attribute.Int("merchant.id", merchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedMonthMethodFailedByMerchantCached(ctx, req); found {
		logSuccess("Successfully found monthly failed transaction methods by merchant from cache", zap.Int("year", year), zap.Int("month", month), zap.Int("merchant.id", merchantID))
		return data, nil
	}

	res, err := s.transactionStatsByMerchantRepository.GetMonthlyTransactionMethodByMerchantFailed(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*models.TransactionMonthlyMethodByMerchantFailedRow](
			s.logger,
			transaction_errors.ErrFailedFindMonthlyMethodByMerchant.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedMonthMethodFailedByMerchantCached(ctx, req, res)

	logSuccess("Successfully found monthly failed transaction methods by merchant", zap.Int("year", year), zap.Int("month", month), zap.Int("merchant.id", merchantID))

	return res, nil
}

func (s *transactionStatsByMerchantQueryService) FindYearlyMethodByMerchantFailed(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*models.TransactionYearlyMethodByMerchantFailedRow, error) {
	const method = "FindYearlyMethodByMerchantFailed"

	year := req.Year
	merchantID := req.MerchantID

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("year", year), attribute.Int("merchant.id", merchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedYearMethodFailedByMerchantCached(ctx, req); found {
		logSuccess("Successfully found yearly failed transaction methods by merchant from cache", zap.Int("year", year), zap.Int("merchant.id", merchantID))
		return data, nil
	}

	res, err := s.transactionStatsByMerchantRepository.GetYearlyTransactionMethodByMerchantFailed(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[[]*models.TransactionYearlyMethodByMerchantFailedRow](
			s.logger,
			transaction_errors.ErrFailedFindYearlyMethodByMerchant.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedYearMethodFailedByMerchantCached(ctx, req, res)

	logSuccess("Successfully found yearly failed transaction methods by merchant", zap.Int("year", year), zap.Int("merchant.id", merchantID))

	return res, nil
}
