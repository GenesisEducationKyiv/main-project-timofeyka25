package http

import (
	"genesis-test/src/app/domain"
	"genesis-test/src/app/presentation/http/response"

	"github.com/gofiber/fiber/v2"
)

type NewsletterHandler struct {
	service   domain.NewsletterService
	presenter ResponsePresenter
}

func NewNewsletterHandler(s domain.NewsletterService, p ResponsePresenter) *NewsletterHandler {
	return &NewsletterHandler{
		service:   s,
		presenter: p,
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
		return h.presenter.PresentError(c.Status(fiber.StatusBadRequest), &response.ErrorResponse{
			Message: err.Error(),
		})
	}

	if len(unsent) > 0 {
		return h.presenter.PresentSendRate(c, &response.SendRateResponse{UnsentEmails: unsent})
	}

	return c.SendStatus(fiber.StatusOK)
}
