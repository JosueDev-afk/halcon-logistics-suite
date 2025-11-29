package models

import (
	"time"

	"gorm.io/gorm"
)

// UserRole represents the role of a user in the system
type UserRole string

const (
	RoleAdmin      UserRole = "Admin"
	RoleSales      UserRole = "Sales"
	RolePurchasing UserRole = "Purchasing"
	RoleWarehouse  UserRole = "Warehouse"
	RoleRoute      UserRole = "Route"
)

// User represents a system user with role-based access
type User struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	Username     string         `gorm:"uniqueIndex;not null" json:"username"`
	PasswordHash string         `gorm:"not null" json:"-"`
	Role         UserRole       `gorm:"type:varchar(20);not null" json:"role"`
	Department   string         `gorm:"type:varchar(100)" json:"department"`
	FullName     string         `gorm:"type:varchar(200)" json:"full_name"`
	Email        string         `gorm:"type:varchar(200)" json:"email"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// OrderStatus represents the current state of an order
type OrderStatus string

const (
	StatusOrdered   OrderStatus = "Ordered"
	StatusInProcess OrderStatus = "In Process"
	StatusInRoute   OrderStatus = "In Route"
	StatusDelivered OrderStatus = "Delivered"
)

// Order represents a customer order with tracking and evidence
type Order struct {
	ID                uint           `gorm:"primarykey" json:"id"`
	InvoiceNumber     string         `gorm:"uniqueIndex;not null" json:"invoice_number"`
	CustomerName      string         `gorm:"type:varchar(200);not null" json:"customer_name"`
	CustomerNumber    string         `gorm:"type:varchar(100);not null;index" json:"customer_number"`
	Status            OrderStatus    `gorm:"type:varchar(20);not null;default:'Ordered'" json:"status"`
	DeliveryAddress   string         `gorm:"type:text" json:"delivery_address"`
	Notes             string         `gorm:"type:text" json:"notes"`
	EvidencePhotoURL  string         `gorm:"type:varchar(500)" json:"evidence_photo_url"`
	IsDeleted         bool           `gorm:"default:false;index" json:"is_deleted"`
	CreatedBy         uint           `gorm:"not null" json:"created_by"`
	CreatedByUser     User           `gorm:"foreignKey:CreatedBy" json:"created_by_user,omitempty"`
	LastModifiedBy    uint           `json:"last_modified_by"`
	LastModifiedUser  User           `gorm:"foreignKey:LastModifiedBy" json:"last_modified_by_user,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}

// TableName specifies the table name for Order model
func (Order) TableName() string {
	return "orders"
}
