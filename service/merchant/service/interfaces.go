package service

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type MerchantQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchants) ([]*models.MerchantRow, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllMerchants) ([]*models.MerchantActiveRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllMerchants) ([]*models.MerchantTrashedRow, *int, error)
	FindById(ctx context.Context, merchantID int) (*models.Merchant, error)
}

type MerchantCommandService interface {
	CreateMerchant(ctx context.Context, req *requests.CreateMerchantRequest) (*models.Merchant, error)
	UpdateMerchant(ctx context.Context, req *requests.UpdateMerchantRequest) (*models.Merchant, error)
	UpdateMerchantStatus(ctx context.Context, req *requests.UpdateMerchantStatusRequest) (*models.Merchant, error)
	TrashedMerchant(ctx context.Context, merchantID int) (*models.Merchant, error)
	RestoreMerchant(ctx context.Context, merchantID int) (*models.Merchant, error)
	DeleteMerchantPermanent(ctx context.Context, merchantID int) (bool, error)
	RestoreAllMerchant(ctx context.Context) (bool, error)
	DeleteAllMerchantPermanent(ctx context.Context) (bool, error)
}

type MerchantDocumentQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*models.MerchantDocumentRow, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*models.MerchantDocumentActiveRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*models.MerchantDocumentTrashedRow, *int, error)
	FindById(ctx context.Context, documentID int) (*models.MerchantDocument, error)
}

type MerchantDocumentCommandService interface {
	CreateMerchantDocument(ctx context.Context, request *requests.CreateMerchantDocumentRequest) (*models.MerchantDocument, error)
	UpdateMerchantDocument(ctx context.Context, request *requests.UpdateMerchantDocumentRequest) (*models.MerchantDocument, error)
	UpdateMerchantDocumentStatus(ctx context.Context, request *requests.UpdateMerchantDocumentStatusRequest) (*models.MerchantDocument, error)
	TrashedMerchantDocument(ctx context.Context, documentID int) (*models.MerchantDocument, error)
	RestoreMerchantDocument(ctx context.Context, documentID int) (*models.MerchantDocument, error)
	DeleteMerchantDocumentPermanent(ctx context.Context, documentID int) (bool, error)
	RestoreAllMerchantDocument(ctx context.Context) (bool, error)
	DeleteAllMerchantDocumentPermanent(ctx context.Context) (bool, error)
}
