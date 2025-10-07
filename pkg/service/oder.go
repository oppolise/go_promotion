package service

import (
	"errors"
	"goPromotion/dto"
	"goPromotion/pkg/model"
	"goPromotion/pkg/repository"
	"goPromotion/utils"
	"time"
)

var (
	ErrNoPromotions       = errors.New("ออเดอร์นี้ยังไม่ได้ใช้โปร")
	ErrNoActivePromotions = errors.New("ไม่มีโปรที่สามารถใช้ได้")
	ErrNoneWithOthers     = errors.New("พบโปร NONE ใช้ร่วมกับโปรอื่นไม่ได้")
)

type OrderItfService interface {
	GetServiceOrder(uint) (*dto.PromotionResult, error)
}

type orderImpService struct {
	orderRepo repository.OrderIrfRepository
}

func NewOrderImpService(orderRepo repository.OrderIrfRepository) OrderItfService {
	return &orderImpService{orderRepo: orderRepo}
}

func (o *orderImpService) GetServiceOrder(id uint) (*dto.PromotionResult, error) {

	now := time.Now()

	order, err := o.orderRepo.GetPepoOrderByID(id)
	if err != nil {
		return nil, err
	}

	if len(order.Promotions) == 0 {
		return nil, ErrNoPromotions
	}

	result := &dto.PromotionResult{
		OrderID:         order.ID,
		RejectedReasons: make(map[uint]string),
		AppliedPromos:   []uint{},
		LineDiscounts:   make(map[uint]int),
	}

	activePromos := []model.Promotion{}

	for _, promo := range order.Promotions {
		if !utils.IsPromotionActive(promo, now) {
			result.RejectedReasons[promo.ID] = "ไม่อยู่ในช่วงเวลา"
			continue
		}
		activePromos = append(activePromos, promo)
	}

	if len(activePromos) == 0 {
		return nil, ErrNoActivePromotions
	}

	hasNone := utils.HasNonePromotion(activePromos)

	if hasNone && len(activePromos) > 1 {
		return nil, ErrNoneWithOthers
	}

	itemPromos := []model.Promotion{}
	orderPromos := []model.Promotion{}

	for _, promo := range activePromos {
		if promo.Condition.ScopeDiscount == "ITEM" {
			itemPromos = append(itemPromos, promo)
		} else if promo.Condition.ScopeDiscount == "ORDER" {
			orderPromos = append(orderPromos, promo)
		}
	}

	eligibleItemPromos := []model.Promotion{}

	for _, promo := range itemPromos {
		ok, reason := utils.CheckConditionProductsITEM(*promo.Condition, order.OrderDetails)
		if !ok {
			result.RejectedReasons[promo.ID] = reason
			continue
		}

		totalQty := utils.SumQuantity(order.OrderDetails)
		if totalQty < int(promo.Condition.MinQuantityItem) {
			result.RejectedReasons[promo.ID] = "จำนวนสินค้าไม่ถึงขั้นต่ำ"
			continue
		}

		totalPrice := utils.SumLineSubtotal(order.OrderDetails)
		if totalPrice < int(promo.Condition.MinPrice) {
			result.RejectedReasons[promo.ID] = "ยอดรวมไม่ถึงขั้นต่ำ"
			continue
		}

		eligibleItemPromos = append(eligibleItemPromos, promo)
	}

	utils.SortPromotionsByPriority(eligibleItemPromos)

	totalItemDiscount := 0
	appliedItemPromos := []uint{}

	for _, promo := range eligibleItemPromos {
		discount := utils.CalculateItemDiscount(promo, order.OrderDetails)
		if discount > 0 {
			totalItemDiscount += discount
			appliedItemPromos = append(appliedItemPromos, promo.ID)
		}
	}

	itemLineDiscounts := utils.DistributeDiscount(order.OrderDetails, totalItemDiscount)
	for detailID, discount := range itemLineDiscounts {
		result.LineDiscounts[detailID] = discount
	}

	originalSubtotal := utils.SumLineSubtotal(order.OrderDetails)
	subtotalAfterItem := originalSubtotal - totalItemDiscount

	eligibleOrderPromos := []model.Promotion{}

	for _, promo := range orderPromos {

		ok, reason := utils.CheckConditionProductsORDER(*promo.Condition, order.OrderDetails)
		if !ok {
			result.RejectedReasons[promo.ID] = reason
			continue
		}

		totalQty := utils.SumQuantity(order.OrderDetails)
		if totalQty < int(promo.Condition.MinQuantityItem) {
			result.RejectedReasons[promo.ID] = "จำนวนสินค้าไม่ถึงขั้นต่ำ"
			continue
		}

		if subtotalAfterItem < int(promo.Condition.MinPrice) {
			result.RejectedReasons[promo.ID] = "ยอดรวมไม่ถึงขั้นต่ำ"
			continue
		}

		eligibleOrderPromos = append(eligibleOrderPromos, promo)
	}
	utils.SortPromotionsByPriority(eligibleOrderPromos)

	totalCartDiscount := 0
	appliedCartPromos := []uint{}
	remainingSubtotal := subtotalAfterItem

	for _, promo := range eligibleOrderPromos {
		discount := utils.CalculateOrderDiscount(promo, remainingSubtotal)
		if discount > 0 {
			totalCartDiscount += discount
			appliedCartPromos = append(appliedCartPromos, promo.ID)
			remainingSubtotal -= discount
		}
	}

	if totalCartDiscount > 0 {
		cartLineDiscounts := utils.DistributeDiscount(order.OrderDetails, totalCartDiscount)
		for detailID, discount := range cartLineDiscounts {
			result.LineDiscounts[detailID] += discount
		}
	}

	result.AppliedPromos = append(appliedItemPromos, appliedCartPromos...)
	result.OriginalSubtotal = originalSubtotal
	result.TotalItemDiscount = totalItemDiscount
	result.SubtotalAfterItem = subtotalAfterItem
	result.TotalCartDiscount = totalCartDiscount
	result.TotalDiscount = totalItemDiscount + totalCartDiscount
	result.FinalNetPrice = originalSubtotal - result.TotalDiscount

	return result, nil

}
