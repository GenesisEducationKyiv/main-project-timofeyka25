package route

import (
	_ "genesis-test/docs" //nolint:revive

	"github.com/gofiber/fiber/v2"
	swagger "github.com/swaggo/fiber-swagger"
)

type NewsletterHandler interface {
	SendEmails(c *fiber.Ctx) error
}

type SubscriptionHandler interface {
	Subscribe(c *fiber.Ctx) error
}

type ExchangeHandler interface {
	GetCurrencyRate(c *fiber.Ctx) error
}

func RegisterExchangeHandler(app *fiber.App, h ExchangeHandler) {
	app.Get("/api/rate", h.GetCurrencyRate)
}

func RegisterNewsletterHandler(app *fiber.App, h NewsletterHandler) {
	app.Post("/api/sendEmails", h.SendEmails)
}

func RegisterSubscriptionHandler(app *fiber.App, h SubscriptionHandler) {
	app.Post("/api/subscribe", h.Subscribe)
}

func RegisterSwaggerHandler(app *fiber.App) {
	app.Get("/swagger/*", swagger.WrapHandler)
}
