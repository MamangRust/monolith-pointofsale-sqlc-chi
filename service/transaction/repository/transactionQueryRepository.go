package repository

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/gormerr"
	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	sharederrors "github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/transaction_errors"
	"gorm.io/gorm"
)

type transactionQueryRepository struct {
	db *gorm.DB
}

func NewTransactionQueryRepository(db *gorm.DB) TransactionQueryRepository {
	return &transactionQueryRepository{db: db}
}

// transactionColumns is the explicit column list shared by the paginated
// transaction listings; it mirrors the original sqlc queries (SELECT * plus
// COUNT(*) OVER () AS total_count).
const transactionColumns = "transaction_id, order_id, merchant_id, payment_method, amount, change_amount, payment_status, created_at, updated_at, deleted_at, COUNT(*) OVER () AS total_count"

func (r *transactionQueryRepository) FindAllTransactions(ctx context.Context, req *requests.FindAllTransaction) ([]*models.TransactionRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var results []*models.TransactionRow
	query := r.db.WithContext(ctx).Table("transactions").
		Select(transactionColumns + ", COUNT(*) OVER () AS total_count")

	if req.Search != "" {
		query = query.Where("(payment_method ILIKE ? OR payment_status ILIKE ?)", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	query = query.Where("deleted_at IS NULL").
		Order("created_at DESC").
		Limit(int(req.PageSize)).
		Offset(int(offset))

	if err := query.Scan(&results).Error; err != nil {
		return nil, nil, transaction_errors.ErrFindAllTransactions.WithInternal(err)
	}

	if len(results) > 0 {
		total := int(results[0].TotalCount)
		return results, &total, nil
	}
	return results, nil, nil
}

func (r *transactionQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllTransaction) ([]*models.TransactionActiveRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var results []*models.TransactionActiveRow
	query := r.db.WithContext(ctx).Table("transactions").
		Select(transactionColumns + ", COUNT(*) OVER () AS total_count")

	if req.Search != "" {
		query = query.Where("(payment_method ILIKE ? OR payment_status ILIKE ?)", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	query = query.Where("deleted_at IS NULL").
		Order("created_at DESC").
		Limit(int(req.PageSize)).
		Offset(int(offset))

	if err := query.Scan(&results).Error; err != nil {
		return nil, nil, transaction_errors.ErrFindByActive.WithInternal(err)
	}

	if len(results) > 0 {
		total := int(results[0].TotalCount)
		return results, &total, nil
	}
	return results, nil, nil
}

func (r *transactionQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllTransaction) ([]*models.TransactionTrashedRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var results []*models.TransactionTrashedRow
	query := r.db.WithContext(ctx).Table("transactions").
		Select(transactionColumns + ", COUNT(*) OVER () AS total_count")

	if req.Search != "" {
		query = query.Where("(payment_method ILIKE ? OR payment_status ILIKE ?)", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	query = query.Where("deleted_at IS NOT NULL").
		Order("created_at DESC").
		Limit(int(req.PageSize)).
		Offset(int(offset))

	if err := query.Scan(&results).Error; err != nil {
		return nil, nil, transaction_errors.ErrFindByTrashed.WithInternal(err)
	}

	if len(results) > 0 {
		total := int(results[0].TotalCount)
		return results, &total, nil
	}
	return results, nil, nil
}

func (r *transactionQueryRepository) FindByMerchant(
	ctx context.Context,
	req *requests.FindAllTransactionByMerchant,
) ([]*models.TransactionByMerchantRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var results []*models.TransactionByMerchantRow
	query := r.db.WithContext(ctx).Table("transactions").
		Select(transactionColumns + ", COUNT(*) OVER () AS total_count")

	if req.Search != "" {
		query = query.Where("(payment_method ILIKE ? OR payment_status ILIKE ?)", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	query = query.Where("deleted_at IS NULL AND merchant_id = ?", req.MerchantID).
		Order("created_at DESC").
		Limit(int(req.PageSize)).
		Offset(int(offset))

	if err := query.Scan(&results).Error; err != nil {
		return nil, nil, transaction_errors.ErrFindByMerchant.WithInternal(err)
	}

	if len(results) > 0 {
		total := int(results[0].TotalCount)
		return results, &total, nil
	}
	return results, nil, nil
}

func (r *transactionQueryRepository) FindById(ctx context.Context, transaction_id int) (*models.Transaction, error) {
	var transaction models.Transaction
	if err := r.db.WithContext(ctx).
		Where("transaction_id = ?", transaction_id).
		First(&transaction).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, transaction_errors.ErrFindById.WithInternal(err)
		}
		return nil, sharederrors.ErrInternal.WithInternal(err)
	}

	return &transaction, nil
}

func (r *transactionQueryRepository) FindByOrderId(ctx context.Context, order_id int) (*models.Transaction, error) {
	var transaction models.Transaction
	if err := r.db.WithContext(ctx).
		Where("order_id = ?", order_id).
		First(&transaction).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, transaction_errors.ErrFindByOrderId.WithInternal(err)
		}
		return nil, sharederrors.ErrInternal.WithInternal(err)
	}

	return &transaction, nil
}
