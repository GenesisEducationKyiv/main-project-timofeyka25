package handler

import (
	_ "genesis-test/docs"
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

	// init base path
	api := app.Group("/api")

	// routes for GET method:
	api.Get("/rate", exchangeHandler.GetCurrencyRate)

	// Routes for POST method:
	api.Post("/sendEmails", newsletterHandler.SendEmails)
	api.Post("/subscribe", newsletterHandler.Subscribe)

	// Swagger route
	app.Get("/swagger/*", swagger.WrapHandler)
}
