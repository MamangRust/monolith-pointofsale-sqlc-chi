package handler

import (
	"context"
	pbcommon "github.com/MamangRust/monolith-point-of-sale-pb/common"
	pbmerchants "github.com/MamangRust/monolith-point-of-sale-pb/merchants"
	"go.uber.org/zap"
	"math"

	"github.com/MamangRust/monolith-point-of-sale-merchant/service"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/merchant_errors"
)

type merchantQueryHandleGrpc struct {
	pbmerchants.UnimplementedMerchantQueryServiceServer
	merchantQuery service.MerchantQueryService
	logger        logger.LoggerInterface
}

func NewMerchantQueryHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pbmerchants.MerchantQueryServiceServer {
	return &merchantQueryHandleGrpc{
		merchantQuery: service.MerchantQuery,
		logger:        logger,
	}
}

func (s *merchantQueryHandleGrpc) FindAll(ctx context.Context, req *pbmerchants.FindAllMerchantRequest) (*pbmerchants.ApiResponsePaginationMerchant, error) {
	s.logger.Info("FindAll merchant called", zap.Int32("page", req.GetPage()))

	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllMerchants{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchants, totalRecords, err := s.merchantQuery.FindAll(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindAll merchant failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindAll merchant success")

	return &pbmerchants.ApiResponsePaginationMerchant{
		Status:     "success",
		Message:    "Successfully fetched merchant record",
		Data:       mapResponsesGetMerchantsRow(merchants),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *merchantQueryHandleGrpc) FindByActive(ctx context.Context, req *pbmerchants.FindAllMerchantRequest) (*pbmerchants.ApiResponsePaginationMerchantDeleteAt, error) {
	s.logger.Info("FindByActive merchants called", zap.Int32("page", req.GetPage()))

	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllMerchants{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	res, totalRecords, err := s.merchantQuery.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByActive merchants failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByActive merchants success")

	return &pbmerchants.ApiResponsePaginationMerchantDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched merchant record",
		Data:       mapResponsesGetMerchantsActiveRow(res),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *merchantQueryHandleGrpc) FindByTrashed(ctx context.Context, req *pbmerchants.FindAllMerchantRequest) (*pbmerchants.ApiResponsePaginationMerchantDeleteAt, error) {
	s.logger.Info("FindByTrashed merchants called")

	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllMerchants{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	res, totalRecords, err := s.merchantQuery.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByTrashed merchants failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByTrashed merchants success")

	return &pbmerchants.ApiResponsePaginationMerchantDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched merchant record",
		Data:       mapResponsesGetMerchantsTrashedRow(res),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *merchantQueryHandleGrpc) FindById(ctx context.Context, request *pbmerchants.FindByIdMerchantRequest) (*pbmerchants.ApiResponseMerchant, error) {
	s.logger.Info("FindById merchant called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, merchant_errors.ErrGrpcInvalidID
	}

	merchant, err := s.merchantQuery.FindById(ctx, id)
	if err != nil {
		s.logger.Error("FindById merchant failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindById merchant success")

	return &pbmerchants.ApiResponseMerchant{
		Status:  "success",
		Message: "Successfully fetched merchant",
		Data:    mapResponseMerchant(merchant),
	}, nil
}
