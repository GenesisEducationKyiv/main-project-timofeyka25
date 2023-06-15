package handler

import (
	"genesis-test/src/app/service"
	"github.com/gofiber/fiber/v2"
)

type ExchangeHandler struct {
	services *service.Services
}

func NewExchangeHandler(s *service.Services) *ExchangeHandler {
	return &ExchangeHandler{
		services: s,
	}
}

// GetCurrencyRate
//
//	@Summary	Get currency rate
//	@Tags		Exchange
//	@Accept		json
//	@Produce	json
//	@Success	200		{integer}   integer
//	@Failure	400		{object}	ErrorResponse
//	@Router		/rate [get]
func (h ExchangeHandler) GetCurrencyRate(c *fiber.Ctx) error {
	rate, err := h.services.Exchange.GetCurrencyRate()
	if err != nil {
		// Return status 400 and error message
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{err.Error()})
	}

	return c.JSON(rate)
}
