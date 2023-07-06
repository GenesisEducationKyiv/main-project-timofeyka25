package http

import (
	"genesis-test/src/app/customerror"
	"genesis-test/src/app/domain"
	"genesis-test/src/app/domain/model"
	"genesis-test/src/app/presentation/http/responses"
	"net/mail"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type SubscriptionHandler struct {
	service domain.SubscriptionService
}

func NewSubscriptionHandler(s domain.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		service: s,
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
func (h *SubscriptionHandler) Subscribe(c *fiber.Ctx) error {
	subscriber := new(model.Subscriber)

	if err := c.BodyParser(subscriber); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Message: err.Error()})
	}

	if _, err := mail.ParseAddress(subscriber.Email); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Message: err.Error()})
	}

	err := h.service.Subscribe(subscriber)
	if err != nil {
		if errors.Is(err, customerror.ErrAlreadyExists) {
			return c.SendStatus(fiber.StatusConflict)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Message: err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}
