package handler

import (
	_ "genesis-test/docs" //nolint:revive
	"genesis-test/src/app/repository"
	"genesis-test/src/app/service"

	"github.com/gofiber/fiber/v2"
	swagger "github.com/swaggo/fiber-swagger"
)

func InitRoutes(app *fiber.App) {
	repos := repository.NewRepositories()
	services := service.NewServices(repos)
	newsletterHandler := NewNewsletterHandler(services)
	exchangeHandler := NewExchangeHandler(services)

	api := app.Group("/api")

	api.Get("/rate", exchangeHandler.GetCurrencyRate)

	api.Post("/sendEmails", newsletterHandler.SendEmails)
	api.Post("/subscribe", newsletterHandler.Subscribe)

	app.Get("/swagger/*", swagger.WrapHandler)
}
