package handler

import (
	"goPromotion/pkg/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type HttpOrderHandler struct {
	service service.OrderItfService
}

func NewHttpOrderHandler(service service.OrderItfService) *HttpOrderHandler {
	return &HttpOrderHandler{service: service}
}

func (h *HttpOrderHandler) GetOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	orderID, err := strconv.Atoi(id)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	order, err := h.service.GetServiceOrder(uint(orderID))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err)
	}

	return c.Status(fiber.StatusOK).JSON(order)
}
