package mencache

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type MerchantQueryCache interface {
	GetCachedMerchants(ctx context.Context, req *requests.FindAllMerchants) ([]*models.MerchantRow, *int, bool)
	SetCachedMerchants(ctx context.Context, req *requests.FindAllMerchants, data []*models.MerchantRow, total *int)

	GetCachedMerchantActive(ctx context.Context, req *requests.FindAllMerchants) ([]*models.MerchantActiveRow, *int, bool)
	SetCachedMerchantActive(ctx context.Context, req *requests.FindAllMerchants, data []*models.MerchantActiveRow, total *int)

	GetCachedMerchantTrashed(ctx context.Context, req *requests.FindAllMerchants) ([]*models.MerchantTrashedRow, *int, bool)
	SetCachedMerchantTrashed(ctx context.Context, req *requests.FindAllMerchants, data []*models.MerchantTrashedRow, total *int)

	GetCachedMerchant(ctx context.Context, id int) (*models.Merchant, bool)
	SetCachedMerchant(ctx context.Context, data *models.Merchant)

	GetCachedMerchantsByUserId(ctx context.Context, id int) ([]*models.Merchant, bool)
	SetCachedMerchantsByUserId(ctx context.Context, userId int, data []*models.Merchant)
}

type MerchantCommandCache interface {
	DeleteCachedMerchant(ctx context.Context, id int)
	DeleteCachedMerchantAllCache(ctx context.Context)
}

type MerchantDocumentQueryCache interface {
	GetCachedMerchantDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*models.MerchantDocumentRow, *int, bool)
	SetCachedMerchantDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments, data []*models.MerchantDocumentRow, total *int)

	GetCachedMerchantDocumentsActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*models.MerchantDocumentActiveRow, *int, bool)
	SetCachedMerchantDocumentsActive(ctx context.Context, req *requests.FindAllMerchantDocuments, data []*models.MerchantDocumentActiveRow, total *int)

	GetCachedMerchantDocumentsTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*models.MerchantDocumentTrashedRow, *int, bool)
	SetCachedMerchantDocumentsTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments, data []*models.MerchantDocumentTrashedRow, total *int)

	GetCachedMerchantDocument(ctx context.Context, id int) (*models.MerchantDocument, bool)
	SetCachedMerchantDocument(ctx context.Context, data *models.MerchantDocument)
}

type MerchantDocumentCommandCache interface {
	DeleteCachedMerchantDocuments(ctx context.Context, id int)
	DeleteCachedMerchantDocumentsAllCache(ctx context.Context)
}
