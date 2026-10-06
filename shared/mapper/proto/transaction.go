package protomapper

import (
	pbcommon "github.com/MamangRust/monolith-point-of-sale-pb/common"
	pbtransactions "github.com/MamangRust/monolith-point-of-sale-pb/transactions"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

type transactionProtoMapper struct{}

func NewTransactionProtoMapper() *transactionProtoMapper {
	return &transactionProtoMapper{}
}

func (t *transactionProtoMapper) ToProtoResponseTransaction(status string, message string, trans *response.TransactionResponse) *pbtransactions.ApiResponseTransaction {
	return &pbtransactions.ApiResponseTransaction{
		Status:  status,
		Message: message,
		Data:    t.mapResponseTransaction(trans),
	}
}

func (t *transactionProtoMapper) ToProtoResponsesTransaction(status string, message string, transList []*response.TransactionResponse) *pbtransactions.ApiResponsesTransaction {
	return &pbtransactions.ApiResponsesTransaction{
		Status:  status,
		Message: message,
		Data:    t.mapResponsesTransaction(transList),
	}
}

func (t *transactionProtoMapper) ToProtoResponseTransactionDeleteAt(status string, message string, trans *response.TransactionResponseDeleteAt) *pbtransactions.ApiResponseTransactionDeleteAt {
	return &pbtransactions.ApiResponseTransactionDeleteAt{
		Status:  status,
		Message: message,
		Data:    t.mapResponseTransactionDeleteAt(trans),
	}
}

func (t *transactionProtoMapper) ToProtoResponseTransactionDelete(status string, message string) *pbtransactions.ApiResponseTransactionDelete {
	return &pbtransactions.ApiResponseTransactionDelete{
		Status:  status,
		Message: message,
	}
}

func (t *transactionProtoMapper) ToProtoResponseTransactionAll(status string, message string) *pbtransactions.ApiResponseTransactionAll {
	return &pbtransactions.ApiResponseTransactionAll{
		Status:  status,
		Message: message,
	}
}

func (t *transactionProtoMapper) ToProtoResponsePaginationTransactionDeleteAt(pagination *pbcommon.PaginationMeta, status string, message string, transactions []*response.TransactionResponseDeleteAt) *pbtransactions.ApiResponsePaginationTransactionDeleteAt {
	return &pbtransactions.ApiResponsePaginationTransactionDeleteAt{
		Status:     status,
		Message:    message,
		Data:       t.mapResponsesTransactionDeleteAt(transactions),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (t *transactionProtoMapper) ToProtoResponsePaginationTransaction(pagination *pbcommon.PaginationMeta, status string, message string, transactions []*response.TransactionResponse) *pbtransactions.ApiResponsePaginationTransaction {
	return &pbtransactions.ApiResponsePaginationTransaction{
		Status:     status,
		Message:    message,
		Data:       t.mapResponsesTransaction(transactions),
		Pagination: mapPaginationMeta(pagination),
	}
}

func (t *transactionProtoMapper) ToProtoResponseMonthAmountSuccess(status string, message string, row []*response.TransactionMonthlyAmountSuccessResponse) *pbtransactions.ApiResponseTransactionMonthAmountSuccess {
	return &pbtransactions.ApiResponseTransactionMonthAmountSuccess{
		Status:  status,
		Message: message,
		Data:    t.mapResponsesTransactionMonthlyAmountSuccess(row),
	}
}

func (t *transactionProtoMapper) ToProtoResponseYearAmountSuccess(status string, message string, row []*response.TransactionYearlyAmountSuccessResponse) *pbtransactions.ApiResponseTransactionYearAmountSuccess {
	return &pbtransactions.ApiResponseTransactionYearAmountSuccess{
		Status:  status,
		Message: message,
		Data:    t.mapResponsesTransactionYearlyAmountSuccess(row),
	}
}

func (t *transactionProtoMapper) ToProtoResponseMonthAmountFailed(status string, message string, row []*response.TransactionMonthlyAmountFailedResponse) *pbtransactions.ApiResponseTransactionMonthAmountFailed {
	return &pbtransactions.ApiResponseTransactionMonthAmountFailed{
		Status:  status,
		Message: message,
		Data:    t.mapResponsesTransactionMonthlyAmountFailed(row),
	}
}

func (t *transactionProtoMapper) ToProtoResponseYearAmountFailed(status string, message string, row []*response.TransactionYearlyAmountFailedResponse) *pbtransactions.ApiResponseTransactionYearAmountFailed {
	return &pbtransactions.ApiResponseTransactionYearAmountFailed{
		Status:  status,
		Message: message,
		Data:    t.mapResponsesTransactionYearlyAmountFailed(row),
	}
}

func (t *transactionProtoMapper) ToProtoResponseMonthMethod(status string, message string, row []*response.TransactionMonthlyMethodResponse) *pbtransactions.ApiResponseTransactionMonthPaymentMethod {
	return &pbtransactions.ApiResponseTransactionMonthPaymentMethod{
		Status:  status,
		Message: message,
		Data:    t.mapResponsesTransactionMonthlyMethod(row),
	}
}

func (t *transactionProtoMapper) ToProtoResponseYearMethod(status string, message string, row []*response.TransactionYearlyMethodResponse) *pbtransactions.ApiResponseTransactionYearPaymentmethod {
	return &pbtransactions.ApiResponseTransactionYearPaymentmethod{
		Status:  status,
		Message: message,
		Data:    t.mapResponsesTransactionYearlyMethod(row),
	}
}

func (t *transactionProtoMapper) mapResponseTransaction(transaction *response.TransactionResponse) *pbtransactions.TransactionResponse {
	return &pbtransactions.TransactionResponse{
		Id:            int32(transaction.ID),
		OrderId:       int32(transaction.OrderID),
		MerchantId:    int32(transaction.MerchantID),
		PaymentMethod: transaction.PaymentMethod,
		Amount:        int32(transaction.Amount),
		ChangeAmount:  int32(transaction.ChangeAmount),
		PaymentStatus: transaction.PaymentStatus,
		CreatedAt:     transaction.CreatedAt,
		UpdatedAt:     transaction.UpdatedAt,
	}
}

func (t *transactionProtoMapper) mapResponsesTransaction(transactions []*response.TransactionResponse) []*pbtransactions.TransactionResponse {
	var mappedTransactions []*pbtransactions.TransactionResponse

	for _, transaction := range transactions {
		mappedTransactions = append(mappedTransactions, t.mapResponseTransaction(transaction))
	}

	return mappedTransactions
}

func (t *transactionProtoMapper) mapResponseTransactionDeleteAt(transaction *response.TransactionResponseDeleteAt) *pbtransactions.TransactionResponseDeleteAt {
	var deletedAt *wrapperspb.StringValue
	if transaction.DeletedAt != nil {
		deletedAt = wrapperspb.String(*transaction.DeletedAt)
	}

	return &pbtransactions.TransactionResponseDeleteAt{
		Id:            int32(transaction.ID),
		OrderId:       int32(transaction.OrderID),
		MerchantId:    int32(transaction.MerchantID),
		PaymentMethod: transaction.PaymentMethod,
		Amount:        int32(transaction.Amount),
		ChangeAmount:  int32(transaction.ChangeAmount),
		PaymentStatus: transaction.PaymentStatus,
		CreatedAt:     transaction.CreatedAt,
		UpdatedAt:     transaction.UpdatedAt,
		DeletedAt:     deletedAt,
	}
}

func (t *transactionProtoMapper) mapResponsesTransactionDeleteAt(transactions []*response.TransactionResponseDeleteAt) []*pbtransactions.TransactionResponseDeleteAt {
	var mappedTransactions []*pbtransactions.TransactionResponseDeleteAt

	for _, transaction := range transactions {
		mappedTransactions = append(mappedTransactions, t.mapResponseTransactionDeleteAt(transaction))
	}

	return mappedTransactions
}

func (s *transactionProtoMapper) mapResponseTransactionMonthAmountSuccess(row *response.TransactionMonthlyAmountSuccessResponse) *pbtransactions.TransactionMonthlyAmountSuccess {
	return &pbtransactions.TransactionMonthlyAmountSuccess{
		Year:         row.Year,
		Month:        row.Month,
		TotalSuccess: int32(row.TotalSuccess),
		TotalAmount:  int32(row.TotalAmount),
	}
}

func (s *transactionProtoMapper) mapResponsesTransactionMonthlyAmountSuccess(rows []*response.TransactionMonthlyAmountSuccessResponse) []*pbtransactions.TransactionMonthlyAmountSuccess {
	var transaction []*pbtransactions.TransactionMonthlyAmountSuccess

	for _, row := range rows {
		transaction = append(transaction, s.mapResponseTransactionMonthAmountSuccess(row))
	}

	return transaction
}

func (s *transactionProtoMapper) mapResponseTransactionYearAmountSuccess(row *response.TransactionYearlyAmountSuccessResponse) *pbtransactions.TransactionYearlyAmountSuccess {
	return &pbtransactions.TransactionYearlyAmountSuccess{
		Year:         row.Year,
		TotalSuccess: int32(row.TotalSuccess),
		TotalAmount:  int32(row.TotalAmount),
	}
}

func (s *transactionProtoMapper) mapResponsesTransactionYearlyAmountSuccess(rows []*response.TransactionYearlyAmountSuccessResponse) []*pbtransactions.TransactionYearlyAmountSuccess {
	var transaction []*pbtransactions.TransactionYearlyAmountSuccess

	for _, row := range rows {
		transaction = append(transaction, s.mapResponseTransactionYearAmountSuccess(row))
	}

	return transaction
}

func (s *transactionProtoMapper) mapResponseTransactionMonthAmountFailed(row *response.TransactionMonthlyAmountFailedResponse) *pbtransactions.TransactionMonthlyAmountFailed {
	return &pbtransactions.TransactionMonthlyAmountFailed{
		Year:        row.Year,
		Month:       row.Month,
		TotalFailed: int32(row.TotalFailed),
		TotalAmount: int32(row.TotalAmount),
	}
}

func (s *transactionProtoMapper) mapResponsesTransactionMonthlyAmountFailed(rows []*response.TransactionMonthlyAmountFailedResponse) []*pbtransactions.TransactionMonthlyAmountFailed {
	var transaction []*pbtransactions.TransactionMonthlyAmountFailed

	for _, row := range rows {
		transaction = append(transaction, s.mapResponseTransactionMonthAmountFailed(row))
	}

	return transaction
}

func (s *transactionProtoMapper) mapResponseTransactionYearAmountFailed(row *response.TransactionYearlyAmountFailedResponse) *pbtransactions.TransactionYearlyAmountFailed {
	return &pbtransactions.TransactionYearlyAmountFailed{
		Year:        row.Year,
		TotalFailed: int32(row.TotalFailed),
		TotalAmount: int32(row.TotalAmount),
	}
}

func (s *transactionProtoMapper) mapResponsesTransactionYearlyAmountFailed(rows []*response.TransactionYearlyAmountFailedResponse) []*pbtransactions.TransactionYearlyAmountFailed {
	var transaction []*pbtransactions.TransactionYearlyAmountFailed

	for _, row := range rows {
		transaction = append(transaction, s.mapResponseTransactionYearAmountFailed(row))
	}

	return transaction
}

func (s *transactionProtoMapper) mapResponseTransactionMonthMethod(row *response.TransactionMonthlyMethodResponse) *pbtransactions.TransactionMonthlyMethod {
	return &pbtransactions.TransactionMonthlyMethod{
		Month:             row.Month,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int32(row.TotalTransactions),
		TotalAmount:       int32(row.TotalAmount),
	}
}

func (s *transactionProtoMapper) mapResponsesTransactionMonthlyMethod(rows []*response.TransactionMonthlyMethodResponse) []*pbtransactions.TransactionMonthlyMethod {
	var transaction []*pbtransactions.TransactionMonthlyMethod

	for _, row := range rows {
		transaction = append(transaction, s.mapResponseTransactionMonthMethod(row))
	}

	return transaction
}

func (s *transactionProtoMapper) mapResponseTransactionYearMethod(row *response.TransactionYearlyMethodResponse) *pbtransactions.TransactionYearlyMethod {
	return &pbtransactions.TransactionYearlyMethod{
		Year:              row.Year,
		PaymentMethod:     row.PaymentMethod,
		TotalTransactions: int32(row.TotalTransactions),
		TotalAmount:       int32(row.TotalAmount),
	}
}

func (s *transactionProtoMapper) mapResponsesTransactionYearlyMethod(rows []*response.TransactionYearlyMethodResponse) []*pbtransactions.TransactionYearlyMethod {
	var transaction []*pbtransactions.TransactionYearlyMethod

	for _, row := range rows {
		transaction = append(transaction, s.mapResponseTransactionYearMethod(row))
	}

	return transaction
}
