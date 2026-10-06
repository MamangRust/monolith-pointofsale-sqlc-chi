package repository

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"gorm.io/gorm"
)

type MerchantQueryRepository interface {
	FindAllMerchants(ctx context.Context, req *requests.FindAllMerchants) ([]*models.MerchantRow, *int, error)
	FindById(ctx context.Context, id int) (*models.Merchant, error)
	FindByActive(ctx context.Context, req *requests.FindAllMerchants) ([]*models.MerchantActiveRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllMerchants) ([]*models.MerchantTrashedRow, *int, error)
}

type MerchantCommandRepository interface {
	CreateMerchant(ctx context.Context, request *requests.CreateMerchantRequest) (*models.Merchant, error)
	UpdateMerchant(ctx context.Context, request *requests.UpdateMerchantRequest) (*models.Merchant, error)
	TrashedMerchant(ctx context.Context, merchant_id int) (*models.Merchant, error)
	RestoreMerchant(ctx context.Context, merchant_id int) (*models.Merchant, error)
	DeleteMerchantPermanent(ctx context.Context, merchant_id int) (bool, error)
	CreateMerchantInTx(ctx context.Context, tx *gorm.DB, request *requests.CreateMerchantRequest) (*models.Merchant, error)
	UpdateMerchantStatusInTx(ctx context.Context, tx *gorm.DB, request *requests.UpdateMerchantStatusRequest) (*models.Merchant, error)
	UpdateMerchantStatus(ctx context.Context, request *requests.UpdateMerchantStatusRequest) (*models.Merchant, error)
	RestoreAllMerchant(ctx context.Context) (bool, error)
	DeleteAllMerchantPermanent(ctx context.Context) (bool, error)
}

type MerchantDocumentQueryRepository interface {
	FindAllDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*models.MerchantDocumentRow, *int, error)
	FindById(ctx context.Context, id int) (*models.MerchantDocument, error)
	FindByActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*models.MerchantDocumentActiveRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*models.MerchantDocumentTrashedRow, *int, error)
}

type MerchantDocumentCommandRepository interface {
	CreateMerchantDocument(ctx context.Context, request *requests.CreateMerchantDocumentRequest) (*models.MerchantDocument, error)
	CreateMerchantDocumentInTx(ctx context.Context, tx *gorm.DB, request *requests.CreateMerchantDocumentRequest) (*models.MerchantDocument, error)
	UpdateMerchantDocumentStatusInTx(ctx context.Context, tx *gorm.DB, request *requests.UpdateMerchantDocumentStatusRequest) (*models.MerchantDocument, error)
	UpdateMerchantDocument(ctx context.Context, request *requests.UpdateMerchantDocumentRequest) (*models.MerchantDocument, error)
	UpdateMerchantDocumentStatus(ctx context.Context, request *requests.UpdateMerchantDocumentStatusRequest) (*models.MerchantDocument, error)
	TrashedMerchantDocument(ctx context.Context, document_id int) (*models.MerchantDocument, error)
	RestoreMerchantDocument(ctx context.Context, document_id int) (*models.MerchantDocument, error)
	DeleteMerchantDocumentPermanent(ctx context.Context, document_id int) (bool, error)
	RestoreAllMerchantDocument(ctx context.Context) (bool, error)
	DeleteAllMerchantDocumentPermanent(ctx context.Context) (bool, error)
}

