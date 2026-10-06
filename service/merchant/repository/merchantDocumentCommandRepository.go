package repository

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/gormerr"
	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"gorm.io/gorm"
)

type merchantDocumentCommandRepository struct {
	db *gorm.DB
}

func NewMerchantDocumentCommandRepository(db *gorm.DB) MerchantDocumentCommandRepository {
	return &merchantDocumentCommandRepository{
		db: db,
	}
}

func (r *merchantDocumentCommandRepository) CreateMerchantDocument(ctx context.Context, request *requests.CreateMerchantDocumentRequest) (*models.MerchantDocument, error) {
	doc := &models.MerchantDocument{
		MerchantID:   int32(request.MerchantID),
		DocumentType: request.DocumentType,
		DocumentUrl:  request.DocumentUrl,
		Status:       "pending",
	}
	if err := r.db.WithContext(ctx).Create(doc).Error; err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return doc, nil
}

func (r *merchantDocumentCommandRepository) CreateMerchantDocumentInTx(ctx context.Context, tx *gorm.DB, request *requests.CreateMerchantDocumentRequest) (*models.MerchantDocument, error) {
	doc := &models.MerchantDocument{
		MerchantID:   int32(request.MerchantID),
		DocumentType: request.DocumentType,
		DocumentUrl:  request.DocumentUrl,
		Status:       "pending",
	}
	if err := tx.WithContext(ctx).Create(doc).Error; err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return doc, nil
}

func (r *merchantDocumentCommandRepository) UpdateMerchantDocumentStatusInTx(ctx context.Context, tx *gorm.DB, request *requests.UpdateMerchantDocumentStatusRequest) (*models.MerchantDocument, error) {
	if request.DocumentID == nil {
		return nil, sharedErrors.ErrBadRequest.WithMessage("document id is required")
	}

	result := tx.WithContext(ctx).Model(&models.MerchantDocument{}).
		Where("document_id = ?", *request.DocumentID).
		Updates(map[string]interface{}{
			"status":     request.Status,
			"note":       request.Note,
			"updated_at": time.Now(),
		})
	if result.Error != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, sharedErrors.ErrInternal.WithInternal(gorm.ErrRecordNotFound)
	}

	var doc models.MerchantDocument
	if err := tx.WithContext(ctx).Where("document_id = ?", *request.DocumentID).First(&doc).Error; err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &doc, nil
}

func (r *merchantDocumentCommandRepository) UpdateMerchantDocument(ctx context.Context, request *requests.UpdateMerchantDocumentRequest) (*models.MerchantDocument, error) {
	if request.DocumentID == nil {
		return nil, sharedErrors.ErrBadRequest.WithMessage("document id is required")
	}

	result := r.db.WithContext(ctx).Model(&models.MerchantDocument{}).
		Where("document_id = ? AND deleted_at IS NULL", *request.DocumentID).
		Updates(map[string]interface{}{
			"document_type": request.DocumentType,
			"document_url":  request.DocumentUrl,
			"status":        request.Status,
			"note":          request.Note,
			"updated_at":    time.Now(),
		})
	if result.Error != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, sharedErrors.ErrInternal.WithInternal(gorm.ErrRecordNotFound)
	}

	var doc models.MerchantDocument
	if err := r.db.WithContext(ctx).Where("document_id = ?", *request.DocumentID).First(&doc).Error; err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &doc, nil
}

func (r *merchantDocumentCommandRepository) UpdateMerchantDocumentStatus(ctx context.Context, request *requests.UpdateMerchantDocumentStatusRequest) (*models.MerchantDocument, error) {
	if request.DocumentID == nil {
		return nil, sharedErrors.ErrBadRequest.WithMessage("document id is required")
	}

	result := r.db.WithContext(ctx).Model(&models.MerchantDocument{}).
		Where("document_id = ? AND deleted_at IS NULL", *request.DocumentID).
		Updates(map[string]interface{}{
			"status":     request.Status,
			"note":       request.Note,
			"updated_at": time.Now(),
		})
	if result.Error != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, sharedErrors.ErrInternal.WithInternal(gorm.ErrRecordNotFound)
	}

	var doc models.MerchantDocument
	if err := r.db.WithContext(ctx).Where("document_id = ?", *request.DocumentID).First(&doc).Error; err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &doc, nil
}

func (r *merchantDocumentCommandRepository) TrashedMerchantDocument(ctx context.Context, documentID int) (*models.MerchantDocument, error) {
	var doc models.MerchantDocument
	if err := r.db.WithContext(ctx).Where("document_id = ?", documentID).First(&doc).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, sharedErrors.ErrInternal.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	if err := r.db.WithContext(ctx).Delete(&doc).Error; err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &doc, nil
}

func (r *merchantDocumentCommandRepository) RestoreMerchantDocument(ctx context.Context, documentID int) (*models.MerchantDocument, error) {
	var doc models.MerchantDocument
	if err := r.db.WithContext(ctx).Unscoped().Where("document_id = ?", documentID).First(&doc).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, sharedErrors.ErrInternal.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	if err := r.db.WithContext(ctx).Unscoped().Model(&doc).Update("deleted_at", nil).Error; err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &doc, nil
}

func (r *merchantDocumentCommandRepository) DeleteMerchantDocumentPermanent(ctx context.Context, documentID int) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("document_id = ?", documentID).Delete(&models.MerchantDocument{})
	if result.Error != nil {
		return false, sharedErrors.ErrInternal.WithInternal(result.Error)
	}
	return result.RowsAffected > 0, nil
}

func (r *merchantDocumentCommandRepository) RestoreAllMerchantDocument(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Model(&models.MerchantDocument{}).Where("deleted_at IS NOT NULL").Update("deleted_at", nil)
	if result.Error != nil {
		return false, sharedErrors.ErrInternal.WithInternal(result.Error)
	}
	return result.RowsAffected > 0, nil
}

func (r *merchantDocumentCommandRepository) DeleteAllMerchantDocumentPermanent(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("deleted_at IS NOT NULL").Delete(&models.MerchantDocument{})
	if result.Error != nil {
		return false, sharedErrors.ErrInternal.WithInternal(result.Error)
	}
	return result.RowsAffected > 0, nil
}
