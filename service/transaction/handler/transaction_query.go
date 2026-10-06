package handler

import (
	"context"
	pbcommon "github.com/MamangRust/monolith-point-of-sale-pb/common"
	pbtransactions "github.com/MamangRust/monolith-point-of-sale-pb/transactions"
	"go.uber.org/zap"
	"math"

	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/transaction_errors"
	"github.com/MamangRust/monolith-point-of-sale-transacton/service"
)

type transactionQueryHandleGrpc struct {
	pbtransactions.UnimplementedTransactionQueryServiceServer
	transactionQuery service.TransactionQueryService
	logger           logger.LoggerInterface
}

func NewTransactionQueryHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pbtransactions.TransactionQueryServiceServer {
	return &transactionQueryHandleGrpc{
		transactionQuery: service.TransactionQuery,
		logger:           logger,
	}
}

func (s *transactionQueryHandleGrpc) FindAll(ctx context.Context, request *pbtransactions.FindAllTransactionRequest) (*pbtransactions.ApiResponsePaginationTransaction, error) {
	s.logger.Info("FindAll transactions called", zap.Int32("page", request.GetPage()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllTransaction{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	transaction, totalRecords, err := s.transactionQuery.FindAllTransactions(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindAll transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindAll transactions success")

	return &pbtransactions.ApiResponsePaginationTransaction{
		Status:     "success",
		Message:    "Successfully fetched transaction",
		Data:       mapResponsesTransaction(transaction),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *transactionQueryHandleGrpc) FindByMerchant(ctx context.Context, request *pbtransactions.FindAllTransactionMerchantRequest) (*pbtransactions.ApiResponsePaginationTransaction, error) {
	s.logger.Info("FindByMerchant transactions called", zap.Int32("merchantId", request.GetMerchantId()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()
	merchant_id := int(request.GetMerchantId())

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllTransactionByMerchant{
		MerchantID: merchant_id,
		Page:       page,
		PageSize:   pageSize,
		Search:     search,
	}

	transaction, totalRecords, err := s.transactionQuery.FindByMerchant(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByMerchant transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByMerchant transactions success")

	return &pbtransactions.ApiResponsePaginationTransaction{
		Status:     "success",
		Message:    "Successfully fetched transaction",
		Data:       mapResponsesTransactionByMerchant(transaction),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *transactionQueryHandleGrpc) FindById(ctx context.Context, request *pbtransactions.FindByIdTransactionRequest) (*pbtransactions.ApiResponseTransaction, error) {
	s.logger.Info("FindById transaction called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, transaction_errors.ErrGrpcInvalidID
	}

	transaction, err := s.transactionQuery.FindById(ctx, id)
	if err != nil {
		s.logger.Error("FindById transaction failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindById transaction success")

	return &pbtransactions.ApiResponseTransaction{
		Status:  "success",
		Message: "Successfully fetched transaction",
		Data:    mapResponseTransaction(transaction),
	}, nil
}

func (s *transactionQueryHandleGrpc) FindByActive(ctx context.Context, request *pbtransactions.FindAllTransactionRequest) (*pbtransactions.ApiResponsePaginationTransactionDeleteAt, error) {
	s.logger.Info("FindByActive transactions called", zap.Int32("page", request.GetPage()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllTransaction{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	transaction, totalRecords, err := s.transactionQuery.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByActive transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByActive transactions success")

	return &pbtransactions.ApiResponsePaginationTransactionDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active transaction",
		Data:       mapResponsesTransactionActive(transaction),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *transactionQueryHandleGrpc) FindByTrashed(ctx context.Context, request *pbtransactions.FindAllTransactionRequest) (*pbtransactions.ApiResponsePaginationTransactionDeleteAt, error) {
	s.logger.Info("FindByTrashed transactions called")

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllTransaction{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	transaction, totalRecords, err := s.transactionQuery.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByTrashed transactions failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByTrashed transactions success")

	return &pbtransactions.ApiResponsePaginationTransactionDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed transaction",
		Data:       mapResponsesTransactionTrashed(transaction),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}
