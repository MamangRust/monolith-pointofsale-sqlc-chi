package models

import "time"

// UserRow - paginated active user listing without deleted_at
type UserRow struct {
	UserID           int32     `gorm:"column:user_id" json:"user_id"`
	Firstname        string    `gorm:"column:firstname" json:"firstname"`
	Lastname         string    `gorm:"column:lastname" json:"lastname"`
	Email            string    `gorm:"column:email" json:"email"`
	Password         string    `gorm:"column:password" json:"password"`
	VerificationCode string    `gorm:"column:verification_code" json:"verification_code"`
	IsVerified       *bool     `gorm:"column:is_verified" json:"is_verified"`
	CreatedAt        time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at" json:"updated_at"`
	TotalCount       int64     `gorm:"column:total_count" json:"total_count"`
}

// UserActiveRow - paginated active user listing with deleted_at
type UserActiveRow struct {
	UserID           int32      `gorm:"column:user_id" json:"user_id"`
	Firstname        string     `gorm:"column:firstname" json:"firstname"`
	Lastname         string     `gorm:"column:lastname" json:"lastname"`
	Email            string     `gorm:"column:email" json:"email"`
	Password         string     `gorm:"column:password" json:"password"`
	VerificationCode string     `gorm:"column:verification_code" json:"verification_code"`
	IsVerified       *bool      `gorm:"column:is_verified" json:"is_verified"`
	CreatedAt        time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt        *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	TotalCount       int64      `gorm:"column:total_count" json:"total_count"`
}

// UserTrashedRow - paginated trashed user listing with deleted_at
type UserTrashedRow struct {
	UserID           int32      `gorm:"column:user_id" json:"user_id"`
	Firstname        string     `gorm:"column:firstname" json:"firstname"`
	Lastname         string     `gorm:"column:lastname" json:"lastname"`
	Email            string     `gorm:"column:email" json:"email"`
	Password         string     `gorm:"column:password" json:"password"`
	VerificationCode string     `gorm:"column:verification_code" json:"verification_code"`
	IsVerified       *bool      `gorm:"column:is_verified" json:"is_verified"`
	CreatedAt        time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt        *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	TotalCount       int64      `gorm:"column:total_count" json:"total_count"`
}
