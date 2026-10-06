package handler

import (
	"context"
	pbcategories "github.com/MamangRust/monolith-point-of-sale-pb/categories"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/MamangRust/monolith-point-of-sale-category/service"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/category_errors"
)

type categoryCommandHandleGrpc struct {
	pbcategories.UnimplementedCategoryCommandServiceServer
	categoryCommand service.CategoryCommandService
	logger          logger.LoggerInterface
}

func NewCategoryCommandHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pbcategories.CategoryCommandServiceServer {
	return &categoryCommandHandleGrpc{
		categoryCommand: service.CategoryCommand,
		logger:          logger,
	}
}

func (s *categoryCommandHandleGrpc) Create(ctx context.Context, request *pbcategories.CreateCategoryRequest) (*pbcategories.ApiResponseCategory, error) {
	s.logger.Info("Create category called", zap.String("name", request.GetName()))

	req := &requests.CreateCategoryRequest{
		Name:        request.GetName(),
		Description: request.GetDescription(),
	}

	if err := req.Validate(); err != nil {
		return nil, category_errors.ErrGrpcValidateCreateCategory
	}

	category, err := s.categoryCommand.CreateCategory(ctx, req)
	if err != nil {
		s.logger.Error("Create category failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Create category success")

	return &pbcategories.ApiResponseCategory{
		Status:  "success",
		Message: "Successfully created category",
		Data:    mapResponseCategory(category),
	}, nil
}

func (s *categoryCommandHandleGrpc) Update(ctx context.Context, request *pbcategories.UpdateCategoryRequest) (*pbcategories.ApiResponseCategory, error) {
	s.logger.Info("Update category called", zap.Int32("id", request.GetCategoryId()))

	id := int(request.GetCategoryId())
	if id <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidId
	}

	req := &requests.UpdateCategoryRequest{
		CategoryID:  &id,
		Name:        request.GetName(),
		Description: request.GetDescription(),
	}

	if err := req.Validate(); err != nil {
		return nil, category_errors.ErrGrpcValidateUpdateCategory
	}

	category, err := s.categoryCommand.UpdateCategory(ctx, req)
	if err != nil {
		s.logger.Error("Update category failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Update category success")

	return &pbcategories.ApiResponseCategory{
		Status:  "success",
		Message: "Successfully updated category",
		Data:    mapResponseCategory(category),
	}, nil
}

func (s *categoryCommandHandleGrpc) TrashedCategory(ctx context.Context, request *pbcategories.FindByIdCategoryRequest) (*pbcategories.ApiResponseCategoryDeleteAt, error) {
	s.logger.Info("TrashedCategory called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidId
	}

	category, err := s.categoryCommand.TrashedCategory(ctx, id)
	if err != nil {
		s.logger.Error("TrashedCategory failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("TrashedCategory success")

	return &pbcategories.ApiResponseCategoryDeleteAt{
		Status:  "success",
		Message: "Successfully trashed category",
		Data:    mapResponseCategoryDeleteAt(category),
	}, nil
}

func (s *categoryCommandHandleGrpc) RestoreCategory(ctx context.Context, request *pbcategories.FindByIdCategoryRequest) (*pbcategories.ApiResponseCategoryDeleteAt, error) {
	s.logger.Info("RestoreCategory called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidId
	}

	category, err := s.categoryCommand.RestoreCategory(ctx, id)
	if err != nil {
		s.logger.Error("RestoreCategory failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreCategory success")

	return &pbcategories.ApiResponseCategoryDeleteAt{
		Status:  "success",
		Message: "Successfully restored category",
		Data:    mapResponseCategoryDeleteAt(category),
	}, nil
}

func (s *categoryCommandHandleGrpc) DeleteCategoryPermanent(ctx context.Context, request *pbcategories.FindByIdCategoryRequest) (*pbcategories.ApiResponseCategoryDelete, error) {
	s.logger.Info("DeleteCategoryPermanent called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, category_errors.ErrGrpcFailedInvalidId
	}

	_, err := s.categoryCommand.DeleteCategoryPermanent(ctx, id)
	if err != nil {
		s.logger.Error("DeleteCategoryPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteCategoryPermanent success")

	return &pbcategories.ApiResponseCategoryDelete{
		Status:  "success",
		Message: "Successfully deleted category permanently",
	}, nil
}

func (s *categoryCommandHandleGrpc) RestoreAllCategory(ctx context.Context, _ *emptypb.Empty) (*pbcategories.ApiResponseCategoryAll, error) {
	s.logger.Info("RestoreAllCategory called")

	_, err := s.categoryCommand.RestoreAllCategories(ctx)
	if err != nil {
		s.logger.Error("RestoreAllCategory failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreAllCategory success")

	return &pbcategories.ApiResponseCategoryAll{
		Status:  "success",
		Message: "Successfully restore all category",
	}, nil
}

func (s *categoryCommandHandleGrpc) DeleteAllCategoryPermanent(ctx context.Context, _ *emptypb.Empty) (*pbcategories.ApiResponseCategoryAll, error) {
	s.logger.Info("DeleteAllCategoryPermanent called")

	_, err := s.categoryCommand.DeleteAllCategoriesPermanent(ctx)
	if err != nil {
		s.logger.Error("DeleteAllCategoryPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteAllCategoryPermanent success")

	return &pbcategories.ApiResponseCategoryAll{
		Status:  "success",
		Message: "Successfully delete category permanen",
	}, nil
}

// Internal map helpers
