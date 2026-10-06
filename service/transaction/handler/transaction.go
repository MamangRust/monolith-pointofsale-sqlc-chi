package handler

import (
	pbcommon "github.com/MamangRust/monolith-point-of-sale-pb/common"
	pbtransactions "github.com/MamangRust/monolith-point-of-sale-pb/transactions"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/convert"
)

func mapPaginationMeta(meta *pbcommon.PaginationMeta) *pbcommon.PaginationMeta {
	if meta == nil {
		return nil
	}
	return &pbcommon.PaginationMeta{
		CurrentPage:  meta.CurrentPage,
		PageSize:     meta.PageSize,
		TotalPages:   meta.TotalPages,
		TotalRecords: meta.TotalRecords,
	}
}

func mapResponseTransaction(transaction *models.Transaction) *pbtransactions.TransactionResponse {
	if transaction == nil {
		return nil
	}
	var createdAtStr string
	if !transaction.CreatedAt.IsZero() {
		createdAtStr = convert.FormatTimePtr(&transaction.CreatedAt)
	}
	var updatedAtStr string
	if !transaction.UpdatedAt.IsZero() {
		updatedAtStr = convert.FormatTimePtr(&transaction.UpdatedAt)
	}
	var changeAmount int32
	if transaction.ChangeAmount != nil {
		changeAmount = *transaction.ChangeAmount
	}
	return &pbtransactions.TransactionResponse{
		Id:            transaction.TransactionID,
		OrderId:       transaction.OrderID,
		MerchantId:    transaction.MerchantID,
		PaymentMethod: transaction.PaymentMethod,
		Amount:        transaction.Amount,
		ChangeAmount:  changeAmount,
		PaymentStatus: transaction.PaymentStatus,
		CreatedAt:     createdAtStr,
		UpdatedAt:     updatedAtStr,
	}
}

func mapResponsesTransaction(transactions []*models.TransactionRow) []*pbtransactions.TransactionResponse {
	var mappedTransactions []*pbtransactions.TransactionResponse
	for _, t := range transactions {
		if t == nil {
			continue
		}
		var createdAtStr string
		if !t.CreatedAt.IsZero() {
			createdAtStr = convert.FormatTimePtr(&t.CreatedAt)
		}
		var updatedAtStr string
		if !t.UpdatedAt.IsZero() {
			updatedAtStr = convert.FormatTimePtr(&t.UpdatedAt)
		}
		var changeAmount int32
		if t.ChangeAmount != nil {
			changeAmount = *t.ChangeAmount
		}
		mappedTransactions = append(mappedTransactions, &pbtransactions.TransactionResponse{
			Id:            t.TransactionID,
			OrderId:       t.OrderID,
			MerchantId:    t.MerchantID,
			PaymentMethod: t.PaymentMethod,
			Amount:        t.Amount,
			ChangeAmount:  changeAmount,
			PaymentStatus: t.PaymentStatus,
			CreatedAt:     createdAtStr,
			UpdatedAt:     updatedAtStr,
		})
	}
	return mappedTransactions
}

func mapResponsesTransactionByMerchant(transactions []*models.TransactionByMerchantRow) []*pbtransactions.TransactionResponse {
	var mappedTransactions []*pbtransactions.TransactionResponse
	for _, t := range transactions {
		if t == nil {
			continue
		}
		var createdAtStr string
		if !t.CreatedAt.IsZero() {
			createdAtStr = convert.FormatTimePtr(&t.CreatedAt)
		}
		var updatedAtStr string
		if !t.UpdatedAt.IsZero() {
			updatedAtStr = convert.FormatTimePtr(&t.UpdatedAt)
		}
		var changeAmount int32
		if t.ChangeAmount != nil {
			changeAmount = *t.ChangeAmount
		}
		mappedTransactions = append(mappedTransactions, &pbtransactions.TransactionResponse{
			Id:            t.TransactionID,
			OrderId:       t.OrderID,
			MerchantId:    t.MerchantID,
			PaymentMethod: t.PaymentMethod,
			Amount:        t.Amount,
			ChangeAmount:  changeAmount,
			PaymentStatus: t.PaymentStatus,
			CreatedAt:     createdAtStr,
			UpdatedAt:     updatedAtStr,
		})
	}
	return mappedTransactions
}

