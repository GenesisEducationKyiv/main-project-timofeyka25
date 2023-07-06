package http

import (
	"genesis-test/src/app/domain"
	"genesis-test/src/app/presentation/http/responses"

	"github.com/gofiber/fiber/v2"
)

type NewsletterHandler struct {
	service domain.NewsletterService
}

func NewNewsletterHandler(s domain.NewsletterService) *NewsletterHandler {
	return &NewsletterHandler{
		service: s,
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
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Message: err.Error()})
	}

	if len(unsent) > 0 {
		return c.JSON(responses.SendRateResponse{UnsentEmails: unsent})
	}

	return c.SendStatus(fiber.StatusOK)
}
