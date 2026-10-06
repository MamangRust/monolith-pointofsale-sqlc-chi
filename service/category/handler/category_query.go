package handler

import (
	"context"
	pbcategories "github.com/MamangRust/monolith-point-of-sale-pb/categories"
	pbcommon "github.com/MamangRust/monolith-point-of-sale-pb/common"
	"go.uber.org/zap"
	"math"

	"github.com/MamangRust/monolith-point-of-sale-category/service"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/category_errors"
)

type categoryQueryHandleGrpc struct {
	pbcategories.UnimplementedCategoryQueryServiceServer
	categoryQuery service.CategoryQueryService
	logger        logger.LoggerInterface
}

func NewCategoryQueryHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pbcategories.CategoryQueryServiceServer {
	return &categoryQueryHandleGrpc{
		categoryQuery: service.CategoryQuery,
		logger:        logger,
	}
}

func (s *categoryQueryHandleGrpc) FindAll(ctx context.Context, request *pbcategories.FindAllCategoryRequest) (*pbcategories.ApiResponsePaginationCategory, error) {
	s.logger.Info("FindAll categories called", zap.Int32("page", request.GetPage()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllCategory{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	category, totalRecords, err := s.categoryQuery.FindAll(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindAll categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindAll categories success", zap.Int("count", len(category)))

	return &pbcategories.ApiResponsePaginationCategory{
		Status:     "success",
		Message:    "Successfully fetched categories",
		Data:       mapResponsesCategory(category),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *categoryQueryHandleGrpc) FindById(ctx context.Context, request *pbcategories.FindByIdCategoryRequest) (*pbcategories.ApiResponseCategory, error) {
	s.logger.Info("FindById category called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidId
	}

	category, err := s.categoryQuery.FindById(ctx, id)
	if err != nil {
		s.logger.Error("FindById category failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindById category success", zap.Int("id", id))

	return &pbcategories.ApiResponseCategory{
		Status:  "success",
		Message: "Successfully fetched category",
		Data:    mapResponseCategory(category),
	}, nil
}

func (s *categoryQueryHandleGrpc) FindByName(ctx context.Context, request *pbcategories.FindByNameCategoryRequest) (*pbcategories.ApiResponseCategory, error) {
	s.logger.Info("FindByName category called", zap.String("name", request.GetName()))

	name := request.GetName()
	if name == "" {
		return nil, category_errors.ErrGrpcFailedInvalidName
	}

	category, err := s.categoryQuery.FindByName(ctx, name)
	if err != nil {
		s.logger.Error("FindByName category failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindByName category success", zap.String("name", name))

	return &pbcategories.ApiResponseCategory{
		Status:  "success",
		Message: "Successfully fetched category",
		Data:    mapResponseCategory(category),
	}, nil
}

func (s *categoryQueryHandleGrpc) FindByIds(ctx context.Context, request *pbcategories.FindByIdsCategoryRequest) (*pbcategories.ApiResponsesCategory, error) {
	s.logger.Info("FindByIds categories called", zap.Int("count", len(request.GetIds())))

	ids := make([]int, 0, len(request.GetIds()))
	for _, id := range request.GetIds() {
		if id <= 0 {
			return nil, category_errors.ErrGrpcFailedInvalidId
		}
		ids = append(ids, int(id))
	}

	categories, err := s.categoryQuery.FindByIds(ctx, ids)
	if err != nil {
		s.logger.Error("FindByIds categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindByIds categories success", zap.Int("count", len(categories)))

	return &pbcategories.ApiResponsesCategory{
		Status:  "success",
		Message: "Successfully fetched categories",
		Data:    mapResponsesCategoryModels(categories),
	}, nil
}

func (s *categoryQueryHandleGrpc) FindByActive(ctx context.Context, request *pbcategories.FindAllCategoryRequest) (*pbcategories.ApiResponsePaginationCategoryDeleteAt, error) {
	s.logger.Info("FindByActive categories called", zap.Int32("page", request.GetPage()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllCategory{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	categories, totalRecords, err := s.categoryQuery.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByActive categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByActive categories success")

	return &pbcategories.ApiResponsePaginationCategoryDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active categories",
		Data:       mapResponsesCategoryActive(categories),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *categoryQueryHandleGrpc) FindByTrashed(ctx context.Context, request *pbcategories.FindAllCategoryRequest) (*pbcategories.ApiResponsePaginationCategoryDeleteAt, error) {
	s.logger.Info("FindByTrashed categories called")

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllCategory{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	categories, totalRecords, err := s.categoryQuery.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByTrashed categories failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByTrashed categories success")

	return &pbcategories.ApiResponsePaginationCategoryDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed categories",
		Data:       mapResponsesCategoryTrashed(categories),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}