func mapResponseTransactionDeleteAt(transaction *models.Transaction) *pbtransactions.TransactionResponseDeleteAt {
	if transaction == nil {
		return nil
	}
	var createdAtStr string
	if !transaction.CreatedAt.IsZero() {
		createdAtStr = convert.FormatTimePtr(&transaction.CreatedAt)
	}
	var updatedAtStr string
	if !transaction.UpdatedAt.IsZero() {
		updatedAtStr = convert.FormatTimePtr(&transaction.UpdatedAt)
	}
	var deletedAt *wrapperspb.StringValue
	if transaction.DeletedAt.Valid {
		deletedAt = convert.TimeToWrappers(&transaction.DeletedAt.Time)
	}
	var changeAmount int32
	if transaction.ChangeAmount != nil {
		changeAmount = *transaction.ChangeAmount
	}

	return &pbtransactions.TransactionResponseDeleteAt{
		Id:            transaction.TransactionID,
		OrderId:       transaction.OrderID,
		MerchantId:    transaction.MerchantID,
		PaymentMethod: transaction.PaymentMethod,
		Amount:        transaction.Amount,
		ChangeAmount:  changeAmount,
		PaymentStatus: transaction.PaymentStatus,
		CreatedAt:     createdAtStr,
		UpdatedAt:     updatedAtStr,
		DeletedAt:     deletedAt,
	}
}

func mapResponsesTransactionActive(transactions []*models.TransactionActiveRow) []*pbtransactions.TransactionResponseDeleteAt {
	var mappedTransactions []*pbtransactions.TransactionResponseDeleteAt
	for _, t := range transactions {
		if t == nil {
			continue
		}
		var createdAtStr string
		if !t.CreatedAt.IsZero() {
			createdAtStr = convert.FormatTimePtr(&t.CreatedAt)
		}
		var updatedAtStr string
		if !t.UpdatedAt.IsZero() {
			updatedAtStr = convert.FormatTimePtr(&t.UpdatedAt)
		}
		var deletedAt *wrapperspb.StringValue
		if t.DeletedAt != nil {
			deletedAt = wrapperspb.String(t.DeletedAt.Format("2006-01-02 15:04:05"))
		}
		var changeAmount int32
		if t.ChangeAmount != nil {
			changeAmount = *t.ChangeAmount
		}
		mappedTransactions = append(mappedTransactions, &pbtransactions.TransactionResponseDeleteAt{
			Id:            t.TransactionID,
			OrderId:       t.OrderID,
			MerchantId:    t.MerchantID,
			PaymentMethod: t.PaymentMethod,
			Amount:        t.Amount,
			ChangeAmount:  changeAmount,
			PaymentStatus: t.PaymentStatus,
			CreatedAt:     createdAtStr,
			UpdatedAt:     updatedAtStr,
			DeletedAt:     deletedAt,
		})
	}
	return mappedTransactions
}

func mapResponsesTransactionTrashed(transactions []*models.TransactionTrashedRow) []*pbtransactions.TransactionResponseDeleteAt {
	var mappedTransactions []*pbtransactions.TransactionResponseDeleteAt
	for _, t := range transactions {
		if t == nil {
			continue
		}
		var createdAtStr string
		if !t.CreatedAt.IsZero() {
			createdAtStr = convert.FormatTimePtr(&t.CreatedAt)
		}
		var updatedAtStr string
		if !t.UpdatedAt.IsZero() {
			updatedAtStr = convert.FormatTimePtr(&t.UpdatedAt)
		}
		var deletedAt *wrapperspb.StringValue
		if t.DeletedAt != nil {
			deletedAt = wrapperspb.String(t.DeletedAt.Format("2006-01-02 15:04:05"))
		}
		var changeAmount int32
		if t.ChangeAmount != nil {
			changeAmount = *t.ChangeAmount
		}
		mappedTransactions = append(mappedTransactions, &pbtransactions.TransactionResponseDeleteAt{
			Id:            t.TransactionID,
			OrderId:       t.OrderID,
			MerchantId:    t.MerchantID,
			PaymentMethod: t.PaymentMethod,
			Amount:        t.Amount,
			ChangeAmount:  changeAmount,
			PaymentStatus: t.PaymentStatus,
			CreatedAt:     createdAtStr,
			UpdatedAt:     updatedAtStr,
			DeletedAt:     deletedAt,
		})
	}
	return mappedTransactions
}

