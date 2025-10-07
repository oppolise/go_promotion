package dto

type OrderResponse struct {
	OrderID  uint `json:"orderId"`
	Price    uint `json:"price"`
	Discount uint `json:"discount"`
	NetPrice uint `json:"netPrice"`
}

type PromotionResult struct {
	OrderID           uint            `json:"order_id"`
	AppliedPromos     []uint          `json:"applied_promos"`
	RejectedReasons   map[uint]string `json:"rejected_reasons"`
	LineDiscounts     map[uint]int    `json:"line_discounts"`
	OriginalSubtotal  int             `json:"original_subtotal"`
	TotalItemDiscount int             `json:"total_item_discount"`
	SubtotalAfterItem int             `json:"subtotal_after_item"`
	TotalCartDiscount int             `json:"total_cart_discount"`
	TotalDiscount     int             `json:"total_discount"`
	FinalNetPrice     int             `json:"final_net_price"`
}
