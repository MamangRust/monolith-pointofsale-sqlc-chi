package service

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	categoryadapter "github.com/MamangRust/monolith-point-of-sale-pkg/adapter/category"
	mencache "github.com/MamangRust/monolith-point-of-sale-product/cache"
	"github.com/MamangRust/monolith-point-of-sale-product/repository"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/product_errors"
	sharederrorhandler "github.com/MamangRust/monolith-point-of-sale-shared/errorhandler"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type productQueryService struct {
	mencache               mencache.ProductQueryCache
	productQueryRepository repository.ProductQueryRepository
	categoryQuery          categoryadapter.Repository
	logger                 logger.LoggerInterface
	observability          observability.TraceLoggerObservability
}

func NewProductQueryService(
	mencache mencache.ProductQueryCache,
	productQueryRepository repository.ProductQueryRepository,
	categoryQuery categoryadapter.Repository,
	logger logger.LoggerInterface,
	obs observability.TraceLoggerObservability,
) *productQueryService {
	return &productQueryService{
		mencache:               mencache,
		productQueryRepository: productQueryRepository,
		categoryQuery:          categoryQuery,
		logger:                 logger,
		observability:          obs,
	}
}

func (s *productQueryService) FindAll(ctx context.Context, req *requests.FindAllProducts) ([]*models.ProductRow, *int, error) {
	const method = "FindAll"
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedProducts(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))
		return data, total, nil
	}

	res, err := s.productQueryRepository.FindAllProducts(ctx, req)
	if err != nil {
		status = "error"
		_, err = sharederrorhandler.HandleError[any](s.logger, err, method, span, zap.Error(err))
		return nil, nil, err
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	}

	s.mencache.SetCachedProducts(ctx, req, res, &totalCount)

	logSuccess("Successfully fetched all products", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return res, &totalCount, nil
}

func (s *productQueryService) FindByMerchant(ctx context.Context, req *requests.ProductByMerchantRequest) ([]*models.ProductMerchantRow, *int, error) {
	const method = "FindByMerchant"
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search
	merchantID := req.MerchantID

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search), attribute.Int("merchant.id", merchantID))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedProductsByMerchant(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search), zap.Int("merchant.id", merchantID))
		return data, total, nil
	}

	res, err := s.productQueryRepository.FindByMerchant(ctx, req)
	if err != nil {
		status = "error"
		_, err = sharederrorhandler.HandleError[any](s.logger, err, method, span, zap.Error(err))
		return nil, nil, err
	}

	if err := s.fillMerchantCategoryNames(ctx, res); err != nil {
		status = "error"
		_, err = sharederrorhandler.HandleError[any](s.logger, err, method, span, zap.Error(err))
		return nil, nil, err
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	}

	s.mencache.SetCachedProductsByMerchant(ctx, req, res, &totalCount)

	logSuccess("Successfully fetched all products by merchant", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search), zap.Int("merchant.id", merchantID))

	return res, &totalCount, nil
}

func (s *productQueryService) FindByCategory(ctx context.Context, req *requests.ProductByCategoryRequest) ([]*models.ProductCategoryNameRow, *int, error) {
	const method = "FindByCategory"
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search
	categoryName := req.CategoryName

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search), attribute.String("category.name", categoryName))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedProductsByCategory(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search), zap.String("category.name", categoryName))
		return data, total, nil
	}

	// Resolve the category name to an id through the Category adapter; the
	// product repository no longer reads the categories table.
	category, err := s.categoryQuery.FindByName(ctx, req.CategoryName)
	if err != nil {
		status = "error"
		_, err = sharederrorhandler.HandleError[any](s.logger, product_errors.ErrFindByCategory, method, span, zap.Error(err))
		return nil, nil, err
	}

	res, err := s.productQueryRepository.FindByCategoryID(ctx, category.CategoryID, req)
	if err != nil {
		status = "error"
		_, err = sharederrorhandler.HandleError[any](s.logger, err, method, span, zap.Error(err))
		return nil, nil, err
	}

	if err := s.fillCategoryNameRows(ctx, res); err != nil {
		status = "error"
		_, err = sharederrorhandler.HandleError[any](s.logger, err, method, span, zap.Error(err))
		return nil, nil, err
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	}

	s.mencache.SetCachedProductsByCategory(ctx, req, res, &totalCount)

	logSuccess("Successfully fetched all products by category", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search), zap.String("category.name", categoryName))

	return res, &totalCount, nil
}

