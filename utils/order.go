package utils

import (
	"goPromotion/pkg/model"
	"sort"
	"time"
)

func IsPromotionActive(promo model.Promotion, now time.Time) bool {
	return !now.Before(promo.DateStart) && !now.After(promo.DateEnd)
}

func HasNonePromotion(promos []model.Promotion) bool {

	for _, p := range promos {
		if p.TypeStack == "NONE_STACKABLE" {
			return true
		}
	}
	return false
}

func CheckConditionProductsITEM(condition model.Condition, orderDetails []model.OrderDetail) (bool, string) {
	if len(condition.Products) == 0 {
		return false, "Promotion นี้ไม่ได้ระบุ product"
	}

	requiredProductIDs := make(map[uint]bool)
	for _, prod := range condition.Products {
		requiredProductIDs[prod.ID] = true
	}

	orderProductIDs := make(map[uint]bool)
	for _, detail := range orderDetails {
		orderProductIDs[detail.ProductID] = true
	}
	for reqID := range requiredProductIDs {
		if !orderProductIDs[reqID] {
			return false, "ไม่มีประเภทสินค้าครบตามเงื่อนไข promotion"
		}
	}

	hasExtraProducts := false
	for orderProdID := range orderProductIDs {
		if !requiredProductIDs[orderProdID] {
			hasExtraProducts = true
			break
		}
	}

	if !hasExtraProducts {
		return true, ""
	}

	if len(condition.Categories) == 0 {
		return false, "ไม่มีประเภทสินค้าตามเงื่อนไข promotion"
	}

	categortIDs := make(map[uint]bool)
	for _, cat := range condition.Categories {
		categortIDs[cat.ID] = true
	}

	for _, detail := range orderDetails {
		if requiredProductIDs[detail.ProductID] {
			continue
		}

		if !categortIDs[detail.Product.CategoryID] {
			return false, "ประเภทสินค้าไม่ตรงตามเงื่อนไข promotion"
		}
	}

	return true, ""
}

func CheckConditionProductsORDER(condition model.Condition, orderDetails []model.OrderDetail) (bool, string) {

	if len(condition.Products) == 0 {
		if len(condition.Categories) > 0 {
			categoryIDs := make(map[uint]bool)
			for _, cat := range condition.Categories {
				categoryIDs[cat.ID] = true
			}

			for _, detail := range orderDetails {
				if !categoryIDs[detail.Product.CategoryID] {
					return false, "ประเภทสินค้าไม่ตรงกับ promotion นี้"
				}
			}
		}
		return true, ""
	}

	requiredProductIDs := make(map[uint]bool)
	for _, prod := range condition.Products {
		requiredProductIDs[prod.ID] = true
	}

	orderProductIDs := make(map[uint]bool)
	for _, detail := range orderDetails {
		orderProductIDs[detail.ProductID] = true
	}

	for reqID := range requiredProductIDs {
		if !orderProductIDs[reqID] {
			return false, "ไม่มีสินค้าครบตามเงื่อนไข promotion"
		}
	}
	hasExtraProducts := false
	for orderProdID := range orderProductIDs {
		if !requiredProductIDs[orderProdID] {
			hasExtraProducts = true
			break
		}
	}

	if !hasExtraProducts {
		return true, ""
	}

	if len(condition.Categories) == 0 {
		return true, ""
	}

	categoryIDs := make(map[uint]bool)
	for _, cat := range condition.Categories {
		categoryIDs[cat.ID] = true
	}

	for _, detail := range orderDetails {
		if requiredProductIDs[detail.ProductID] {
			continue
		}

		if !categoryIDs[detail.Product.CategoryID] {
			return false, "ประเภทสินค้านี้ไม่ตรงตาม promotion"
		}
	}

	return true, ""
}

func SumQuantity(orderDetails []model.OrderDetail) int {
	total := 0
	for _, detail := range orderDetails {
		total += int(detail.Quantity)
	}
	return total
}

func SumLineSubtotal(orderDetails []model.OrderDetail) int {
	total := 0
	for _, detail := range orderDetails {
		total += int(detail.TotalPrice)
	}
	return total
}

func SortPromotionsByPriority(promos []model.Promotion) {
	sort.Slice(promos, func(i, j int) bool {
		condI := promos[i].Condition
		condJ := promos[j].Condition

		if condI.DiscountUnit != condJ.DiscountUnit {
			return condI.DiscountUnit == "PERCENT"
		}

		return condI.DiscountValue > condJ.DiscountValue
	})
}

func CalculateItemDiscount(promo model.Promotion, orderDetails []model.OrderDetail) int {

	condition := promo.Condition
	totalPrice := SumLineSubtotal(orderDetails)
	discount := 0

	if condition.DiscountUnit == "PERCENT" {
		discount = (totalPrice * condition.DiscountValue) / 100
	} else if condition.DiscountUnit == "CURRENCY" {
		discount = condition.DiscountValue
	}

	if condition.MaxDiscount > 0 && discount > int(condition.MaxDiscount) {
		discount = int(condition.MaxDiscount)
	}

	if discount > totalPrice {
		discount = totalPrice
	}

	return discount
}

func CalculateOrderDiscount(promo model.Promotion, subtotal int) int {
	condition := promo.Condition
	discount := 0

	if condition.DiscountUnit == "PERCENT" {
		discount = (subtotal * condition.DiscountValue) / 100
	} else if condition.DiscountUnit == "CURRENCY" {
		discount = condition.DiscountValue
	}

	if condition.MaxDiscount > 0 && discount > int(condition.MaxDiscount) {
		discount = int(condition.MaxDiscount)
	}

	if discount > subtotal {
		discount = subtotal
	}

	return discount
}

func DistributeDiscount(orderDetails []model.OrderDetail, totalDiscount int) map[uint]int {
	lineDiscounts := make(map[uint]int)
	totalPrice := SumLineSubtotal(orderDetails)

	if totalPrice == 0 {
		return lineDiscounts
	}

	remainingDiscount := totalDiscount

	for i, detail := range orderDetails {
		var lineDiscount int

		if i == len(orderDetails)-1 {
			lineDiscount = remainingDiscount
		} else {
			lineDiscount = int(detail.TotalPrice) * totalDiscount / totalPrice
			remainingDiscount -= lineDiscount
		}
		if lineDiscount > int(detail.TotalPrice) {
			lineDiscount = int(detail.TotalPrice)
		}

		lineDiscounts[detail.OrderID] = lineDiscount
	}

	return lineDiscounts
}
