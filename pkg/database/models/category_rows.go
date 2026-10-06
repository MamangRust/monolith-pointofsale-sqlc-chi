package models

import "time"

// CategoryRow - paginated category listing without deleted_at
type CategoryRow struct {
	CategoryID   int32     `gorm:"column:category_id" json:"category_id"`
	Name         string    `gorm:"column:name" json:"name"`
	Description  *string   `gorm:"column:description" json:"description"`
	SlugCategory *string   `gorm:"column:slug_category" json:"slug_category"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
	TotalCount   int64     `gorm:"column:total_count" json:"total_count"`
}

// CategoryActiveRow - paginated active category listing with deleted_at
type CategoryActiveRow struct {
	CategoryID   int32      `gorm:"column:category_id" json:"category_id"`
	Name         string     `gorm:"column:name" json:"name"`
	Description  *string    `gorm:"column:description" json:"description"`
	SlugCategory *string    `gorm:"column:slug_category" json:"slug_category"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	TotalCount   int64      `gorm:"column:total_count" json:"total_count"`
}

// CategoryTrashedRow - paginated trashed category listing with deleted_at
type CategoryTrashedRow struct {
	CategoryID   int32      `gorm:"column:category_id" json:"category_id"`
	Name         string     `gorm:"column:name" json:"name"`
	Description  *string    `gorm:"column:description" json:"description"`
	SlugCategory *string    `gorm:"column:slug_category" json:"slug_category"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	TotalCount   int64      `gorm:"column:total_count" json:"total_count"`
}

// CategoryMonthlyTotalPriceRow - monthly revenue totals across two comparison periods
type CategoryMonthlyTotalPriceRow struct {
	Year         string `gorm:"column:year" json:"year"`
	Month        string `gorm:"column:month" json:"month"`
	TotalRevenue int32  `gorm:"column:total_revenue" json:"total_revenue"`
}

// CategoryMonthlyTotalPriceByIdRow - monthly revenue totals by category_id
type CategoryMonthlyTotalPriceByIdRow struct {
	Year         string `gorm:"column:year" json:"year"`
	Month        string `gorm:"column:month" json:"month"`
	TotalRevenue int32  `gorm:"column:total_revenue" json:"total_revenue"`
}

// CategoryMonthlyTotalPriceByMerchantRow - monthly revenue totals by merchant_id
type CategoryMonthlyTotalPriceByMerchantRow struct {
	Year         string `gorm:"column:year" json:"year"`
	Month        string `gorm:"column:month" json:"month"`
	TotalRevenue int32  `gorm:"column:total_revenue" json:"total_revenue"`
}

// CategoryYearlyTotalPriceRow - yearly revenue totals
type CategoryYearlyTotalPriceRow struct {
	Year         string `gorm:"column:year" json:"year"`
	TotalRevenue int32  `gorm:"column:total_revenue" json:"total_revenue"`
}

// CategoryYearlyTotalPriceByIdRow - yearly revenue totals by category_id
type CategoryYearlyTotalPriceByIdRow struct {
	Year         string `gorm:"column:year" json:"year"`
	TotalRevenue int32  `gorm:"column:total_revenue" json:"total_revenue"`
}

// CategoryYearlyTotalPriceByMerchantRow - yearly revenue totals by merchant_id
type CategoryYearlyTotalPriceByMerchantRow struct {
	Year         string `gorm:"column:year" json:"year"`
	TotalRevenue int32  `gorm:"column:total_revenue" json:"total_revenue"`
}

// CategoryMonthlyRow - monthly sales activity for all categories
type CategoryMonthlyRow struct {
	Month        string `gorm:"column:month" json:"month"`
	CategoryID   int32  `gorm:"column:category_id" json:"category_id"`
	CategoryName string `gorm:"column:category_name" json:"category_name"`
	OrderCount   int64  `gorm:"column:order_count" json:"order_count"`
	ItemsSold    int64  `gorm:"column:items_sold" json:"items_sold"`
	TotalRevenue int32  `gorm:"column:total_revenue" json:"total_revenue"`
}

// CategoryMonthlyByIdRow - monthly sales activity by category_id
type CategoryMonthlyByIdRow struct {
	Month        string `gorm:"column:month" json:"month"`
	CategoryID   int32  `gorm:"column:category_id" json:"category_id"`
	CategoryName string `gorm:"column:category_name" json:"category_name"`
	OrderCount   int64  `gorm:"column:order_count" json:"order_count"`
	ItemsSold    int64  `gorm:"column:items_sold" json:"items_sold"`
	TotalRevenue int32  `gorm:"column:total_revenue" json:"total_revenue"`
}

// CategoryMonthlyByMerchantRow - monthly sales activity by merchant_id
type CategoryMonthlyByMerchantRow struct {
	Month        string `gorm:"column:month" json:"month"`
	CategoryID   int32  `gorm:"column:category_id" json:"category_id"`
	CategoryName string `gorm:"column:category_name" json:"category_name"`
	OrderCount   int64  `gorm:"column:order_count" json:"order_count"`
	ItemsSold    int64  `gorm:"column:items_sold" json:"items_sold"`
	TotalRevenue int32  `gorm:"column:total_revenue" json:"total_revenue"`
}

// CategoryYearlyRow - yearly sales performance for categories over a 5-year span
type CategoryYearlyRow struct {
	Year               string `gorm:"column:year" json:"year"`
	CategoryID         int32  `gorm:"column:category_id" json:"category_id"`
	CategoryName       string `gorm:"column:category_name" json:"category_name"`
	OrderCount         int64  `gorm:"column:order_count" json:"order_count"`
	ItemsSold          int64  `gorm:"column:items_sold" json:"items_sold"`
	TotalRevenue       int32  `gorm:"column:total_revenue" json:"total_revenue"`
	UniqueProductsSold int64  `gorm:"column:unique_products_sold" json:"unique_products_sold"`
}

// CategoryYearlyByIdRow - yearly sales performance by category_id
type CategoryYearlyByIdRow struct {
	Year               string `gorm:"column:year" json:"year"`
	CategoryID         int32  `gorm:"column:category_id" json:"category_id"`
	CategoryName       string `gorm:"column:category_name" json:"category_name"`
	OrderCount         int64  `gorm:"column:order_count" json:"order_count"`
	ItemsSold          int64  `gorm:"column:items_sold" json:"items_sold"`
	TotalRevenue       int32  `gorm:"column:total_revenue" json:"total_revenue"`
	UniqueProductsSold int64  `gorm:"column:unique_products_sold" json:"unique_products_sold"`
}

// CategoryYearlyByMerchantRow - yearly sales performance by merchant_id
type CategoryYearlyByMerchantRow struct {
	Year               string `gorm:"column:year" json:"year"`
	CategoryID         int32  `gorm:"column:category_id" json:"category_id"`
	CategoryName       string `gorm:"column:category_name" json:"category_name"`
	OrderCount         int64  `gorm:"column:order_count" json:"order_count"`
	ItemsSold          int64  `gorm:"column:items_sold" json:"items_sold"`
	TotalRevenue       int32  `gorm:"column:total_revenue" json:"total_revenue"`
	UniqueProductsSold int64  `gorm:"column:unique_products_sold" json:"unique_products_sold"`
}
