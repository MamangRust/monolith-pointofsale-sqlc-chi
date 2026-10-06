package models

import "time"

// TransactionRow - paginated transaction listing (all/active) with total_count
type TransactionRow struct {
	TransactionID int32      `gorm:"column:transaction_id" json:"transaction_id"`
	OrderID       int32      `gorm:"column:order_id" json:"order_id"`
	MerchantID    int32      `gorm:"column:merchant_id" json:"merchant_id"`
	PaymentMethod string     `gorm:"column:payment_method" json:"payment_method"`
	Amount        int32      `gorm:"column:amount" json:"amount"`
	ChangeAmount  *int32     `gorm:"column:change_amount" json:"change_amount"`
	PaymentStatus string     `gorm:"column:payment_status" json:"payment_status"`
	CreatedAt     time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt     *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	TotalCount    int64      `gorm:"column:total_count" json:"total_count"`
}

// TransactionActiveRow - paginated active transaction listing
type TransactionActiveRow struct {
	TransactionID int32      `gorm:"column:transaction_id" json:"transaction_id"`
	OrderID       int32      `gorm:"column:order_id" json:"order_id"`
	MerchantID    int32      `gorm:"column:merchant_id" json:"merchant_id"`
	PaymentMethod string     `gorm:"column:payment_method" json:"payment_method"`
	Amount        int32      `gorm:"column:amount" json:"amount"`
	ChangeAmount  *int32     `gorm:"column:change_amount" json:"change_amount"`
	PaymentStatus string     `gorm:"column:payment_status" json:"payment_status"`
	CreatedAt     time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt     *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	TotalCount    int64      `gorm:"column:total_count" json:"total_count"`
}

// TransactionTrashedRow - paginated trashed transaction listing
type TransactionTrashedRow struct {
	TransactionID int32      `gorm:"column:transaction_id" json:"transaction_id"`
	OrderID       int32      `gorm:"column:order_id" json:"order_id"`
	MerchantID    int32      `gorm:"column:merchant_id" json:"merchant_id"`
	PaymentMethod string     `gorm:"column:payment_method" json:"payment_method"`
	Amount        int32      `gorm:"column:amount" json:"amount"`
	ChangeAmount  *int32     `gorm:"column:change_amount" json:"change_amount"`
	PaymentStatus string     `gorm:"column:payment_status" json:"payment_status"`
	CreatedAt     time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt     *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	TotalCount    int64      `gorm:"column:total_count" json:"total_count"`
}

// TransactionByMerchantRow - paginated transaction listing by merchant
type TransactionByMerchantRow struct {
	TransactionID int32      `gorm:"column:transaction_id" json:"transaction_id"`
	OrderID       int32      `gorm:"column:order_id" json:"order_id"`
	MerchantID    int32      `gorm:"column:merchant_id" json:"merchant_id"`
	PaymentMethod string     `gorm:"column:payment_method" json:"payment_method"`
	Amount        int32      `gorm:"column:amount" json:"amount"`
	ChangeAmount  *int32     `gorm:"column:change_amount" json:"change_amount"`
	PaymentStatus string     `gorm:"column:payment_status" json:"payment_status"`
	CreatedAt     time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt     *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	TotalCount    int64      `gorm:"column:total_count" json:"total_count"`
}

// TransactionMonthlyAmountSuccessRow - monthly successful transaction metrics
type TransactionMonthlyAmountSuccessRow struct {
	Year         string `gorm:"column:year" json:"year"`
	Month        string `gorm:"column:month" json:"month"`
	TotalSuccess int64  `gorm:"column:total_success" json:"total_success"`
	TotalAmount  int32  `gorm:"column:total_amount" json:"total_amount"`
}

// TransactionMonthlyAmountFailedRow - monthly failed transaction metrics
type TransactionMonthlyAmountFailedRow struct {
	Year        string `gorm:"column:year" json:"year"`
	Month       string `gorm:"column:month" json:"month"`
	TotalFailed int64  `gorm:"column:total_failed" json:"total_failed"`
	TotalAmount int32  `gorm:"column:total_amount" json:"total_amount"`
}

// TransactionYearlyAmountSuccessRow - yearly successful transaction metrics
type TransactionYearlyAmountSuccessRow struct {
	Year         string `gorm:"column:year" json:"year"`
	TotalSuccess int32  `gorm:"column:total_success" json:"total_success"`
	TotalAmount  int32  `gorm:"column:total_amount" json:"total_amount"`
}

// TransactionYearlyAmountFailedRow - yearly failed transaction metrics
type TransactionYearlyAmountFailedRow struct {
	Year        string `gorm:"column:year" json:"year"`
	TotalFailed int32  `gorm:"column:total_failed" json:"total_failed"`
	TotalAmount int32  `gorm:"column:total_amount" json:"total_amount"`
}

// TransactionMonthlyMethodSuccessRow - monthly payment method usage (success)
type TransactionMonthlyMethodSuccessRow struct {
	Month             string  `gorm:"column:month" json:"month"`
	PaymentMethod     string  `gorm:"column:payment_method" json:"payment_method"`
	TotalTransactions int64   `gorm:"column:total_transactions" json:"total_transactions"`
	TotalAmount       float64 `gorm:"column:total_amount" json:"total_amount"`
}

