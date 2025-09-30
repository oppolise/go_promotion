package dto

type OrderResponse struct {
	OrderID  uint `json:"orderId"`
	Price    uint `json:"price"`
	Discount uint `json:"discount"`
	NetPrice uint `json:"netPrice"`
}
