package http

import (
	"genesis-test/src/app/domain/model"
	"genesis-test/src/app/presentation/http/response"
	"genesis-test/src/config"

	"github.com/gofiber/fiber/v2"
)

type ExchangeService interface {
	GetCurrencyRate(pair *model.CurrencyPair) (float64, error)
}

type ExchangeHandler struct {
	service   ExchangeService
	presenter ResponsePresenter
	logger    Logger
}

func NewExchangeHandler(
	s ExchangeService,
	p ResponsePresenter,
	l Logger,
) *ExchangeHandler {
	return &ExchangeHandler{
		service:   s,
		presenter: p,
		logger:    l,
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
	pair := &model.CurrencyPair{
		BaseCurrency:  model.CurrencyFactory{}.CreateCurrency(config.Get().BaseCurrency),
		QuoteCurrency: model.CurrencyFactory{}.CreateCurrency(config.Get().QuoteCurrency),
	}
	rate, err := h.service.GetCurrencyRate(pair)
	if err != nil {
		h.logger.Error(err.Error())
		return h.presenter.PresentError(c.Status(fiber.StatusBadRequest), &response.ErrorResponse{
			Message: err.Error(),
		})
	}
	return h.presenter.PresentExchangeRate(c, &response.RateResponse{Rate: rate})
}
