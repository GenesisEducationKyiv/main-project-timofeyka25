package http

import (
	"genesis-test/src/app/presentation/http/response"

	"github.com/gofiber/fiber/v2"
)

type ResponsePresenter interface {
	PresentExchangeRate(c *fiber.Ctx, r *response.RateResponse) error
	PresentSendRate(c *fiber.Ctx, r *response.SendRateResponse) error
	PresentError(c *fiber.Ctx, r *response.ErrorResponse) error
}

type Logger interface {
	Info(msg string)
	Debug(msg string)
	Error(msg string)
}
