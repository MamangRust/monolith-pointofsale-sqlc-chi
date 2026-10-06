package models

import "time"

// OrderRow - paginated active order listing without deleted_at
type OrderRow struct {
	OrderID    int32     `gorm:"column:order_id" json:"order_id"`
	MerchantID int32     `gorm:"column:merchant_id" json:"merchant_id"`
	CashierID  int32     `gorm:"column:cashier_id" json:"cashier_id"`
	TotalPrice int64     `gorm:"column:total_price" json:"total_price"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
	TotalCount int64     `gorm:"column:total_count" json:"total_count"`
}

// OrderActiveRow - paginated active order listing with deleted_at
type OrderActiveRow struct {
	OrderID    int32      `gorm:"column:order_id" json:"order_id"`
	MerchantID int32      `gorm:"column:merchant_id" json:"merchant_id"`
	CashierID  int32      `gorm:"column:cashier_id" json:"cashier_id"`
	TotalPrice int64      `gorm:"column:total_price" json:"total_price"`
	CreatedAt  time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt  *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	TotalCount int64      `gorm:"column:total_count" json:"total_count"`
}

// OrderTrashedRow - paginated trashed order listing with deleted_at
type OrderTrashedRow struct {
	OrderID    int32      `gorm:"column:order_id" json:"order_id"`
	MerchantID int32      `gorm:"column:merchant_id" json:"merchant_id"`
	CashierID  int32      `gorm:"column:cashier_id" json:"cashier_id"`
	TotalPrice int64      `gorm:"column:total_price" json:"total_price"`
	CreatedAt  time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt  *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	TotalCount int64      `gorm:"column:total_count" json:"total_count"`
}

// OrderMerchantRow - paginated order listing for a specific merchant
type OrderMerchantRow struct {
	OrderID    int32     `gorm:"column:order_id" json:"order_id"`
	MerchantID int32     `gorm:"column:merchant_id" json:"merchant_id"`
	CashierID  int32     `gorm:"column:cashier_id" json:"cashier_id"`
	TotalPrice int64     `gorm:"column:total_price" json:"total_price"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
	TotalCount int64     `gorm:"column:total_count" json:"total_count"`
}

// MonthlyTotalRevenueRow - monthly revenue summary row
type MonthlyTotalRevenueRow struct {
	Year         string `gorm:"column:year" json:"year"`
	Month        string `gorm:"column:month" json:"month"`
	TotalRevenue int32  `gorm:"column:total_revenue" json:"total_revenue"`
}

// YearlyTotalRevenueRow - yearly revenue summary row
type YearlyTotalRevenueRow struct {
	Year         string `gorm:"column:year" json:"year"`
	TotalRevenue int32  `gorm:"column:total_revenue" json:"total_revenue"`
}

// MonthlyOrderRow - monthly order summary row
type MonthlyOrderRow struct {
	Month          string  `gorm:"column:month" json:"month"`
	OrderCount     int64   `gorm:"column:order_count" json:"order_count"`
	TotalRevenue   float64 `gorm:"column:total_revenue" json:"total_revenue"`
	TotalItemsSold int64   `gorm:"column:total_items_sold" json:"total_items_sold"`
}

// YearlyOrderRow - yearly order summary row
type YearlyOrderRow struct {
	Year               string  `gorm:"column:year" json:"year"`
	OrderCount         int64   `gorm:"column:order_count" json:"order_count"`
	TotalRevenue       float64 `gorm:"column:total_revenue" json:"total_revenue"`
	TotalItemsSold     int64   `gorm:"column:total_items_sold" json:"total_items_sold"`
	ActiveCashiers     int64   `gorm:"column:active_cashiers" json:"active_cashiers"`
	UniqueProductsSold int64   `gorm:"column:unique_products_sold" json:"unique_products_sold"`
}

// MonthlyTotalRevenueByMerchantRow - monthly revenue summary row by merchant
type MonthlyTotalRevenueByMerchantRow struct {
	Year         string `gorm:"column:year" json:"year"`
	Month        string `gorm:"column:month" json:"month"`
	TotalRevenue int32  `gorm:"column:total_revenue" json:"total_revenue"`
}

// YearlyTotalRevenueByMerchantRow - yearly revenue summary row by merchant
type YearlyTotalRevenueByMerchantRow struct {
	Year         string `gorm:"column:year" json:"year"`
	TotalRevenue int32  `gorm:"column:total_revenue" json:"total_revenue"`
}

// MonthlyOrderByMerchantRow - monthly order summary row by merchant
type MonthlyOrderByMerchantRow struct {
	Month          string  `gorm:"column:month" json:"month"`
	OrderCount     int64   `gorm:"column:order_count" json:"order_count"`
	TotalRevenue   float64 `gorm:"column:total_revenue" json:"total_revenue"`
	TotalItemsSold int64   `gorm:"column:total_items_sold" json:"total_items_sold"`
}

// YearlyOrderByMerchantRow - yearly order summary row by merchant
type YearlyOrderByMerchantRow struct {
	Year               string  `gorm:"column:year" json:"year"`
	OrderCount         int64   `gorm:"column:order_count" json:"order_count"`
	TotalRevenue       float64 `gorm:"column:total_revenue" json:"total_revenue"`
	TotalItemsSold     int64   `gorm:"column:total_items_sold" json:"total_items_sold"`
	ActiveCashiers     int64   `gorm:"column:active_cashiers" json:"active_cashiers"`
	UniqueProductsSold int64   `gorm:"column:unique_products_sold" json:"unique_products_sold"`
}

// OrderByMerchantRow - order listing filtered by merchant
type OrderByMerchantRow struct {
	OrderID    int32     `gorm:"column:order_id" json:"order_id"`
	MerchantID int32     `gorm:"column:merchant_id" json:"merchant_id"`
	CashierID  int32     `gorm:"column:cashier_id" json:"cashier_id"`
	TotalPrice int64     `gorm:"column:total_price" json:"total_price"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
	TotalCount int64     `gorm:"column:total_count" json:"total_count"`
}
