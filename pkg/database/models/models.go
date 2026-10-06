package models

import (
	"time"

	"gorm.io/gorm"
)

// User maps to the "users" table.
type User struct {
	UserID           int32          `gorm:"column:user_id;primaryKey;autoIncrement" json:"user_id"`
	Firstname        string         `gorm:"column:firstname;size:100;not null" json:"firstname"`
	Lastname         string         `gorm:"column:lastname;size:100;not null" json:"lastname"`
	Email            string         `gorm:"column:email;size:100;uniqueIndex;not null" json:"email"`
	Password         string         `gorm:"column:password;size:100;not null" json:"password"`
	VerificationCode string         `gorm:"column:verification_code;size:100;not null" json:"verification_code"`
	IsVerified       *bool          `gorm:"column:is_verified;default:false" json:"is_verified"`
	CreatedAt        time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (User) TableName() string { return "users" }

// Role maps to the "roles" table.
type Role struct {
	RoleID    int32          `gorm:"column:role_id;primaryKey;autoIncrement" json:"role_id"`
	RoleName  string         `gorm:"column:role_name;size:50;uniqueIndex;not null" json:"role_name"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (Role) TableName() string { return "roles" }

// UserRole maps to the "user_roles" table.
type UserRole struct {
	UserRoleID int32          `gorm:"column:user_role_id;primaryKey;autoIncrement" json:"user_role_id"`
	UserID     int32          `gorm:"column:user_id;not null" json:"user_id"`
	RoleID     int32          `gorm:"column:role_id;not null" json:"role_id"`
	CreatedAt  time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (UserRole) TableName() string { return "user_roles" }

// RefreshToken maps to the "refresh_tokens" table.
type RefreshToken struct {
	RefreshTokenID int32          `gorm:"column:refresh_token_id;primaryKey;autoIncrement" json:"refresh_token_id"`
	UserID         int32          `gorm:"column:user_id;not null" json:"user_id"`
	Token          string         `gorm:"column:token;size:255;uniqueIndex;not null" json:"token"`
	Expiration     time.Time      `gorm:"column:expiration;not null" json:"expiration"`
	CreatedAt      time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (RefreshToken) TableName() string { return "refresh_tokens" }

// ResetToken maps to the "reset_tokens" table.
type ResetToken struct {
	ID         int32     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID     int64     `gorm:"column:user_id;not null" json:"user_id"`
	Token      string    `gorm:"column:token;type:text;not null" json:"token"`
	ExpiryDate time.Time `gorm:"column:expiry_date;not null" json:"expiry_date"`
}

func (ResetToken) TableName() string { return "reset_tokens" }

// Merchant maps to the "merchants" table.
type Merchant struct {
	MerchantID   int32          `gorm:"column:merchant_id;primaryKey;autoIncrement" json:"merchant_id"`
	UserID       int32          `gorm:"column:user_id;not null" json:"user_id"`
	Name         string         `gorm:"column:name;size:255;not null" json:"name"`
	Description  *string        `gorm:"column:description;type:text" json:"description"`
	Address      *string        `gorm:"column:address;type:text" json:"address"`
	ContactEmail *string        `gorm:"column:contact_email;size:100" json:"contact_email"`
	ContactPhone *string        `gorm:"column:contact_phone;size:20" json:"contact_phone"`
	Status       string         `gorm:"column:status;size:20;not null;default:active" json:"status"`
	CreatedAt    time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (Merchant) TableName() string { return "merchants" }

// MerchantDocument maps to the "merchant_documents" table.
type MerchantDocument struct {
	DocumentID   int32          `gorm:"column:document_id;primaryKey;autoIncrement" json:"document_id"`
	MerchantID   int32          `gorm:"column:merchant_id;not null" json:"merchant_id"`
	DocumentType string         `gorm:"column:document_type;size:50;not null" json:"document_type"`
	DocumentUrl  string         `gorm:"column:document_url;type:text;not null" json:"document_url"`
	Status       string         `gorm:"column:status;size:20;not null;default:pending" json:"status"`
	Note         *string        `gorm:"column:note;type:text" json:"note"`
	UploadedAt   time.Time      `gorm:"column:uploaded_at" json:"uploaded_at"`
	CreatedAt    time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (MerchantDocument) TableName() string { return "merchant_documents" }

// Cashier maps to the "cashiers" table.
type Cashier struct {
	CashierID int32          `gorm:"column:cashier_id;primaryKey;autoIncrement" json:"cashier_id"`
	MerchantID int32         `gorm:"column:merchant_id;not null" json:"merchant_id"`
	UserID     int32         `gorm:"column:user_id;not null" json:"user_id"`
	Name       string        `gorm:"column:name;size:100;not null" json:"name"`
	CreatedAt  time.Time     `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time     `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (Cashier) TableName() string { return "cashiers" }

// Category maps to the "categories" table.
type Category struct {
	CategoryID   int32          `gorm:"column:category_id;primaryKey;autoIncrement" json:"category_id"`
	Name         string         `gorm:"column:name;size:100;not null" json:"name"`
	Description  *string        `gorm:"column:description;type:text" json:"description"`
	SlugCategory *string        `gorm:"column:slug_category;size:100;uniqueIndex" json:"slug_category"`
	CreatedAt    time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (Category) TableName() string { return "categories" }

// Product maps to the "products" table.
type Product struct {
	ProductID    int32          `gorm:"column:product_id;primaryKey;autoIncrement" json:"product_id"`
	MerchantID   int32          `gorm:"column:merchant_id;not null" json:"merchant_id"`
	CategoryID   int32          `gorm:"column:category_id;not null" json:"category_id"`
	Name         string         `gorm:"column:name;size:255;not null" json:"name"`
	Description  *string        `gorm:"column:description;type:text" json:"description"`
	Price        int32          `gorm:"column:price;not null" json:"price"`
	CountInStock int32          `gorm:"column:count_in_stock;not null;default:0" json:"count_in_stock"`
	Brand        *string        `gorm:"column:brand;size:100" json:"brand"`
	Weight       *int32         `gorm:"column:weight" json:"weight"`
	SlugProduct  *string        `gorm:"column:slug_product;size:100;uniqueIndex" json:"slug_product"`
	ImageProduct *string        `gorm:"column:image_product;type:text" json:"image_product"`
	Barcode      *string        `gorm:"column:barcode;size:50;uniqueIndex" json:"barcode"`
	CreatedAt    time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (Product) TableName() string { return "products" }

// Order maps to the "orders" table.
type Order struct {
	OrderID    int32          `gorm:"column:order_id;primaryKey;autoIncrement" json:"order_id"`
	MerchantID int32          `gorm:"column:merchant_id;not null" json:"merchant_id"`
	CashierID  int32          `gorm:"column:cashier_id;not null" json:"cashier_id"`
	TotalPrice int64          `gorm:"column:total_price;not null" json:"total_price"`
	CreatedAt  time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (Order) TableName() string { return "orders" }

// OrderItem maps to the "order_items" table.
type OrderItem struct {
	OrderItemID int32          `gorm:"column:order_item_id;primaryKey;autoIncrement" json:"order_item_id"`
	OrderID     int32          `gorm:"column:order_id;not null" json:"order_id"`
	ProductID   int32          `gorm:"column:product_id;not null" json:"product_id"`
	Quantity    int32          `gorm:"column:quantity;not null" json:"quantity"`
	Price       int32          `gorm:"column:price;not null" json:"price"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (OrderItem) TableName() string { return "order_items" }

// Transaction maps to the "transactions" table.
type Transaction struct {
	TransactionID int32          `gorm:"column:transaction_id;primaryKey;autoIncrement" json:"transaction_id"`
	OrderID       int32          `gorm:"column:order_id;not null" json:"order_id"`
	MerchantID    int32          `gorm:"column:merchant_id;not null" json:"merchant_id"`
	PaymentMethod string         `gorm:"column:payment_method;size:50;not null" json:"payment_method"`
	Amount        int32          `gorm:"column:amount;not null" json:"amount"`
	ChangeAmount  *int32         `gorm:"column:change_amount;default:0" json:"change_amount"`
	PaymentStatus string         `gorm:"column:payment_status;size:20;not null;default:completed" json:"payment_status"`
	CreatedAt     time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}

func (Transaction) TableName() string { return "transactions" }

// OutboxEvent maps to the "outbox_events" table.
type OutboxEvent struct {
	OutboxID      int64     `gorm:"column:outbox_id;primaryKey;autoIncrement" json:"outbox_id"`
	Topic         string    `gorm:"column:topic;size:255;not null" json:"topic"`
	EventKey      string    `gorm:"column:event_key;size:255;not null" json:"event_key"`
	Payload       []byte    `gorm:"column:payload;type:jsonb;not null" json:"payload"`
	Status        string    `gorm:"column:status;size:20;not null;default:pending" json:"status"`
	Attempts      int32     `gorm:"column:attempts;not null;default:0" json:"attempts"`
	NextAttemptAt time.Time `gorm:"column:next_attempt_at;not null" json:"next_attempt_at"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (OutboxEvent) TableName() string { return "outbox_events" }

// ConsumerInbox maps to the "consumer_inbox" table with a composite primary
// key of (consumer_name, event_key).
type ConsumerInbox struct {
	ConsumerName       string     `gorm:"column:consumer_name;primaryKey;size:128;not null" json:"consumer_name"`
	EventKey           string     `gorm:"column:event_key;primaryKey;size:255;not null" json:"event_key"`
	Topic              string     `gorm:"column:topic;size:255;not null;default:''" json:"topic"`
	PartitionID        int32      `gorm:"column:partition_id;not null;default:-1" json:"partition_id"`
	MessageOffset      int64      `gorm:"column:message_offset;not null;default:-1" json:"message_offset"`
	Status             string     `gorm:"column:status;size:16;not null;default:processing" json:"status"`
	Attempts           int32      `gorm:"column:attempts;not null;default:1" json:"attempts"`
	LeaseUntil         time.Time  `gorm:"column:lease_until;not null" json:"lease_until"`
	LastError          string     `gorm:"column:last_error;type:text;not null;default:''" json:"last_error"`
	ProcessedAt        *time.Time `gorm:"column:processed_at" json:"processed_at"`
	ReservationVersion int64      `gorm:"column:reservation_version;not null;default:0" json:"reservation_version"`
}

func (ConsumerInbox) TableName() string { return "consumer_inbox" }
