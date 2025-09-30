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
	err := o.db.First(&order, id)

	if err.Error != nil {
		return nil, err.Error
	}

	return &order, nil
}
