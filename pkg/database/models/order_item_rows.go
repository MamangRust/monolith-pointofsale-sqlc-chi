package models

import "time"

// OrderItemRow - paginated order item listing
type OrderItemRow struct {
	OrderItemID int32     `gorm:"column:order_item_id" json:"order_item_id"`
	OrderID     int32     `gorm:"column:order_id" json:"order_id"`
	ProductID   int32     `gorm:"column:product_id" json:"product_id"`
	Quantity    int32     `gorm:"column:quantity" json:"quantity"`
	Price       int32     `gorm:"column:price" json:"price"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
	TotalCount  int64     `gorm:"column:total_count" json:"total_count"`
}

// OrderItemActiveRow - paginated active (non-deleted) order item listing
type OrderItemActiveRow struct {
	OrderItemID int32      `gorm:"column:order_item_id" json:"order_item_id"`
	OrderID     int32      `gorm:"column:order_id" json:"order_id"`
	ProductID   int32      `gorm:"column:product_id" json:"product_id"`
	Quantity    int32      `gorm:"column:quantity" json:"quantity"`
	Price       int32      `gorm:"column:price" json:"price"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	TotalCount  int64      `gorm:"column:total_count" json:"total_count"`
}

// OrderItemTrashedRow - paginated trashed order item listing
type OrderItemTrashedRow struct {
	OrderItemID int32      `gorm:"column:order_item_id" json:"order_item_id"`
	OrderID     int32      `gorm:"column:order_id" json:"order_id"`
	ProductID   int32      `gorm:"column:product_id" json:"product_id"`
	Quantity    int32      `gorm:"column:quantity" json:"quantity"`
	Price       int32      `gorm:"column:price" json:"price"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	TotalCount  int64      `gorm:"column:total_count" json:"total_count"`
}
