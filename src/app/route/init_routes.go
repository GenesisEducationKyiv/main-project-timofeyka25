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
	exchangeService "genesis-test/src/app/service/exchange"
	newsletterService "genesis-test/src/app/service/newsletter"
	"genesis-test/src/app/service/subscription"
	"genesis-test/src/config"
	"genesis-test/src/logger"
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
	exchangerProvider := getExchangeChains()

	return &service.Repositories{
		Newsletter: newsletterRepo,
		Storage:    csvStorage,
		Exchange:   exchangerProvider,
	}
}

func getExchangeChains() service.ExchangeChain {
	newLogger := logger.NewZapLogger(config.Get().LogPath)
	exchangeLogger := exchange.NewExchangeLogger(newLogger)

	coinbaseChain := exchange.CoinbaseFactory{}.CreateCoinbaseFactory()
	binanceChain := exchange.BinanceFactory{}.CreateBinanceFactory()
	btcTradeChain := exchange.BTCTradeUAFactory{}.CreateBTCTradeUAFactory()
	coingeckoChain := exchange.CoingeckoFactory{}.CreateCoingeckoFactory()

	coinbaseLoggingChain := exchange.NewLoggingWrapper(coinbaseChain, exchangeLogger)
	binanceLoggingChain := exchange.NewLoggingWrapper(binanceChain, exchangeLogger)
	btcTradeLoggingChain := exchange.NewLoggingWrapper(btcTradeChain, exchangeLogger)
	coingeckoLoggingChain := exchange.NewLoggingWrapper(coingeckoChain, exchangeLogger)

	coinbaseLoggingChain.SetNext(binanceLoggingChain)
	binanceLoggingChain.SetNext(btcTradeLoggingChain)
	btcTradeLoggingChain.SetNext(coingeckoLoggingChain)

	return coinbaseLoggingChain
}

func createServices(repos *service.Repositories) *handler.Services {
	cfg := config.Get()
	BTCUAHPair := &domain.CurrencyPair{
		BaseCurrency:  cfg.BaseCurrency,
		QuoteCurrency: cfg.QuoteCurrency,
	}
	newExchangeService := exchangeService.NewExchangeService(BTCUAHPair, repos.Exchange)
	newNewsletterService := newsletterService.NewNewsletterService(repos, BTCUAHPair)
	subscriptionService := subscription.NewSubscriptionService(repos.Storage)

	return &handler.Services{
		Subscription: subscriptionService,
		Newsletter:   newNewsletterService,
		Exchange:     newExchangeService,
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
