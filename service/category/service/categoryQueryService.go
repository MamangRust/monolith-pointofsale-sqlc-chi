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

type categoryQueryDeps struct {
	Cache         mencache.CategoryQueryCache
	CategoryQuery repository.CategoryQueryRepository
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

type categoryQueryService struct {
	mencache      mencache.CategoryQueryCache
	categoryQuery repository.CategoryQueryRepository
	logger        logger.LoggerInterface
	observability observability.TraceLoggerObservability
}

func NewCategoryQueryService(params *categoryQueryDeps) CategoryQueryService {
	return &categoryQueryService{
		mencache:      params.Cache,
		categoryQuery: params.CategoryQuery,
		logger:        params.Logger,
		observability: params.Observability,
	}
}

func (s *categoryQueryService) FindAll(ctx context.Context, req *requests.FindAllCategory) ([]*models.CategoryRow, *int, error) {
	const method = "FindAll"
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", req.Search),
	)
	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedCategoriesCache(ctx, req); found {
		logSuccess("Successfully fetched categories from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	category, err := s.categoryQueryRepository().FindAllCategory(ctx, req)
	if err != nil {
		status = "error"
		_, mappedErr := sharederrorhandler.HandleError[[]*models.CategoryRow](
			s.logger,
			err,
			method,
			span,
			zap.Error(err),
		)
		return nil, nil, mappedErr
	}

	var totalCount int
	if len(category) > 0 {
		totalCount = int(category[0].TotalCount)
	}

	s.mencache.SetCachedCategoriesCache(ctx, req, category, &totalCount)
	logSuccess("Successfully fetched categories", zap.Int("totalRecords", totalCount), zap.Int("page", page), zap.Int("pageSize", pageSize))
	return category, &totalCount, nil
}

func (s *categoryQueryService) FindByActive(ctx context.Context, req *requests.FindAllCategory) ([]*models.CategoryActiveRow, *int, error) {
	const method = "FindByActive"
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", req.Search),
	)
	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedCategoryActiveCache(ctx, req); found {
		logSuccess("Successfully fetched active categories from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	category, err := s.categoryQueryRepository().FindByActive(ctx, req)
	if err != nil {
		status = "error"
		_, mappedErr := sharederrorhandler.HandleError[[]*models.CategoryActiveRow](
			s.logger,
			err,
			method,
			span,
			zap.Error(err),
		)
		return nil, nil, mappedErr
	}

	var totalCount int
	if len(category) > 0 {
		totalCount = int(category[0].TotalCount)
	}

	s.mencache.SetCachedCategoryActiveCache(ctx, req, category, &totalCount)
	logSuccess("Successfully fetched active categories", zap.Int("totalRecords", totalCount), zap.Int("page", page), zap.Int("pageSize", pageSize))
	return category, &totalCount, nil
}

func (s *categoryQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllCategory) ([]*models.CategoryTrashedRow, *int, error) {
	const method = "FindByTrashed"
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", req.Search),
	)
	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedCategoryTrashedCache(ctx, req); found {
		logSuccess("Successfully fetched trashed categories from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	categories, err := s.categoryQueryRepository().FindByTrashed(ctx, req)
	if err != nil {
		status = "error"
		_, mappedErr := sharederrorhandler.HandleError[[]*models.CategoryTrashedRow](
			s.logger,
			err,
			method,
			span,
			zap.Error(err),
		)
		return nil, nil, mappedErr
	}

	var totalCount int
	if len(categories) > 0 {
		totalCount = int(categories[0].TotalCount)
	}

	s.mencache.SetCachedCategoryTrashedCache(ctx, req, categories, &totalCount)
	logSuccess("Successfully fetched trashed categories", zap.Int("totalRecords", totalCount), zap.Int("page", page), zap.Int("pageSize", pageSize))
	return categories, &totalCount, nil
}

func (s *categoryQueryService) FindById(ctx context.Context, category_id int) (*models.Category, error) {
	const method = "FindById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("category.id", category_id),
	)
	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedCategoryCache(ctx, category_id); found {
		logSuccess("Successfully fetched category from cache", zap.Int("category.id", category_id))
		return data, nil
	}

	category, err := s.categoryQueryRepository().FindById(ctx, category_id)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*models.Category](
			s.logger,
			err,
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.SetCachedCategoryCache(ctx, category)
	logSuccess("Successfully fetched category", zap.Int("category.id", category_id))
	return category, nil
}

func (s *categoryQueryService) FindByName(ctx context.Context, name string) (*models.Category, error) {
	const method = "FindByName"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.String("category.name", name),
	)
	defer func() {
		end(status)
	}()

	category, err := s.categoryQueryRepository().FindByName(ctx, name)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*models.Category](
			s.logger,
			err,
			method,
			span,
			zap.Error(err),
		)
	}

	logSuccess("Successfully fetched category by name", zap.String("category.name", name))
	return category, nil
}

func (s *categoryQueryService) FindByIds(ctx context.Context, category_ids []int) ([]*models.Category, error) {
	const method = "FindByIds"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("category.count", len(category_ids)),
	)
	defer func() {
		end(status)
	}()

	categories, err := s.categoryQueryRepository().FindByIds(ctx, category_ids)
	if err != nil {
		status = "error"
		_, mappedErr := sharederrorhandler.HandleError[[]*models.Category](
			s.logger,
			err,
			method,
			span,
			zap.Error(err),
		)
		return nil, mappedErr
	}

	logSuccess("Successfully fetched categories by IDs", zap.Int("count", len(categories)))
	return categories, nil
}

func (s *categoryQueryService) categoryQueryRepository() repository.CategoryQueryRepository {
	return s.categoryQuery
}

func (s *categoryQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}
