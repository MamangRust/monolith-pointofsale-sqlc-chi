package repository

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/gormerr"
	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"gorm.io/gorm"
)

type transactionCommandRepository struct {
	db *gorm.DB
}

func NewTransactionCommandRepository(db *gorm.DB) TransactionCommandRepository {
	return &transactionCommandRepository{db: db}
}

func (r *transactionCommandRepository) CreateTransaction(ctx context.Context, request *requests.CreateTransactionRequest) (*models.Transaction, error) {
	txn := &models.Transaction{
		OrderID:       int32(request.OrderID),
		MerchantID:    int32(request.MerchantID),
		PaymentMethod: request.PaymentMethod,
		Amount:        int32(request.Amount),
		PaymentStatus: *request.PaymentStatus,
	}
	if err := r.db.WithContext(ctx).Create(txn).Error; err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return txn, nil
}

func (r *transactionCommandRepository) CreateTransactionInTx(ctx context.Context, tx *gorm.DB, request *requests.CreateTransactionRequest) (*models.Transaction, error) {
	txn := &models.Transaction{
		OrderID:       int32(request.OrderID),
		MerchantID:    int32(request.MerchantID),
		PaymentMethod: request.PaymentMethod,
		Amount:        int32(request.Amount),
		PaymentStatus: *request.PaymentStatus,
	}
	if err := tx.WithContext(ctx).Create(txn).Error; err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return txn, nil
}

func (r *transactionCommandRepository) UpdateTransaction(ctx context.Context, request *requests.UpdateTransactionRequest) (*models.Transaction, error) {
	result := r.db.WithContext(ctx).Model(&models.Transaction{}).
		Where("transaction_id = ? AND deleted_at IS NULL", *request.TransactionID).
		Updates(map[string]interface{}{
			"order_id":       int32(request.OrderID),
			"merchant_id":    int32(request.MerchantID),
			"payment_method": request.PaymentMethod,
			"amount":         int32(request.Amount),
			"payment_status": *request.PaymentStatus,
		})
	if result.Error != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, sharedErrors.ErrInternal.WithInternal(gorm.ErrRecordNotFound)
	}
	var txn models.Transaction
	if err := r.db.WithContext(ctx).Where("transaction_id = ?", *request.TransactionID).First(&txn).Error; err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &txn, nil
}

func (r *transactionCommandRepository) TrashTransaction(ctx context.Context, transaction_id int) (*models.Transaction, error) {
	var txn models.Transaction
	if err := r.db.WithContext(ctx).Where("transaction_id = ?", transaction_id).First(&txn).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, sharedErrors.ErrInternal.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	if err := r.db.WithContext(ctx).Delete(&txn).Error; err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &txn, nil
}

func (r *transactionCommandRepository) RestoreTransaction(ctx context.Context, transaction_id int) (*models.Transaction, error) {
	var txn models.Transaction
	if err := r.db.WithContext(ctx).Unscoped().Where("transaction_id = ?", transaction_id).First(&txn).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, sharedErrors.ErrInternal.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	if err := r.db.WithContext(ctx).Unscoped().Model(&txn).Update("deleted_at", nil).Error; err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &txn, nil
}

func (r *transactionCommandRepository) DeleteTransactionPermanently(ctx context.Context, transaction_id int) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("transaction_id = ?", transaction_id).Delete(&models.Transaction{})
	if result.Error != nil {
		return false, sharedErrors.ErrInternal.WithInternal(result.Error)
	}
	return result.RowsAffected > 0, nil
}

func (r *transactionCommandRepository) RestoreAllTransactions(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Model(&models.Transaction{}).Where("deleted_at IS NOT NULL").Update("deleted_at", nil)
	if result.Error != nil {
		return false, sharedErrors.ErrInternal.WithInternal(result.Error)
	}
	return result.RowsAffected > 0, nil
}

func (r *transactionCommandRepository) DeleteAllTransactionPermanent(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().Where("deleted_at IS NOT NULL").Delete(&models.Transaction{})
	if result.Error != nil {
		return false, sharedErrors.ErrInternal.WithInternal(result.Error)
	}
	return result.RowsAffected > 0, nil
}
