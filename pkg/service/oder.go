package service

import (
	"goPromotion/dto"
	"goPromotion/pkg/repository"
)

type OrderItfService interface {
	GetServiceOrder(uint) (*dto.OrderResponse, error)
}

type orderImpService struct {
	orderRepo repository.OrderIrfRepository
}

func NewOrderImpService(orderRepo repository.OrderIrfRepository) OrderItfService {
	return &orderImpService{orderRepo: orderRepo}
}

func (o *orderImpService) GetServiceOrder(id uint) (*dto.OrderResponse, error) {
	order, err := o.orderRepo.GetPepoOrderByID(id)

	if err != nil {
		return nil, err
	}

	oderResponse := dto.OrderResponse{
		OrderID:  order.ID,
		Price:    order.TotalPrice,
		Discount: order.TotalDiscount,
		NetPrice: order.NetPrice,
	}

	return &oderResponse, nil
}
