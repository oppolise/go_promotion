package model

import (
	"time"
)

// Category model
type Category struct {
	ID   uint   `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"type:varchar(255);not null"`

	// Relationships
	Products []Product `gorm:"foreignKey:CategoryID"`
}

// Product model
type Product struct {
	ID         uint    `gorm:"primaryKey;autoIncrement"`
	CategoryID uint    `gorm:"not null;index"`
	Quantity   int     `gorm:"not null;default:0"`
	Name       string  `gorm:"type:varchar(255);not null"`
	Price      float64 `gorm:"type:decimal(10,2);not null"`

	// Relationships
	Category     Category      `gorm:"foreignKey:CategoryID"`
	OrderDetails []OrderDetail `gorm:"foreignKey:ProductID"`
}

// Payment model
type Payment struct {
	ID   uint   `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"type:varchar(255);not null"`

	// Relationships
	Conditions []Condition `gorm:"foreignKey:PaymentID"`
}

// Condition model
type Condition struct {
	ID                 uint   `gorm:"primaryKey;autoIncrement"`
	PaymentID          uint   `gorm:"index"`
	ScopeDiscount      string `gorm:"type:varchar(100)"`
	MinQuantityItem    uint
	MaxQuantityItem    uint
	MinPrice           uint
	MaxDiscount        uint
	DiscountUnit       string `gorm:"type:varchar(50)"`
	DiscountValue      int
	RepeatableDiscount bool `gorm:"default:false"`
	DateStart          time.Time
	DateEnd            time.Time

	// Relationships
	Payment    Payment    `gorm:"foreignKey:PaymentID"`
	Promotion  *Promotion `gorm:"foreignKey:ConditionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Products   []Product  `gorm:"many2many:condition_product;"`
	Categories []Category `gorm:"many2many:condition_category;"`
}

// Promotion model
type Promotion struct {
	ID                  uint   `gorm:"primaryKey;autoIncrement"`
	ConditionID         uint   `gorm:"uniqueIndex;not null"`
	Code                string `gorm:"type:varchar(100);uniqueIndex;not null"`
	Description         string `gorm:"type:text"`
	Quantity            uint
	MaxQuantityPerDay   uint
	MaxQuantityPerMonth uint
	TypeStack           string    `gorm:"type:varchar(50)"`
	DateStart           time.Time `gorm:"not null"`
	DateEnd             time.Time `gorm:"not null"`

	// Relationships
	Condition *Condition `gorm:"foreignKey:ConditionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Orders    []Order    `gorm:"many2many:order_promotion;"`
}

// Order model
type Order struct {
	ID            uint `gorm:"primaryKey;autoIncrement"`
	TotalPrice    uint `gorm:"not null;default:0"`
	TotalDiscount uint `gorm:"default:0"`
	NetPrice      uint `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time

	// Relationships
	OrderDetails []OrderDetail `gorm:"foreignKey:OrderID"`
	Promotions   []Promotion   `gorm:"many2many:order_promotion;"`
}

// OrderDetail model
type OrderDetail struct {
	OrderID    uint `gorm:"primaryKey;not null"`
	ProductID  uint `gorm:"primaryKey;not null"`
	Quantity   uint `gorm:"not null"`
	TotalPrice uint `gorm:"not null"`

	// Relationships
	Order   Order   `gorm:"foreignKey:OrderID"`
	Product Product `gorm:"foreignKey:ProductID"`
}