// TransactionMonthlyMethodFailedRow - monthly payment method usage (failed)
type TransactionMonthlyMethodFailedRow struct {
	Month             string  `gorm:"column:month" json:"month"`
	PaymentMethod     string  `gorm:"column:payment_method" json:"payment_method"`
	TotalTransactions int64   `gorm:"column:total_transactions" json:"total_transactions"`
	TotalAmount       float64 `gorm:"column:total_amount" json:"total_amount"`
}

// TransactionYearlyMethodSuccessRow - yearly payment method usage (success)
type TransactionYearlyMethodSuccessRow struct {
	Year              string  `gorm:"column:year" json:"year"`
	PaymentMethod     string  `gorm:"column:payment_method" json:"payment_method"`
	TotalTransactions int64   `gorm:"column:total_transactions" json:"total_transactions"`
	TotalAmount       float64 `gorm:"column:total_amount" json:"total_amount"`
}

// TransactionYearlyMethodFailedRow - yearly payment method usage (failed)
type TransactionYearlyMethodFailedRow struct {
	Year              string  `gorm:"column:year" json:"year"`
	PaymentMethod     string  `gorm:"column:payment_method" json:"payment_method"`
	TotalTransactions int64   `gorm:"column:total_transactions" json:"total_transactions"`
	TotalAmount       float64 `gorm:"column:total_amount" json:"total_amount"`
}

// TransactionMonthlyAmountSuccessByMerchantRow - monthly success metrics by merchant
type TransactionMonthlyAmountSuccessByMerchantRow struct {
	Year         string `gorm:"column:year" json:"year"`
	Month        string `gorm:"column:month" json:"month"`
	TotalSuccess int64  `gorm:"column:total_success" json:"total_success"`
	TotalAmount  int32  `gorm:"column:total_amount" json:"total_amount"`
}

// TransactionMonthlyAmountFailedByMerchantRow - monthly failed metrics by merchant
type TransactionMonthlyAmountFailedByMerchantRow struct {
	Year        string `gorm:"column:year" json:"year"`
	Month       string `gorm:"column:month" json:"month"`
	TotalFailed int64  `gorm:"column:total_failed" json:"total_failed"`
	TotalAmount int32  `gorm:"column:total_amount" json:"total_amount"`
}

// TransactionYearlyAmountSuccessByMerchantRow - yearly success metrics by merchant
type TransactionYearlyAmountSuccessByMerchantRow struct {
	Year         string `gorm:"column:year" json:"year"`
	TotalSuccess int32  `gorm:"column:total_success" json:"total_success"`
	TotalAmount  int32  `gorm:"column:total_amount" json:"total_amount"`
}

// TransactionYearlyAmountFailedByMerchantRow - yearly failed metrics by merchant
type TransactionYearlyAmountFailedByMerchantRow struct {
	Year        string `gorm:"column:year" json:"year"`
	TotalFailed int32  `gorm:"column:total_failed" json:"total_failed"`
	TotalAmount int32  `gorm:"column:total_amount" json:"total_amount"`
}

// TransactionMonthlyMethodByMerchantSuccessRow - monthly method usage by merchant (success)
type TransactionMonthlyMethodByMerchantSuccessRow struct {
	Month             string  `gorm:"column:month" json:"month"`
	PaymentMethod     string  `gorm:"column:payment_method" json:"payment_method"`
	TotalTransactions int64   `gorm:"column:total_transactions" json:"total_transactions"`
	TotalAmount       float64 `gorm:"column:total_amount" json:"total_amount"`
}

// TransactionMonthlyMethodByMerchantFailedRow - monthly method usage by merchant (failed)
type TransactionMonthlyMethodByMerchantFailedRow struct {
	Month             string  `gorm:"column:month" json:"month"`
	PaymentMethod     string  `gorm:"column:payment_method" json:"payment_method"`
	TotalTransactions int64   `gorm:"column:total_transactions" json:"total_transactions"`
	TotalAmount       float64 `gorm:"column:total_amount" json:"total_amount"`
}

// TransactionYearlyMethodByMerchantSuccessRow - yearly method usage by merchant (success)
type TransactionYearlyMethodByMerchantSuccessRow struct {
	Year              string  `gorm:"column:year" json:"year"`
	PaymentMethod     string  `gorm:"column:payment_method" json:"payment_method"`
	TotalTransactions int64   `gorm:"column:total_transactions" json:"total_transactions"`
	TotalAmount       float64 `gorm:"column:total_amount" json:"total_amount"`
}

// TransactionYearlyMethodByMerchantFailedRow - yearly method usage by merchant (failed)
type TransactionYearlyMethodByMerchantFailedRow struct {
	Year              string  `gorm:"column:year" json:"year"`
	PaymentMethod     string  `gorm:"column:payment_method" json:"payment_method"`
	TotalTransactions int64   `gorm:"column:total_transactions" json:"total_transactions"`
	TotalAmount       float64 `gorm:"column:total_amount" json:"total_amount"`
}
