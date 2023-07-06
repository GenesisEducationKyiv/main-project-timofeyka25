package http

import (
	"genesis-test/src/app/domain"
	"genesis-test/src/app/presentation/http/responses"
	"github.com/gofiber/fiber/v2"
)

type ExchangeHandler struct {
	service domain.ExchangeService
}

func NewExchangeHandler(s domain.ExchangeService) *ExchangeHandler {
	return &ExchangeHandler{
		service: s,
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
func (h *ExchangeHandler) GetCurrencyRate(c *fiber.Ctx) error {
	rate, err := h.service.GetCurrencyRate()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Message: err.Error()})
	}

	return c.JSON(rate)
}
