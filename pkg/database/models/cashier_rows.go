package models

import "time"

// CashierRow - paginated cashier listing without deleted_at
type CashierRow struct {
	CashierID  int32     `gorm:"column:cashier_id" json:"cashier_id"`
	MerchantID int32     `gorm:"column:merchant_id" json:"merchant_id"`
	UserID     int32     `gorm:"column:user_id" json:"user_id"`
	Name       string    `gorm:"column:name" json:"name"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
	TotalCount int64     `gorm:"column:total_count" json:"total_count"`
}

// CashierByMerchantRow - paginated cashier listing filtered by merchant without deleted_at
type CashierByMerchantRow struct {
	CashierID  int32     `gorm:"column:cashier_id" json:"cashier_id"`
	MerchantID int32     `gorm:"column:merchant_id" json:"merchant_id"`
	UserID     int32     `gorm:"column:user_id" json:"user_id"`
	Name       string    `gorm:"column:name" json:"name"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
	TotalCount int64     `gorm:"column:total_count" json:"total_count"`
}

// CashierActiveRow - paginated active cashier listing with deleted_at
type CashierActiveRow struct {
	CashierID  int32      `gorm:"column:cashier_id" json:"cashier_id"`
	MerchantID int32      `gorm:"column:merchant_id" json:"merchant_id"`
	UserID     int32      `gorm:"column:user_id" json:"user_id"`
	Name       string     `gorm:"column:name" json:"name"`
	CreatedAt  time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt  *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	TotalCount int64      `gorm:"column:total_count" json:"total_count"`
}

// CashierTrashedRow - paginated trashed cashier listing with deleted_at
type CashierTrashedRow struct {
	CashierID  int32      `gorm:"column:cashier_id" json:"cashier_id"`
	MerchantID int32      `gorm:"column:merchant_id" json:"merchant_id"`
	UserID     int32      `gorm:"column:user_id" json:"user_id"`
	Name       string     `gorm:"column:name" json:"name"`
	CreatedAt  time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt  *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	TotalCount int64      `gorm:"column:total_count" json:"total_count"`
}

// CashierMonthlyTotalSalesRow - monthly total sales aggregated across cashiers,
// optionally filtered by cashier or merchant
type CashierMonthlyTotalSalesRow struct {
	Year       string `gorm:"column:year" json:"year"`
	Month      string `gorm:"column:month" json:"month"`
	TotalSales int32  `gorm:"column:total_sales" json:"total_sales"`
}

// CashierYearlyTotalSalesRow - yearly total sales aggregated across cashiers,
// optionally filtered by cashier or merchant
type CashierYearlyTotalSalesRow struct {
	Year       string `gorm:"column:year" json:"year"`
	TotalSales int32  `gorm:"column:total_sales" json:"total_sales"`
}

// CashierMonthlySaleRow - monthly per-cashier sales activity (numeric total)
type CashierMonthlySaleRow struct {
	CashierID   int32   `gorm:"column:cashier_id" json:"cashier_id"`
	CashierName string  `gorm:"column:cashier_name" json:"cashier_name"`
	Month       string  `gorm:"column:month" json:"month"`
	OrderCount  int64   `gorm:"column:order_count" json:"order_count"`
	TotalSales  float64 `gorm:"column:total_sales" json:"total_sales"`
}

// CashierMonthlySaleByIdRow - monthly per-cashier sales activity filtered by cashier id
type CashierMonthlySaleByIdRow struct {
	CashierID   int32  `gorm:"column:cashier_id" json:"cashier_id"`
	CashierName string `gorm:"column:cashier_name" json:"cashier_name"`
	Month       string `gorm:"column:month" json:"month"`
	OrderCount  int64  `gorm:"column:order_count" json:"order_count"`
	TotalSales  int64  `gorm:"column:total_sales" json:"total_sales"`
}

// CashierMonthlySaleByMerchantRow - monthly per-cashier sales activity filtered by merchant id
type CashierMonthlySaleByMerchantRow struct {
	CashierID   int32  `gorm:"column:cashier_id" json:"cashier_id"`
	CashierName string `gorm:"column:cashier_name" json:"cashier_name"`
	Month       string `gorm:"column:month" json:"month"`
	OrderCount  int64  `gorm:"column:order_count" json:"order_count"`
	TotalSales  int64  `gorm:"column:total_sales" json:"total_sales"`
}

// CashierYearlySaleRow - yearly per-cashier sales performance over a 5-year span
type CashierYearlySaleRow struct {
	Year        string `gorm:"column:year" json:"year"`
	CashierID   int32  `gorm:"column:cashier_id" json:"cashier_id"`
	CashierName string `gorm:"column:cashier_name" json:"cashier_name"`
	OrderCount  int64  `gorm:"column:order_count" json:"order_count"`
	TotalSales  int64  `gorm:"column:total_sales" json:"total_sales"`
}
