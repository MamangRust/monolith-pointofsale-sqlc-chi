package repository

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/gormerr"
	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"gorm.io/gorm"
)

type merchantDocumentQueryRepository struct {
	db *gorm.DB
}

func NewMerchantDocumentQueryRepository(db *gorm.DB) MerchantDocumentQueryRepository {
	return &merchantDocumentQueryRepository{
		db: db,
	}
}

func (r *merchantDocumentQueryRepository) FindAllDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*models.MerchantDocumentRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*models.MerchantDocumentRow
	query := r.db.WithContext(ctx).Table("merchant_documents").
		Select("document_id, merchant_id, document_type, document_url, status, note, uploaded_at, created_at, updated_at, COUNT(*) OVER() AS total_count").
		Where("deleted_at IS NULL").
		Order("document_id")

	if req.Search != "" {
		query = query.Where("(document_type ILIKE ? OR status ILIKE ? OR note ILIKE ?)",
			"%"+req.Search+"%", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	query = query.Limit(int(req.PageSize)).Offset(int(offset))
	if err := query.Scan(&results).Error; err != nil {
		return nil, nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	var total int
	if len(results) > 0 {
		total = int(results[0].TotalCount)
	}
	return results, &total, nil
}

func (r *merchantDocumentQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*models.MerchantDocumentActiveRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*models.MerchantDocumentActiveRow
	query := r.db.WithContext(ctx).Table("merchant_documents").
		Select("document_id, merchant_id, document_type, document_url, status, note, uploaded_at, created_at, updated_at, deleted_at, COUNT(*) OVER() AS total_count").
		Where("deleted_at IS NULL AND status != 'deleted'").
		Order("document_id")

	if req.Search != "" {
		query = query.Where("(document_type ILIKE ? OR status ILIKE ? OR note ILIKE ?)",
			"%"+req.Search+"%", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	query = query.Limit(int(req.PageSize)).Offset(int(offset))
	if err := query.Scan(&results).Error; err != nil {
		return nil, nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	var total int
	if len(results) > 0 {
		total = int(results[0].TotalCount)
	}
	return results, &total, nil
}

func (r *merchantDocumentQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*models.MerchantDocumentTrashedRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*models.MerchantDocumentTrashedRow
	query := r.db.WithContext(ctx).Table("merchant_documents").
		Select("document_id, merchant_id, document_type, document_url, status, note, uploaded_at, created_at, updated_at, deleted_at, COUNT(*) OVER() AS total_count").
		Where("deleted_at IS NOT NULL").
		Order("document_id")

	if req.Search != "" {
		query = query.Where("(document_type ILIKE ? OR status ILIKE ? OR note ILIKE ?)",
			"%"+req.Search+"%", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	query = query.Limit(int(req.PageSize)).Offset(int(offset))
	if err := query.Scan(&results).Error; err != nil {
		return nil, nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	var total int
	if len(results) > 0 {
		total = int(results[0].TotalCount)
	}
	return results, &total, nil
}

func (r *merchantDocumentQueryRepository) FindById(ctx context.Context, id int) (*models.MerchantDocument, error) {
	var doc models.MerchantDocument
	if err := r.db.WithContext(ctx).
		Where("document_id = ? AND deleted_at IS NULL", id).
		First(&doc).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, sharedErrors.ErrInternal.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &doc, nil
}
