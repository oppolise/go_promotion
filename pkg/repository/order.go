package repository

import (
	"goPromotion/pkg/model"

	"gorm.io/gorm"
)

type OrderIrfRepository interface {
	GetPepoOrderByID(uint) (*model.Order, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderIrfRepository {
	return &orderRepository{db: db}
}

func (o *orderRepository) GetPepoOrderByID(id uint) (*model.Order, error) {

	order := model.Order{}
	err := o.db.
		Preload("OrderDetails").
		Preload("OrderDetails.Product"). // ถ้าต้องใช้ชื่อสินค้า
		Preload("Promotions", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "condition_id") // เก็บเฉพาะ field ที่ต้องการ
		}).
		Preload("Promotions.Condition", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "payment_id", "ScopeDiscount", "min_price") // แล้วแต่ logic
		}).
		First(&order, id).Error

	if err != nil {
		return nil, err
	}

	return &order, nil
}
