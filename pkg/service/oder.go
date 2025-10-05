package service

import (
	// "goPromotion/dto"

	"goPromotion/pkg/model"
	"goPromotion/pkg/repository"
)

// type OrderItfService interface {
// 	GetServiceOrder(uint) (*dto.OrderResponse, error)
// }

type OrderItfService interface {
	GetServiceOrder(uint) (*model.Order, error)
}

type orderImpService struct {
	orderRepo repository.OrderIrfRepository
}

func NewOrderImpService(orderRepo repository.OrderIrfRepository) OrderItfService {
	return &orderImpService{orderRepo: orderRepo}
}

func (o *orderImpService) GetServiceOrder(id uint) (*model.Order, error) {

	order, err := o.orderRepo.GetPepoOrderByID(id)
	if err != nil {
		return nil, err
	}

	return order, nil

	// return &oderResponse, nil
}

// func (o *orderRepository) GetOrderForValidation(ctx context.Context, id uint) (*model.Order, error) {
//     var order model.Order
//     err := o.db.WithContext(ctx).
//         Select("id", "total_price", "total_discount", "net_price", "created_at"). // order columns
//         Preload("OrderDetails", func(db *gorm.DB) *gorm.DB {
//             return db.Select("order_id", "product_id", "quantity", "total_price")
//         }).
//         Preload("Promotions", func(db *gorm.DB) *gorm.DB {
//             return db.Select("id", "condition_id", "code", "type_stack")
//         }).
//         Preload("Promotions.Condition", func(db *gorm.DB) *gorm.DB {
//             return db.Select(
//                 "id", "scope_discount", "payment_id", "min_quantity_item",
//                 "max_quantity_item", "min_price", "max_discount",
//                 "discount_unit", "discount_value", "repeatable_discount",
//                 "date_start", "date_end",
//             )
//         }).
//         First(&order, id).Error
//     if err != nil {
//         return nil, err
//     }
//     return &order, nil
// }

// func (s *orderService) ValidateOrder(ctx context.Context, orderID uint) error {
//     order, err := s.repo.GetOrderForValidation(ctx, orderID)
//     if err != nil {
//         return err
//     }

//     if len(order.Promotions) > maxPromotionsAllowed {
//         return ErrTooManyPromotions
//     }

//     for _, promo := range order.Promotions {
//         cond := promo.Condition
//         if cond == nil {
//             return ErrConditionMissing
//         }
//         if !cond.DateStart.Before(time.Now()) || !cond.DateEnd.After(time.Now()) {
//             return ErrConditionExpired
//         }
//         // ตรวจ scope, product type, min price ฯลฯ ตาม logic คุณ
//     }

//     // ตรวจ order details / ราคาสินค้า ฯลฯ
//     return nil