func mapResponseTransactionMonthAmountSuccess(row *models.TransactionMonthlyAmountSuccessRow) *pbtransactions.TransactionMonthlyAmountSuccess {
	if row == nil {
		return nil
	}
	return &pbtransactions.TransactionMonthlyAmountSuccess{
		Year:         row.Year,
		Month:        row.Month,
		TotalSuccess: int32(row.TotalSuccess),
		TotalAmount:  int32(row.TotalAmount),
	}
}

func mapResponsesTransactionMonthlyAmountSuccess(rows []*models.TransactionMonthlyAmountSuccessRow) []*pbtransactions.TransactionMonthlyAmountSuccess {
	var transaction []*pbtransactions.TransactionMonthlyAmountSuccess
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionMonthAmountSuccess(row))
	}
	return transaction
}

func mapResponseTransactionYearAmountSuccess(row *models.TransactionYearlyAmountSuccessRow) *pbtransactions.TransactionYearlyAmountSuccess {
	if row == nil {
		return nil
	}
	return &pbtransactions.TransactionYearlyAmountSuccess{
		Year:         row.Year,
		TotalSuccess: int32(row.TotalSuccess),
		TotalAmount:  int32(row.TotalAmount),
	}
}

func mapResponsesTransactionYearlyAmountSuccess(rows []*models.TransactionYearlyAmountSuccessRow) []*pbtransactions.TransactionYearlyAmountSuccess {
	var transaction []*pbtransactions.TransactionYearlyAmountSuccess
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionYearAmountSuccess(row))
	}
	return transaction
}

func mapResponseTransactionMonthAmountFailed(row *models.TransactionMonthlyAmountFailedRow) *pbtransactions.TransactionMonthlyAmountFailed {
	if row == nil {
		return nil
	}
	return &pbtransactions.TransactionMonthlyAmountFailed{
		Year:        row.Year,
		Month:       row.Month,
		TotalFailed: int32(row.TotalFailed),
		TotalAmount: int32(row.TotalAmount),
	}
}

func mapResponsesTransactionMonthlyAmountFailed(rows []*models.TransactionMonthlyAmountFailedRow) []*pbtransactions.TransactionMonthlyAmountFailed {
	var transaction []*pbtransactions.TransactionMonthlyAmountFailed
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionMonthAmountFailed(row))
	}
	return transaction
}

func mapResponseTransactionYearAmountFailed(row *models.TransactionYearlyAmountFailedRow) *pbtransactions.TransactionYearlyAmountFailed {
	if row == nil {
		return nil
	}
	return &pbtransactions.TransactionYearlyAmountFailed{
		Year:        row.Year,
		TotalFailed: int32(row.TotalFailed),
		TotalAmount: int32(row.TotalAmount),
	}
}

func mapResponsesTransactionYearlyAmountFailed(rows []*models.TransactionYearlyAmountFailedRow) []*pbtransactions.TransactionYearlyAmountFailed {
	var transaction []*pbtransactions.TransactionYearlyAmountFailed
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionYearAmountFailed(row))
	}
	return transaction
}

func mapResponseTransactionMonthAmountSuccessByMerchant(row *models.TransactionMonthlyAmountSuccessByMerchantRow) *pbtransactions.TransactionMonthlyAmountSuccess {
	if row == nil {
		return nil
	}
	return &pbtransactions.TransactionMonthlyAmountSuccess{
		Year:         row.Year,
		Month:        row.Month,
		TotalSuccess: int32(row.TotalSuccess),
		TotalAmount:  int32(row.TotalAmount),
	}
}

func mapResponsesTransactionMonthlyAmountSuccessByMerchant(rows []*models.TransactionMonthlyAmountSuccessByMerchantRow) []*pbtransactions.TransactionMonthlyAmountSuccess {
	var transaction []*pbtransactions.TransactionMonthlyAmountSuccess
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionMonthAmountSuccessByMerchant(row))
	}
	return transaction
}

