package http

import (
	"genesis-test/src/app/presentation/http/response"

	"github.com/gofiber/fiber/v2"
)

type NewsletterService interface {
	SendCurrencyRate() ([]string, error)
}

type NewsletterHandler struct {
	service   NewsletterService
	presenter ResponsePresenter
	logger    Logger
}

func NewNewsletterHandler(
	s NewsletterService,
	p ResponsePresenter,
	l Logger,
) *NewsletterHandler {
	return &NewsletterHandler{
		service:   s,
		presenter: p,
		logger:    l,
	}
}

// SendEmails
//
//	@Summary	Send currency rate to subscribed emails
//	@Tags		Newsletter
//	@Accept		json
//	@Success	200		{object}	fiber.Map
//	@Failure	400		{object}	ErrorResponse
//	@Router		/sendEmails [post]
func (h NewsletterHandler) SendEmails(c *fiber.Ctx) error {
	unsent, err := h.service.SendCurrencyRate()
	if err != nil {
		h.logger.Error(err.Error())
		return h.presenter.PresentError(c.Status(fiber.StatusBadRequest), &response.ErrorResponse{
			Message: err.Error(),
		})
	}

	if len(unsent) > 0 {
		return h.presenter.PresentSendRate(c, &response.SendRateResponse{UnsentEmails: unsent})
	}

	return c.SendStatus(fiber.StatusOK)
}
