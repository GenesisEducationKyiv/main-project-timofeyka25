package route

import (
	_ "genesis-test/docs" //nolint:revive
	"genesis-test/src/app/api"
	"genesis-test/src/app/presentation/http"
	"genesis-test/src/app/presentation/http/middleware"

	"github.com/gofiber/fiber/v2"
	swagger "github.com/swaggo/fiber-swagger"
)

func InitRoutes(app *fiber.App) {
	newPersistence := api.CreatePersistence()
	newServices := api.CreateServices(newPersistence)
	newsletterHandler := http.NewNewsletterHandler(newServices.Newsletter)
	exchangeHandler := http.NewExchangeHandler(newServices.Exchange)
	subscriptionHandler := http.NewSubscriptionHandler(newServices.Subscription)

	middleware.FiberMiddleware(app)

	apiGroup := app.Group("/api")

	apiGroup.Get("/rate", exchangeHandler.GetCurrencyRate)

	apiGroup.Post("/sendEmails", newsletterHandler.SendEmails)
	apiGroup.Post("/subscribe", subscriptionHandler.Subscribe)

	app.Get("/swagger/*", swagger.WrapHandler)
}
