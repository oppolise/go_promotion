package repository

import (
	"goPromotion/pkg/model"

	"gorm.io/gorm"
)

type orderImpRepository interface {
	GetOrderByID(uint) (*model.Order, error)
	UpdeteOrder(uint, *model.Order) (*model.Order, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) orderImpRepository {
	return &orderRepository{db: db}
}

func (o *orderRepository) GetOrderByID(id uint) (*model.Order, error) {

	return nil, nil
}

func (o *orderRepository) UpdeteOrder(id uint, order *model.Order) (*model.Order, error) {

	return nil, nil
}