func mapResponseTransactionYearAmountSuccessByMerchant(row *models.TransactionYearlyAmountSuccessByMerchantRow) *pbtransactions.TransactionYearlyAmountSuccess {
	if row == nil {
		return nil
	}
	return &pbtransactions.TransactionYearlyAmountSuccess{
		Year:         row.Year,
		TotalSuccess: int32(row.TotalSuccess),
		TotalAmount:  int32(row.TotalAmount),
	}
}

func mapResponsesTransactionYearlyAmountSuccessByMerchant(rows []*models.TransactionYearlyAmountSuccessByMerchantRow) []*pbtransactions.TransactionYearlyAmountSuccess {
	var transaction []*pbtransactions.TransactionYearlyAmountSuccess
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionYearAmountSuccessByMerchant(row))
	}
	return transaction
}

func mapResponseTransactionMonthAmountFailedByMerchant(row *models.TransactionMonthlyAmountFailedByMerchantRow) *pbtransactions.TransactionMonthlyAmountFailed {
	if row == nil {
		return nil
	}
	return &pbtransactions.TransactionMonthlyAmountFailed{
		Year:        row.Year,
		Month:       row.Month,
		TotalFailed: int32(row.TotalFailed),
		TotalAmount: int32(row.TotalAmount),
	}
}

func mapResponsesTransactionMonthlyAmountFailedByMerchant(rows []*models.TransactionMonthlyAmountFailedByMerchantRow) []*pbtransactions.TransactionMonthlyAmountFailed {
	var transaction []*pbtransactions.TransactionMonthlyAmountFailed
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionMonthAmountFailedByMerchant(row))
	}
	return transaction
}

func mapResponseTransactionYearAmountFailedByMerchant(row *models.TransactionYearlyAmountFailedByMerchantRow) *pbtransactions.TransactionYearlyAmountFailed {
	if row == nil {
		return nil
	}
	return &pbtransactions.TransactionYearlyAmountFailed{
		Year:        row.Year,
		TotalFailed: int32(row.TotalFailed),
		TotalAmount: int32(row.TotalAmount),
	}
}

func mapResponsesTransactionYearlyAmountFailedByMerchant(rows []*models.TransactionYearlyAmountFailedByMerchantRow) []*pbtransactions.TransactionYearlyAmountFailed {
	var transaction []*pbtransactions.TransactionYearlyAmountFailed
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionYearAmountFailedByMerchant(row))
	}
	return transaction
}

func mapResponseTransactionMonthMethodSuccess(row *models.TransactionMonthlyMethodSuccessRow) *pbtransactions.TransactionMonthlyMethod {
	if row == nil {
		return nil
	}
	return &pbtransactions.TransactionMonthlyMethod{
		Month:             row.Month,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int32(row.TotalTransactions),
		TotalAmount:       int32(row.TotalAmount),
	}
}

func mapResponsesTransactionMonthlyMethodSuccess(rows []*models.TransactionMonthlyMethodSuccessRow) []*pbtransactions.TransactionMonthlyMethod {
	var transaction []*pbtransactions.TransactionMonthlyMethod
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionMonthMethodSuccess(row))
	}
	return transaction
}

func mapResponseTransactionYearlyMethodSuccess(row *models.TransactionYearlyMethodSuccessRow) *pbtransactions.TransactionYearlyMethod {
	if row == nil {
		return nil
	}
	return &pbtransactions.TransactionYearlyMethod{
		Year:              row.Year,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int32(row.TotalTransactions),
		TotalAmount:       int32(row.TotalAmount),
	}
}

func mapResponsesTransactionYearlyMethodSuccess(rows []*models.TransactionYearlyMethodSuccessRow) []*pbtransactions.TransactionYearlyMethod {
	var transaction []*pbtransactions.TransactionYearlyMethod
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionYearlyMethodSuccess(row))
	}
	return transaction
}

func mapResponseTransactionMonthMethodFailed(row *models.TransactionMonthlyMethodFailedRow) *pbtransactions.TransactionMonthlyMethod {
	if row == nil {
		return nil
	}
	return &pbtransactions.TransactionMonthlyMethod{
		Month:             row.Month,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int32(row.TotalTransactions),
		TotalAmount:       int32(row.TotalAmount),
	}
}

