package handler

import (
	"genesis-test/src/app/domain"
	"genesis-test/src/app/errors"
	"genesis-test/src/app/service"
	"github.com/gofiber/fiber/v2"
	"net/mail"
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
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{errors.ErrInvalidEmail.Error()})

	}

	err := h.services.Newsletter.Subscribe(subscriber)
	if err != nil {
		if err == errors.ErrAlreadyExists {
			// Return status 409 if subscriber already exists
			return c.SendStatus(fiber.StatusConflict)
		}
		// Return status 500 and error message for other errors
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
		// Return status 400 and error message
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{err.Error()})
	}

	if len(unsent) > 0 {
		// Return status 200 with the list of unsent recipients
		return c.JSON(fiber.Map{
			"unsent": unsent,
		})
	}

	// Return status 200 if all subscribers were successfully sent the email
	return c.SendStatus(fiber.StatusOK)
}
