package http

import (
	"genesis-test/src/app/customerror"
	"genesis-test/src/app/domain/model"
	"genesis-test/src/app/presentation/http/response"
	"net/mail"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type SubscriptionService interface {
	Subscribe(subscriber *model.Subscriber) error
}

type SubscriptionHandler struct {
	service   SubscriptionService
	presenter ResponsePresenter
	logger    Logger
}

func NewSubscriptionHandler(
	s SubscriptionService,
	p ResponsePresenter,
	l Logger,
) *SubscriptionHandler {
	return &SubscriptionHandler{
		service:   s,
		presenter: p,
		logger:    l,
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
		h.logger.Error(err.Error())
		return h.presenter.PresentError(c.Status(fiber.StatusBadRequest),
			&response.ErrorResponse{
				Message: err.Error(),
			})
	}

	if _, err := mail.ParseAddress(subscriber.Email); err != nil {
		h.logger.Error(err.Error())
		return h.presenter.PresentError(c.Status(fiber.StatusBadRequest),
			&response.ErrorResponse{
				Message: err.Error(),
			})
	}

	err := h.service.Subscribe(subscriber)
	if err != nil {
		h.logger.Error(err.Error())
		if errors.Is(err, customerror.ErrAlreadyExists) {
			return c.SendStatus(fiber.StatusConflict)
		}
		return h.presenter.PresentError(c.Status(fiber.StatusInternalServerError),
			&response.ErrorResponse{
				Message: err.Error(),
			})
	}

	return c.SendStatus(fiber.StatusOK)
}
