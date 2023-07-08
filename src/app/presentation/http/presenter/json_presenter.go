package presenter

import (
	"genesis-test/src/app/presentation/http"
	"genesis-test/src/app/presentation/http/response"
	"github.com/gofiber/fiber/v2"
)

type JSONPresenter struct{}

func NewJSONPresenter() http.ResponsePresenter {
	return &JSONPresenter{}
}

func (J JSONPresenter) PresentExchangeRate(c *fiber.Ctx, r *response.RateResponse) error {
	return c.JSON(r.Rate)
}

func (J JSONPresenter) PresentSendRate(c *fiber.Ctx, r *response.SendRateResponse) error {
	return c.JSON(&fiber.Map{"unsent": r.UnsentEmails})
}

func (J JSONPresenter) PresentError(c *fiber.Ctx, r *response.ErrorResponse) error {
	return c.JSON(&fiber.Map{"message": r.Message})
}
