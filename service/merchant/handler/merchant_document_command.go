package handler

import (
	"context"
	pbmerchants_document "github.com/MamangRust/monolith-point-of-sale-pb/merchants_document"
	merchantdocument_errors "github.com/MamangRust/monolith-point-of-sale-shared/errors/merchant_document_errors"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/MamangRust/monolith-point-of-sale-merchant/service"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
)

type merchantDocumentCommandHandleGrpc struct {
	pbmerchants_document.UnimplementedMerchantDocumentCommandServiceServer
	merchantDocumentCommand service.MerchantDocumentCommandService
	logger                  logger.LoggerInterface
}

func NewMerchantDocumentCommandHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pbmerchants_document.MerchantDocumentCommandServiceServer {
	return &merchantDocumentCommandHandleGrpc{
		merchantDocumentCommand: service.MerchantDocumentCommand,
		logger:                  logger,
	}
}

func (s *merchantDocumentCommandHandleGrpc) Create(ctx context.Context, req *pbmerchants_document.CreateMerchantDocumentRequest) (*pbmerchants_document.ApiResponseMerchantDocument, error) {
	s.logger.Info("Create merchant document called", zap.Int32("merchantId", req.GetMerchantId()))

	request := requests.CreateMerchantDocumentRequest{
		MerchantID:   int(req.GetMerchantId()),
		DocumentType: req.GetDocumentType(),
		DocumentUrl:  req.GetDocumentUrl(),
	}

	if err := request.Validate(); err != nil {
		return nil, merchantdocument_errors.ErrGrpcValidateCreateMerchantDocument
	}

	document, err := s.merchantDocumentCommand.CreateMerchantDocument(ctx, &request)
	if err != nil {
		s.logger.Error("Create merchant document failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Create merchant document success")

	return &pbmerchants_document.ApiResponseMerchantDocument{
		Status:  "success",
		Message: "Successfully created merchant document",
		Data:    mapMerchantDocument(document),
	}, nil
}

func (s *merchantDocumentCommandHandleGrpc) Update(ctx context.Context, req *pbmerchants_document.UpdateMerchantDocumentRequest) (*pbmerchants_document.ApiResponseMerchantDocument, error) {
	s.logger.Info("Update merchant document called", zap.Int32("id", req.GetDocumentId()))

	id := int(req.GetDocumentId())
	if id <= 0 {
		return nil, merchantdocument_errors.ErrGrpcMerchantInvalidID
	}

	request := requests.UpdateMerchantDocumentRequest{
		DocumentID:   &id,
		MerchantID:   int(req.GetMerchantId()),
		DocumentType: req.GetDocumentType(),
		DocumentUrl:  req.GetDocumentUrl(),
		Status:       req.GetStatus(),
		Note:         req.GetNote(),
	}

	if err := request.Validate(); err != nil {
		return nil, merchantdocument_errors.ErrGrpcFailedUpdateMerchantDocument
	}

	document, err := s.merchantDocumentCommand.UpdateMerchantDocument(ctx, &request)
	if err != nil {
		s.logger.Error("Update merchant document failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Update merchant document success")

	return &pbmerchants_document.ApiResponseMerchantDocument{
		Status:  "success",
		Message: "Successfully updated merchant document",
		Data:    mapMerchantDocument(document),
	}, nil
}

func (s *merchantDocumentCommandHandleGrpc) UpdateStatus(ctx context.Context, req *pbmerchants_document.UpdateMerchantDocumentStatusRequest) (*pbmerchants_document.ApiResponseMerchantDocument, error) {
	s.logger.Info("UpdateStatus merchant document called", zap.Int32("id", req.GetDocumentId()))

	id := int(req.GetDocumentId())
	if id <= 0 {
		return nil, merchantdocument_errors.ErrGrpcMerchantInvalidID
	}

	request := requests.UpdateMerchantDocumentStatusRequest{
		DocumentID: &id,
		MerchantID: int(req.GetMerchantId()),
		Status:     req.GetStatus(),
		Note:       req.GetNote(),
	}

	if err := request.Validate(); err != nil {
		return nil, merchantdocument_errors.ErrGrpcFailedUpdateMerchantDocument
	}

	document, err := s.merchantDocumentCommand.UpdateMerchantDocumentStatus(ctx, &request)
	if err != nil {
		s.logger.Error("UpdateStatus merchant document failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("UpdateStatus merchant document success")

	return &pbmerchants_document.ApiResponseMerchantDocument{
		Status:  "success",
		Message: "Successfully updated merchant document status",
		Data:    mapMerchantDocument(document),
	}, nil
}

func (s *merchantDocumentCommandHandleGrpc) Trashed(ctx context.Context, req *pbmerchants_document.TrashedMerchantDocumentRequest) (*pbmerchants_document.ApiResponseMerchantDocument, error) {
	s.logger.Info("Trashed merchant document called", zap.Int32("id", req.GetDocumentId()))

	id := int(req.GetDocumentId())
	if id <= 0 {
		return nil, merchantdocument_errors.ErrGrpcMerchantInvalidID
	}

	document, err := s.merchantDocumentCommand.TrashedMerchantDocument(ctx, id)
	if err != nil {
		s.logger.Error("Trashed merchant document failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Trashed merchant document success")

	return &pbmerchants_document.ApiResponseMerchantDocument{
		Status:  "success",
		Message: "Successfully trashed merchant document",
		Data:    mapMerchantDocument(document),
	}, nil
}

func (s *merchantDocumentCommandHandleGrpc) Restore(ctx context.Context, req *pbmerchants_document.RestoreMerchantDocumentRequest) (*pbmerchants_document.ApiResponseMerchantDocument, error) {
	s.logger.Info("Restore merchant document called", zap.Int32("id", req.GetDocumentId()))

	id := int(req.GetDocumentId())
	if id <= 0 {
		return nil, merchantdocument_errors.ErrGrpcMerchantInvalidID
	}

	document, err := s.merchantDocumentCommand.RestoreMerchantDocument(ctx, id)
	if err != nil {
		s.logger.Error("Restore merchant document failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Restore merchant document success")

	return &pbmerchants_document.ApiResponseMerchantDocument{
		Status:  "success",
		Message: "Successfully restored merchant document",
		Data:    mapMerchantDocument(document),
	}, nil
}

func (s *merchantDocumentCommandHandleGrpc) DeletePermanent(ctx context.Context, req *pbmerchants_document.DeleteMerchantDocumentPermanentRequest) (*pbmerchants_document.ApiResponseMerchantDocumentDelete, error) {
	s.logger.Info("DeletePermanent merchant document called", zap.Int32("id", req.GetDocumentId()))

	id := int(req.GetDocumentId())
	if id <= 0 {
		return nil, merchantdocument_errors.ErrGrpcMerchantInvalidID
	}

	_, err := s.merchantDocumentCommand.DeleteMerchantDocumentPermanent(ctx, id)
	if err != nil {
		s.logger.Error("DeletePermanent merchant document failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeletePermanent merchant document success")

	return &pbmerchants_document.ApiResponseMerchantDocumentDelete{
		Status:  "success",
		Message: "Successfully permanently deleted merchant document",
	}, nil
}

func (s *merchantDocumentCommandHandleGrpc) RestoreAll(ctx context.Context, _ *emptypb.Empty) (*pbmerchants_document.ApiResponseMerchantDocumentAll, error) {
	s.logger.Info("RestoreAll merchant documents called")

	_, err := s.merchantDocumentCommand.RestoreAllMerchantDocument(ctx)
	if err != nil {
		s.logger.Error("RestoreAll merchant documents failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreAll merchant documents success")

	return &pbmerchants_document.ApiResponseMerchantDocumentAll{
		Status:  "success",
		Message: "Successfully restored all merchant documents",
	}, nil
}

func (s *merchantDocumentCommandHandleGrpc) DeleteAllPermanent(ctx context.Context, _ *emptypb.Empty) (*pbmerchants_document.ApiResponseMerchantDocumentAll, error) {
	s.logger.Info("DeleteAllPermanent merchant documents called")

	_, err := s.merchantDocumentCommand.DeleteAllMerchantDocumentPermanent(ctx)
	if err != nil {
		s.logger.Error("DeleteAllPermanent merchant documents failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteAllPermanent merchant documents success")

	return &pbmerchants_document.ApiResponseMerchantDocumentAll{
		Status:  "success",
		Message: "Successfully permanently deleted all merchant documents",
	}, nil
}

// Map helpers