func mapResponsesTransactionMonthlyMethodFailed(rows []*models.TransactionMonthlyMethodFailedRow) []*pbtransactions.TransactionMonthlyMethod {
	var transaction []*pbtransactions.TransactionMonthlyMethod
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionMonthMethodFailed(row))
	}
	return transaction
}

func mapResponseTransactionYearlyMethodFailed(row *models.TransactionYearlyMethodFailedRow) *pbtransactions.TransactionYearlyMethod {
	if row == nil {
		return nil
	}
	return &pbtransactions.TransactionYearlyMethod{
		Year:              row.Year,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int32(row.TotalTransactions),
		TotalAmount:       int32(row.TotalAmount),
	}
}

func mapResponsesTransactionYearlyMethodFailed(rows []*models.TransactionYearlyMethodFailedRow) []*pbtransactions.TransactionYearlyMethod {
	var transaction []*pbtransactions.TransactionYearlyMethod
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionYearlyMethodFailed(row))
	}
	return transaction
}

func mapResponseTransactionMonthMethodByMerchantSuccess(row *models.TransactionMonthlyMethodByMerchantSuccessRow) *pbtransactions.TransactionMonthlyMethod {
	if row == nil {
		return nil
	}
	return &pbtransactions.TransactionMonthlyMethod{
		Month:             row.Month,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int32(row.TotalTransactions),
		TotalAmount:       int32(row.TotalAmount),
	}
}

func mapResponsesTransactionMonthlyMethodByMerchantSuccess(rows []*models.TransactionMonthlyMethodByMerchantSuccessRow) []*pbtransactions.TransactionMonthlyMethod {
	var transaction []*pbtransactions.TransactionMonthlyMethod
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionMonthMethodByMerchantSuccess(row))
	}
	return transaction
}

func mapResponseTransactionYearlyMethodByMerchantSuccess(row *models.TransactionYearlyMethodByMerchantSuccessRow) *pbtransactions.TransactionYearlyMethod {
	if row == nil {
		return nil
	}
	return &pbtransactions.TransactionYearlyMethod{
		Year:              row.Year,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int32(row.TotalTransactions),
		TotalAmount:       int32(row.TotalAmount),
	}
}

func mapResponsesTransactionYearlyMethodByMerchantSuccess(rows []*models.TransactionYearlyMethodByMerchantSuccessRow) []*pbtransactions.TransactionYearlyMethod {
	var transaction []*pbtransactions.TransactionYearlyMethod
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionYearlyMethodByMerchantSuccess(row))
	}
	return transaction
}

func mapResponseTransactionMonthMethodByMerchantFailed(row *models.TransactionMonthlyMethodByMerchantFailedRow) *pbtransactions.TransactionMonthlyMethod {
	if row == nil {
		return nil
	}
	return &pbtransactions.TransactionMonthlyMethod{
		Month:             row.Month,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int32(row.TotalTransactions),
		TotalAmount:       int32(row.TotalAmount),
	}
}

func mapResponsesTransactionMonthlyMethodByMerchantFailed(rows []*models.TransactionMonthlyMethodByMerchantFailedRow) []*pbtransactions.TransactionMonthlyMethod {
	var transaction []*pbtransactions.TransactionMonthlyMethod
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionMonthMethodByMerchantFailed(row))
	}
	return transaction
}

func mapResponseTransactionYearlyMethodByMerchantFailed(row *models.TransactionYearlyMethodByMerchantFailedRow) *pbtransactions.TransactionYearlyMethod {
	if row == nil {
		return nil
	}
	return &pbtransactions.TransactionYearlyMethod{
		Year:              row.Year,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int32(row.TotalTransactions),
		TotalAmount:       int32(row.TotalAmount),
	}
}

func mapResponsesTransactionYearlyMethodByMerchantFailed(rows []*models.TransactionYearlyMethodByMerchantFailedRow) []*pbtransactions.TransactionYearlyMethod {
	var transaction []*pbtransactions.TransactionYearlyMethod
	for _, row := range rows {
		transaction = append(transaction, mapResponseTransactionYearlyMethodByMerchantFailed(row))
	}
	return transaction
}
