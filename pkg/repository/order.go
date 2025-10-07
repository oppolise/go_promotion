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
		Preload("OrderDetails.Product").
		Preload("Promotions").
		Preload("Promotions.Condition").
		Preload("Promotions.Condition.Payment").
		Preload("Promotions.Condition.Products"). // ← สำหรับ CP
		Preload("Promotions.Condition.Categories").
		First(&order, id).Error

	if err != nil {
		return nil, err
	}

	return &order, nil
}
