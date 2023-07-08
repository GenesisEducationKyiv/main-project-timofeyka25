package http

import (
	"genesis-test/src/app/domain"
	"genesis-test/src/app/presentation/http/response"

	"github.com/gofiber/fiber/v2"
)

type ExchangeHandler struct {
	service   domain.ExchangeService
	presenter ResponsePresenter
}

func NewExchangeHandler(s domain.ExchangeService, p ResponsePresenter) *ExchangeHandler {
	return &ExchangeHandler{
		service:   s,
		presenter: p,
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
		return h.presenter.PresentError(c.Status(fiber.StatusBadRequest), &response.ErrorResponse{
			Message: err.Error(),
		})
	}

	return h.presenter.PresentExchangeRate(c, &response.RateResponse{Rate: rate})
}
