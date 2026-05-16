package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Email        string         `gorm:"uniqueIndex;not null;size:255" json:"email"`
	Username     string         `gorm:"uniqueIndex;not null;size:64" json:"username"`
	PasswordHash string         `gorm:"not null;size:255" json:"-"`
	DisplayName  string         `gorm:"size:255" json:"display_name"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Roles        []Role         `gorm:"many2many:user_roles;" json:"roles,omitempty"`
}

type Role struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"uniqueIndex;not null;size:64" json:"name"`
	Label       string         `gorm:"size:255" json:"label"`
	Description string         `gorm:"size:512" json:"description,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Permissions []RolePermission `gorm:"foreignKey:RoleID" json:"permissions,omitempty"`
}

type RolePermission struct {
	RoleID     uint   `gorm:"primaryKey" json:"role_id"`
	Permission string `gorm:"primaryKey;size:128" json:"permission"`
}

type Product struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null;size:255" json:"name"`
	Slug        string         `gorm:"uniqueIndex;not null;size:255" json:"slug"`
	Description string         `gorm:"type:text" json:"description"`
	PriceCents  int64          `gorm:"not null" json:"price_cents"`
	Active      bool           `gorm:"not null;default:true;index" json:"active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Images      []ProductImage `gorm:"foreignKey:ProductID" json:"images,omitempty"`
}

type ProductImage struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ProductID   uint      `gorm:"not null;index" json:"product_id"`
	Path        string    `gorm:"not null;size:512" json:"path"`
	SortOrder   int       `gorm:"not null;default:0" json:"sort_order"`
	ExternalURL string    `gorm:"size:1024" json:"external_url,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// StoreBranch is a physical Tofex pickup location (not a bank branch).
type StoreBranch struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"not null;size:255" json:"name"`
	Address   string         `gorm:"type:text;not null" json:"address"`
	City      string         `gorm:"size:128" json:"city,omitempty"`
	Phone     string         `gorm:"size:64" json:"phone,omitempty"`
	Hours     string         `gorm:"size:255" json:"hours,omitempty"`
	MapsURL   string         `gorm:"size:1024" json:"maps_url,omitempty"`
	Latitude  *float64       `json:"latitude,omitempty"`
	Longitude *float64       `json:"longitude,omitempty"`
	Active    bool           `gorm:"not null;default:true;index" json:"active"`
	SortOrder int            `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type PaymentAccount struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	BankName       string         `gorm:"not null;size:255" json:"bank_name"`
	AccountName    string         `gorm:"not null;size:255" json:"account_name"`
	AccountNumber  string         `gorm:"not null;size:64" json:"account_number"`
	Branch         string         `gorm:"size:255" json:"branch,omitempty"`
	Currency       string         `gorm:"not null;size:8;default:NGN" json:"currency"`
	DisplayLabel   string         `gorm:"size:255" json:"display_label,omitempty"`
	Active         bool           `gorm:"not null;default:true;index" json:"active"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

type Order struct {
	ID                   uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	CustomerName         string         `gorm:"not null;size:255" json:"customer_name"`
	CustomerEmail        string         `gorm:"not null;size:255;index" json:"customer_email"`
	CustomerPhone        string         `gorm:"size:64" json:"customer_phone,omitempty"`
	PaymentAccountID     uint           `gorm:"not null;index" json:"payment_account_id"`
	PaymentSnapshotJSON  string         `gorm:"type:json;not null" json:"payment_snapshot"`
	TotalCents           int64          `gorm:"not null" json:"total_cents"`
	Status               string         `gorm:"not null;size:32;default:'pending';index" json:"status"`
	CustomerNote         string         `gorm:"type:text" json:"customer_note,omitempty"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`
	Items                []OrderItem    `gorm:"foreignKey:OrderID;references:ID" json:"items,omitempty"`
	Receipts             []OrderReceipt `gorm:"foreignKey:OrderID;references:ID" json:"receipts,omitempty"`
}

type OrderItem struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	OrderID          uuid.UUID `gorm:"type:char(36);not null;index" json:"order_id"`
	ProductID        uint      `gorm:"not null;index" json:"product_id"`
	ProductNameSnap  string    `gorm:"not null;size:255" json:"product_name"`
	UnitPriceCents   int64     `gorm:"not null" json:"unit_price_cents"`
	Quantity         int       `gorm:"not null" json:"quantity"`
}

type OrderReceipt struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	OrderID   uuid.UUID `gorm:"type:char(36);not null;index" json:"order_id"`
	FilePath  string    `gorm:"not null;size:512" json:"path"`
	MIME      string    `gorm:"size:128" json:"mime"`
	SizeBytes int64     `json:"size_bytes"`
	CreatedAt time.Time `json:"created_at"`
}

type AuditLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       *uint     `gorm:"index" json:"user_id,omitempty"`
	Action       string    `gorm:"not null;size:128;index" json:"action"`
	ResourceType string    `gorm:"size:64;index" json:"resource_type,omitempty"`
	ResourceID   string    `gorm:"size:64" json:"resource_id,omitempty"`
	MetadataJSON string    `gorm:"type:json" json:"metadata,omitempty"`
	IP           string    `gorm:"size:64" json:"ip,omitempty"`
	UserAgent    string    `gorm:"size:512" json:"user_agent,omitempty"`
	CreatedAt    time.Time `gorm:"index" json:"created_at"`
}

type PasswordReset struct {
	ID        uint       `gorm:"primaryKey" json:"-"`
	Email     string     `gorm:"not null;size:255;index" json:"-"`
	TokenHash string     `gorm:"not null;size:128;index" json:"-"`
	ExpiresAt time.Time  `gorm:"not null" json:"-"`
	UsedAt    *time.Time `json:"-"`
	CreatedAt time.Time  `json:"-"`
}
