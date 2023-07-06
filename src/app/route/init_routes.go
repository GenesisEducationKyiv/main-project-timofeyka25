package route

import (
	_ "genesis-test/docs" //nolint:revive
	"genesis-test/src/app/domain"
	"genesis-test/src/app/handler"
	"genesis-test/src/app/handler/middleware"
	"genesis-test/src/app/persistence/exchange"
	"genesis-test/src/app/persistence/newsletter"
	"genesis-test/src/app/persistence/storage"
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

type persistence struct {
	sender  service.NewsletterSender
	storage service.EmailStorage
	chain   service.ExchangeChain
}

type services struct {
	exchange     handler.ExchangeService
	newsletter   handler.NewsletterService
	subscription handler.SubscriptionService
}

func createPersistence() *persistence {
	cfg := config.Get()
	smtpMailer := mailer.NewSMTPMailer(
		cfg.SMTPServer,
		cfg.SMTPPort,
		cfg.SMTPUsername,
		cfg.SMTPPassword)
	csvStorage := storage.NewCsvRepository(cfg.StorageFile)
	newsletterSender := newsletter.NewNewsletterSender(smtpMailer)
	exchangerChain := getExchangeChains()

	return &persistence{
		sender:  newsletterSender,
		storage: csvStorage,
		chain:   exchangerChain,
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

func createServices(persistence *persistence) *services {
	cfg := config.Get()
	BTCUAHPair := &domain.CurrencyPair{
		BaseCurrency:  cfg.BaseCurrency,
		QuoteCurrency: cfg.QuoteCurrency,
	}
	newExchangeService := exchangeService.NewExchangeService(
		BTCUAHPair,
		persistence.chain)
	newNewsletterService := newsletterService.NewNewsletterService(
		persistence.chain,
		persistence.storage,
		persistence.sender,
		BTCUAHPair)
	subscriptionService := subscription.NewSubscriptionService(persistence.storage)

	return &services{
		exchange:     newExchangeService,
		newsletter:   newNewsletterService,
		subscription: subscriptionService,
	}
}

func InitRoutes(app *fiber.App) {
	newPersistence := createPersistence()
	newServices := createServices(newPersistence)
	newsletterHandler := handler.NewNewsletterHandler(newServices.newsletter)
	exchangeHandler := handler.NewExchangeHandler(newServices.exchange)
	subscriptionHandler := handler.NewSubscriptionHandler(newServices.subscription)

	middleware.FiberMiddleware(app)

	api := app.Group("/api")

	api.Get("/rate", exchangeHandler.GetCurrencyRate)

	api.Post("/sendEmails", newsletterHandler.SendEmails)
	api.Post("/subscribe", subscriptionHandler.Subscribe)

	app.Get("/swagger/*", swagger.WrapHandler)
}
