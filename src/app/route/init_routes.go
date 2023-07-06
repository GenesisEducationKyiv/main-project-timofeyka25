package route

import (
	_ "genesis-test/docs" //nolint:revive
	"genesis-test/src/app/domain"
	"genesis-test/src/app/handler"
	"genesis-test/src/app/handler/middleware"
	"genesis-test/src/app/repository/exchange"
	"genesis-test/src/app/repository/newsletter"
	"genesis-test/src/app/repository/storage"
	"genesis-test/src/app/service"
	"genesis-test/src/config"
	"genesis-test/src/pkg/mailer"

	"github.com/gofiber/fiber/v2"
	swagger "github.com/swaggo/fiber-swagger"
)

func createRepositories() *service.Repositories {
	cfg := config.Get()
	smtpMailer := mailer.NewSMTPMailer(
		cfg.SMTPServer,
		cfg.SMTPPort,
		cfg.SMTPUsername,
		cfg.SMTPPassword)
	csvStorage := storage.NewCsvStorage(cfg.StorageFile)
	newsletterRepo := newsletter.NewNewsletterRepository(smtpMailer)
	exchangeRepo := exchange.NewExchangeCoinbaseRepository(cfg.CryptoAPIFormatURL)

	return &service.Repositories{
		Newsletter: newsletterRepo,
		Storage:    csvStorage,
		Exchange:   exchangeRepo,
	}
}

func createServices(repos *service.Repositories) *handler.Services {
	cfg := config.Get()
	BTCUAHPair := &domain.CurrencyPair{
		BaseCurrency:  cfg.BaseCurrency,
		QuoteCurrency: cfg.QuoteCurrency,
	}
	exchangeService := service.NewExchangeService(BTCUAHPair, repos.Exchange)
	newsletterService := service.NewNewsletterService(repos, BTCUAHPair)
	subscriptionService := service.NewSubscriptionService(repos.Storage)

	return &handler.Services{
		Subscription: subscriptionService,
		Newsletter:   newsletterService,
		Exchange:     exchangeService,
	}
}

func InitRoutes(app *fiber.App) {
	repos := createRepositories()
	services := createServices(repos)
	newsletterHandler := handler.NewNewsletterHandler(services.Newsletter)
	exchangeHandler := handler.NewExchangeHandler(services.Exchange)
	subscriptionHandler := handler.NewSubscriptionHandler(services.Subscription)

	middleware.FiberMiddleware(app)

	api := app.Group("/api")

	api.Get("/rate", exchangeHandler.GetCurrencyRate)

	api.Post("/sendEmails", newsletterHandler.SendEmails)
	api.Post("/subscribe", subscriptionHandler.Subscribe)

	app.Get("/swagger/*", swagger.WrapHandler)
}
