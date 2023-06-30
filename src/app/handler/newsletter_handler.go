package handler

import (
	"genesis-test/src/app/customerror"
	"genesis-test/src/app/domain"
	"genesis-test/src/app/service"
	"net/mail"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type NewsletterHandler struct {
	services *service.Services
}

func NewNewsletterHandler(s *service.Services) *NewsletterHandler {
	return &NewsletterHandler{
		services: s,
	}
}

// Subscribe
//
//	@Summary	Subscribe to newsletter
//	@Tags		Newsletter
//	@Param		input	body 	domain.Subscriber true "Email to subscribe"
//	@Accept		json
//	@Success	200
//	@Failure	400	{object}	ErrorResponse
//	@Failure	409
//	@Failure	500 {object}	ErrorResponse
//	@Router		/subscribe [post]
func (h NewsletterHandler) Subscribe(c *fiber.Ctx) error {
	subscriber := new(domain.Subscriber)

	if err := c.BodyParser(subscriber); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{err.Error()})
	}

	if _, err := mail.ParseAddress(subscriber.Email); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{err.Error()})
	}

	err := h.services.Newsletter.Subscribe(subscriber)
	if err != nil {
		if errors.Is(err, customerror.ErrAlreadyExists) {
			return c.SendStatus(fiber.StatusConflict)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
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
	unsent, err := h.services.Newsletter.SendEmails()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{err.Error()})
	}

	if len(unsent) > 0 {
		return c.JSON(fiber.Map{
			"unsent": unsent,
		})
	}

	return c.SendStatus(fiber.StatusOK)
}
