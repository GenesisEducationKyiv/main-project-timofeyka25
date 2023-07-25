package route

import "github.com/gofiber/fiber/v2"

type NewsletterHandler interface {
	SendEmails(c *fiber.Ctx) error
}

type SubscriptionHandler interface {
	Subscribe(c *fiber.Ctx) error
}

type ExchangeHandler interface {
	GetCurrencyRate(c *fiber.Ctx) error
}
