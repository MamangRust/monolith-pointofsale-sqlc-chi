package models

import "time"

// ProductRow - paginated product listing with total_count
type ProductRow struct {
	ProductID    int32      `gorm:"column:product_id" json:"product_id"`
	MerchantID   int32      `gorm:"column:merchant_id" json:"merchant_id"`
	CategoryID   int32      `gorm:"column:category_id" json:"category_id"`
	Name         string     `gorm:"column:name" json:"name"`
	Description  *string    `gorm:"column:description" json:"description"`
	Price        int32      `gorm:"column:price" json:"price"`
	CountInStock int32      `gorm:"column:count_in_stock" json:"count_in_stock"`
	Brand        *string    `gorm:"column:brand" json:"brand"`
	Weight       *int32     `gorm:"column:weight" json:"weight"`
	SlugProduct  *string    `gorm:"column:slug_product" json:"slug_product"`
	ImageProduct *string    `gorm:"column:image_product" json:"image_product"`
	Barcode      *string    `gorm:"column:barcode" json:"barcode"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	TotalCount   int64      `gorm:"column:total_count" json:"total_count"`
}

// ProductActiveRow - paginated active (non-deleted) product listing
type ProductActiveRow struct {
	ProductID    int32      `gorm:"column:product_id" json:"product_id"`
	MerchantID   int32      `gorm:"column:merchant_id" json:"merchant_id"`
	CategoryID   int32      `gorm:"column:category_id" json:"category_id"`
	Name         string     `gorm:"column:name" json:"name"`
	Description  *string    `gorm:"column:description" json:"description"`
	Price        int32      `gorm:"column:price" json:"price"`
	CountInStock int32      `gorm:"column:count_in_stock" json:"count_in_stock"`
	Brand        *string    `gorm:"column:brand" json:"brand"`
	Weight       *int32     `gorm:"column:weight" json:"weight"`
	SlugProduct  *string    `gorm:"column:slug_product" json:"slug_product"`
	ImageProduct *string    `gorm:"column:image_product" json:"image_product"`
	Barcode      *string    `gorm:"column:barcode" json:"barcode"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	TotalCount   int64      `gorm:"column:total_count" json:"total_count"`
}

// ProductTrashedRow - paginated trashed (soft-deleted) product listing
type ProductTrashedRow struct {
	ProductID    int32      `gorm:"column:product_id" json:"product_id"`
	MerchantID   int32      `gorm:"column:merchant_id" json:"merchant_id"`
	CategoryID   int32      `gorm:"column:category_id" json:"category_id"`
	Name         string     `gorm:"column:name" json:"name"`
	Description  *string    `gorm:"column:description" json:"description"`
	Price        int32      `gorm:"column:price" json:"price"`
	CountInStock int32      `gorm:"column:count_in_stock" json:"count_in_stock"`
	Brand        *string    `gorm:"column:brand" json:"brand"`
	Weight       *int32     `gorm:"column:weight" json:"weight"`
	SlugProduct  *string    `gorm:"column:slug_product" json:"slug_product"`
	ImageProduct *string    `gorm:"column:image_product" json:"image_product"`
	Barcode      *string    `gorm:"column:barcode" json:"barcode"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	TotalCount   int64      `gorm:"column:total_count" json:"total_count"`
}

// ProductMerchantRow - paginated products owned by a merchant, joined with category name
type ProductMerchantRow struct {
	TotalCount   int64     `gorm:"column:total_count" json:"total_count"`
	ProductID    int32     `gorm:"column:product_id" json:"product_id"`
	CategoryID   int32     `gorm:"column:category_id" json:"category_id"`
	Name         string    `gorm:"column:name" json:"name"`
	Description  *string   `gorm:"column:description" json:"description"`
	Price        int32     `gorm:"column:price" json:"price"`
	CountInStock int32     `gorm:"column:count_in_stock" json:"count_in_stock"`
	Brand        *string   `gorm:"column:brand" json:"brand"`
	ImageProduct *string   `gorm:"column:image_product" json:"image_product"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	CategoryName string    `gorm:"column:category_name" json:"category_name"`
}

// ProductCategoryNameRow - paginated products under a category name, joined with category name
type ProductCategoryNameRow struct {
	TotalCount   int64      `gorm:"column:total_count" json:"total_count"`
	ProductID    int32      `gorm:"column:product_id" json:"product_id"`
	MerchantID   int32      `gorm:"column:merchant_id" json:"merchant_id"`
	CategoryID   int32      `gorm:"column:category_id" json:"category_id"`
	SlugProduct  *string    `gorm:"column:slug_product" json:"slug_product"`
	Weight       *int32     `gorm:"column:weight" json:"weight"`
	Name         string     `gorm:"column:name" json:"name"`
	Description  *string    `gorm:"column:description" json:"description"`
	Price        int32      `gorm:"column:price" json:"price"`
	CountInStock int32      `gorm:"column:count_in_stock" json:"count_in_stock"`
	Brand        *string    `gorm:"column:brand" json:"brand"`
	ImageProduct *string    `gorm:"column:image_product" json:"image_product"`
	Barcode      *string    `gorm:"column:barcode" json:"barcode"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	CategoryName string     `gorm:"column:category_name" json:"category_name"`
}
