package models

import "time"

// MerchantRow - paginated merchant listing without deleted_at
type MerchantRow struct {
	MerchantID   int32     `gorm:"column:merchant_id" json:"merchant_id"`
	UserID       int32     `gorm:"column:user_id" json:"user_id"`
	Name         string    `gorm:"column:name" json:"name"`
	Description  *string   `gorm:"column:description" json:"description"`
	Address      *string   `gorm:"column:address" json:"address"`
	ContactEmail *string   `gorm:"column:contact_email" json:"contact_email"`
	ContactPhone *string   `gorm:"column:contact_phone" json:"contact_phone"`
	Status       string    `gorm:"column:status" json:"status"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
	TotalCount   int64     `gorm:"column:total_count" json:"total_count"`
}

// MerchantActiveRow - paginated active merchant listing with deleted_at
type MerchantActiveRow struct {
	MerchantID   int32      `gorm:"column:merchant_id" json:"merchant_id"`
	UserID       int32      `gorm:"column:user_id" json:"user_id"`
	Name         string     `gorm:"column:name" json:"name"`
	Description  *string    `gorm:"column:description" json:"description"`
	Address      *string    `gorm:"column:address" json:"address"`
	ContactEmail *string    `gorm:"column:contact_email" json:"contact_email"`
	ContactPhone *string    `gorm:"column:contact_phone" json:"contact_phone"`
	Status       string     `gorm:"column:status" json:"status"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	TotalCount   int64      `gorm:"column:total_count" json:"total_count"`
}

// MerchantTrashedRow - paginated trashed merchant listing with deleted_at
type MerchantTrashedRow struct {
	MerchantID   int32      `gorm:"column:merchant_id" json:"merchant_id"`
	UserID       int32      `gorm:"column:user_id" json:"user_id"`
	Name         string     `gorm:"column:name" json:"name"`
	Description  *string    `gorm:"column:description" json:"description"`
	Address      *string    `gorm:"column:address" json:"address"`
	ContactEmail *string    `gorm:"column:contact_email" json:"contact_email"`
	ContactPhone *string    `gorm:"column:contact_phone" json:"contact_phone"`
	Status       string     `gorm:"column:status" json:"status"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	TotalCount   int64      `gorm:"column:total_count" json:"total_count"`
}

// MerchantDocumentRow - paginated merchant document listing without deleted_at
type MerchantDocumentRow struct {
	DocumentID   int32     `gorm:"column:document_id" json:"document_id"`
	MerchantID   int32     `gorm:"column:merchant_id" json:"merchant_id"`
	DocumentType string    `gorm:"column:document_type" json:"document_type"`
	DocumentUrl  string    `gorm:"column:document_url" json:"document_url"`
	Status       string    `gorm:"column:status" json:"status"`
	Note         *string   `gorm:"column:note" json:"note"`
	UploadedAt   time.Time `gorm:"column:uploaded_at" json:"uploaded_at"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
	TotalCount   int64     `gorm:"column:total_count" json:"total_count"`
}

// MerchantDocumentActiveRow - paginated active merchant document listing without deleted_at
type MerchantDocumentActiveRow struct {
	DocumentID   int32     `gorm:"column:document_id" json:"document_id"`
	MerchantID   int32     `gorm:"column:merchant_id" json:"merchant_id"`
	DocumentType string    `gorm:"column:document_type" json:"document_type"`
	DocumentUrl  string    `gorm:"column:document_url" json:"document_url"`
	Status       string    `gorm:"column:status" json:"status"`
	Note         *string   `gorm:"column:note" json:"note"`
	UploadedAt   time.Time `gorm:"column:uploaded_at" json:"uploaded_at"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
	TotalCount   int64     `gorm:"column:total_count" json:"total_count"`
}

// MerchantDocumentTrashedRow - paginated trashed merchant document listing with deleted_at
type MerchantDocumentTrashedRow struct {
	DocumentID   int32      `gorm:"column:document_id" json:"document_id"`
	MerchantID   int32      `gorm:"column:merchant_id" json:"merchant_id"`
	DocumentType string     `gorm:"column:document_type" json:"document_type"`
	DocumentUrl  string     `gorm:"column:document_url" json:"document_url"`
	Status       string     `gorm:"column:status" json:"status"`
	Note         *string    `gorm:"column:note" json:"note"`
	UploadedAt   time.Time  `gorm:"column:uploaded_at" json:"uploaded_at"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	TotalCount   int64      `gorm:"column:total_count" json:"total_count"`
}
