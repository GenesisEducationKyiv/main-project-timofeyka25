package route

import (
	_ "genesis-test/docs" //nolint:revive
	"genesis-test/src/app/api"
	"genesis-test/src/app/presentation/http"
	"genesis-test/src/app/presentation/http/middleware"
	"genesis-test/src/app/presentation/http/presenter"

	"github.com/gofiber/fiber/v2"
	swagger "github.com/swaggo/fiber-swagger"
)

func InitRoutes(app *fiber.App) {
	newPersistence := api.CreatePersistence()
	newServices := api.CreateServices(newPersistence)

	jsonPresenter := presenter.NewJSONPresenter()
	newsletterHandler := http.NewNewsletterHandler(newServices.Newsletter, jsonPresenter)
	exchangeHandler := http.NewExchangeHandler(newServices.Exchange, jsonPresenter)
	subscriptionHandler := http.NewSubscriptionHandler(newServices.Subscription, jsonPresenter)

	middleware.FiberMiddleware(app)

	apiGroup := app.Group("/api")

	apiGroup.Get("/rate", exchangeHandler.GetCurrencyRate)

	apiGroup.Post("/sendEmails", newsletterHandler.SendEmails)
	apiGroup.Post("/subscribe", subscriptionHandler.Subscribe)

	app.Get("/swagger/*", swagger.WrapHandler)
}