func (s *productQueryService) FindById(ctx context.Context, productID int) (*models.Product, error) {
	const method = "FindById"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.Int("product.id", productID))

	defer func() {
		end(status)
	}()

	if data, found := s.mencache.GetCachedProduct(ctx, productID); found {
		logSuccess("Data found in cache", zap.Int("product.id", productID))
		return data, nil
	}

	product, err := s.productQueryRepository.FindById(ctx, productID)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*models.Product](s.logger, err, method, span, zap.Error(err))
	}

	s.mencache.SetCachedProduct(ctx, product)

	logSuccess("Successfully fetched product by id", zap.Int("product.id", productID))

	return product, nil
}

func (s *productQueryService) FindByActive(ctx context.Context, req *requests.FindAllProducts) ([]*models.ProductActiveRow, *int, error) {
	const method = "FindByActive"
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedProductActive(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))
		return data, total, nil
	}

	res, err := s.productQueryRepository.FindByActive(ctx, req)
	if err != nil {
		status = "error"
		_, err = sharederrorhandler.HandleError[any](s.logger, err, method, span, zap.Error(err))
		return nil, nil, err
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	}

	s.mencache.SetCachedProductActive(ctx, req, res, &totalCount)

	logSuccess("Successfully fetched all products", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return res, &totalCount, nil
}

func (s *productQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllProducts) ([]*models.ProductTrashedRow, *int, error) {
	const method = "FindByTrashed"
	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page), attribute.Int("pageSize", pageSize), attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.mencache.GetCachedProductTrashed(ctx, req); found {
		logSuccess("Data found in cache", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))
		return data, total, nil
	}

	res, err := s.productQueryRepository.FindByTrashed(ctx, req)
	if err != nil {
		status = "error"
		_, err = sharederrorhandler.HandleError[any](s.logger, err, method, span, zap.Error(err))
		return nil, nil, err
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	}

	s.mencache.SetCachedProductTrashed(ctx, req, res, &totalCount)

	logSuccess("Successfully fetched all products", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.String("search", search))

	return res, &totalCount, nil
}

// fillMerchantCategoryNames resolves the category names of one page of
// products with a single batched Category adapter call (avoids N+1 gRPC calls).
func (s *productQueryService) fillMerchantCategoryNames(ctx context.Context, rows []*models.ProductMerchantRow) error {
	ids := make([]int, 0, len(rows))
	seen := make(map[int32]struct{}, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		if _, ok := seen[row.CategoryID]; ok {
			continue
		}
		seen[row.CategoryID] = struct{}{}
		ids = append(ids, int(row.CategoryID))
	}

	if len(ids) == 0 {
		return nil
	}

	categories, err := s.categoryQuery.FindByIds(ctx, ids)
	if err != nil {
		return err
	}

	nameByID := make(map[int32]string, len(categories))
	for _, category := range categories {
		if category == nil {
			continue
		}
		nameByID[category.CategoryID] = category.Name
	}

	for _, row := range rows {
		if row == nil {
			continue
		}
		row.CategoryName = nameByID[row.CategoryID]
	}

	return nil
}

// fillCategoryNameRows is the ProductCategoryNameRow counterpart of
// fillMerchantCategoryNames.
func (s *productQueryService) fillCategoryNameRows(ctx context.Context, rows []*models.ProductCategoryNameRow) error {
	ids := make([]int, 0, len(rows))
	seen := make(map[int32]struct{}, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		if _, ok := seen[row.CategoryID]; ok {
			continue
		}
		seen[row.CategoryID] = struct{}{}
		ids = append(ids, int(row.CategoryID))
	}

	if len(ids) == 0 {
		return nil
	}

	categories, err := s.categoryQuery.FindByIds(ctx, ids)
	if err != nil {
		return err
	}

	nameByID := make(map[int32]string, len(categories))
	for _, category := range categories {
		if category == nil {
			continue
		}
		nameByID[category.CategoryID] = category.Name
	}

	for _, row := range rows {
		if row == nil {
			continue
		}
		row.CategoryName = nameByID[row.CategoryID]
	}

	return nil
}

func (s *productQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}
