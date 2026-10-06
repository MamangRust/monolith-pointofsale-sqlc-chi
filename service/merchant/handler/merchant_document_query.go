package handler

import (
	"context"
	pbcommon "github.com/MamangRust/monolith-point-of-sale-pb/common"
	pbmerchants_document "github.com/MamangRust/monolith-point-of-sale-pb/merchants_document"
	merchantdocument_errors "github.com/MamangRust/monolith-point-of-sale-shared/errors/merchant_document_errors"
	"go.uber.org/zap"
	"math"

	"github.com/MamangRust/monolith-point-of-sale-merchant/service"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
)

type merchantDocumentQueryHandleGrpc struct {
	pbmerchants_document.UnimplementedMerchantDocumentQueryServiceServer
	merchantDocumentQuery service.MerchantDocumentQueryService
	logger                logger.LoggerInterface
}

func NewMerchantDocumentQueryHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pbmerchants_document.MerchantDocumentQueryServiceServer {
	return &merchantDocumentQueryHandleGrpc{
		merchantDocumentQuery: service.MerchantDocumentQuery,
		logger:                logger,
	}
}

func (s *merchantDocumentQueryHandleGrpc) FindAll(ctx context.Context, req *pbmerchants_document.FindAllMerchantDocumentsRequest) (*pbmerchants_document.ApiResponsePaginationMerchantDocument, error) {
	s.logger.Info("FindAll merchant documents called", zap.Int32("page", req.GetPage()))

	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllMerchantDocuments{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	documents, totalRecords, err := s.merchantDocumentQuery.FindAll(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindAll merchant documents failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindAll merchant documents success")

	return &pbmerchants_document.ApiResponsePaginationMerchantDocument{
		Status:     "success",
		Message:    "Successfully fetched merchant documents",
		Data:       mapResponsesGetMerchantDocumentsRow(documents),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *merchantDocumentQueryHandleGrpc) FindById(ctx context.Context, req *pbmerchants_document.FindMerchantDocumentByIdRequest) (*pbmerchants_document.ApiResponseMerchantDocument, error) {
	s.logger.Info("FindById merchant document called", zap.Int32("id", req.GetDocumentId()))

	id := int(req.GetDocumentId())
	if id <= 0 {
		return nil, merchantdocument_errors.ErrGrpcMerchantInvalidID
	}

	document, err := s.merchantDocumentQuery.FindById(ctx, id)
	if err != nil {
		s.logger.Error("FindById merchant document failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindById merchant document success")

	return &pbmerchants_document.ApiResponseMerchantDocument{
		Status:  "success",
		Message: "Successfully fetched merchant document",
		Data:    mapMerchantDocument(document),
	}, nil
}

func (s *merchantDocumentQueryHandleGrpc) FindAllActive(ctx context.Context, req *pbmerchants_document.FindAllMerchantDocumentsRequest) (*pbmerchants_document.ApiResponsePaginationMerchantDocument, error) {
	s.logger.Info("FindAllActive merchant documents called", zap.Int32("page", req.GetPage()))

	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllMerchantDocuments{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	documents, totalRecords, err := s.merchantDocumentQuery.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindAllActive merchant documents failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindAllActive merchant documents success")

	return &pbmerchants_document.ApiResponsePaginationMerchantDocument{
		Status:     "success",
		Message:    "Successfully fetched active merchant documents",
		Data:       mapResponsesGetActiveMerchantDocumentsRow(documents),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *merchantDocumentQueryHandleGrpc) FindAllTrashed(ctx context.Context, req *pbmerchants_document.FindAllMerchantDocumentsRequest) (*pbmerchants_document.ApiResponsePaginationMerchantDocumentAt, error) {
	s.logger.Info("FindAllTrashed merchant documents called")

	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())
	search := req.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllMerchantDocuments{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	documents, totalRecords, err := s.merchantDocumentQuery.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindAllTrashed merchant documents failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindAllTrashed merchant documents success")

	return &pbmerchants_document.ApiResponsePaginationMerchantDocumentAt{
		Status:     "success",
		Message:    "Successfully fetched trashed merchant documents",
		Data:       mapResponsesGetTrashedMerchantDocumentsRow(documents),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}
