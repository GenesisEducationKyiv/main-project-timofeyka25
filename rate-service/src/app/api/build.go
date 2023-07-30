package api

import (
	exchangeServicePkg "genesis-test/src/app/application/exchange"
	newsletterServicePkg "genesis-test/src/app/application/newsletter"
	subscriptionServicePkg "genesis-test/src/app/application/subscription"
	"genesis-test/src/app/domain/model"
	"genesis-test/src/app/persistence/exchange"
	exchangeLogger "genesis-test/src/app/persistence/exchange/logger"
	"genesis-test/src/app/persistence/newsletter"
	"genesis-test/src/app/persistence/storage"
	"genesis-test/src/app/presentation/http"
	"genesis-test/src/app/presentation/http/middleware"
	"genesis-test/src/app/presentation/http/presenter"
	"genesis-test/src/app/presentation/http/route"
	"genesis-test/src/config"
	"genesis-test/src/logger"
	"genesis-test/src/pkg/mailer"

	"github.com/gofiber/fiber/v2"
)

type Logger interface {
	Info(msg string)
	Debug(msg string)
	Error(msg string)
}

var appLogger Logger //nolint:gochecknoglobals

func initLogger() {
	appLogger = logger.NewZapRabbitMQLogger()
}

type newsletterSender interface {
	MultipleSending(subscribers []string, message *model.EmailMessage) ([]string, error)
}

type csvStorage interface {
	GetAllEmails() ([]string, error)
	AddEmail(newEmail string) error
}

type exchangeProvider interface {
	GetCurrencyRate(pair *model.CurrencyPair) (*model.CurrencyRate, error)
}

type persistence struct {
	Sender    newsletterSender
	Storage   csvStorage
	Providers exchangeProvider
}

type exchangeService interface {
	GetCurrencyRate(pair *model.CurrencyPair) (float64, error)
}

type newsletterService interface {
	SendCurrencyRate() ([]string, error)
}

type subscriptionService interface {
	Subscribe(subscriber *model.Subscriber) error
}

type services struct {
	Exchange     exchangeService
	Newsletter   newsletterService
	Subscription subscriptionService
}

func createPersistence() *persistence {
	cfg := config.Get()
	smtpMailer := mailer.NewSMTPMailer(
		cfg.SMTPServer,
		cfg.SMTPPort,
		cfg.SMTPUsername,
		cfg.SMTPPassword)
	newCsvStorage := storage.NewCsvRepository(cfg.StorageFile)
	newNewsletterSender := newsletter.NewNewsletterSender(
		smtpMailer,
		appLogger,
	)
	newExchangerProviders := getExchangeProviders()

	return &persistence{
		Sender:    newNewsletterSender,
		Storage:   newCsvStorage,
		Providers: newExchangerProviders,
	}
}

func getExchangeProviders() exchangeProvider {
	newExchangeLogger := exchangeLogger.NewExchangeLogger(appLogger)

	coinbaseProvider := exchange.CoinbaseFactory{}.CreateCoinbaseFactory()
	binanceProvider := exchange.BinanceFactory{}.CreateBinanceFactory()
	btcTradeProvider := exchange.BTCTradeUAFactory{}.CreateBTCTradeUAFactory()
	coingeckoProvider := exchange.CoingeckoFactory{}.CreateCoingeckoFactory()

	coinbaseLoggingProvider := exchange.NewLoggingWrapper(coinbaseProvider, newExchangeLogger)
	binanceLoggingProvider := exchange.NewLoggingWrapper(binanceProvider, newExchangeLogger)
	btcTradeLoggingProvider := exchange.NewLoggingWrapper(btcTradeProvider, newExchangeLogger)
	coingeckoLoggingProvider := exchange.NewLoggingWrapper(coingeckoProvider, newExchangeLogger)

	coinbaseLoggingProviderNode := exchange.NewProviderNode(coinbaseLoggingProvider)
	binanceLoggingProviderNode := exchange.NewProviderNode(binanceLoggingProvider)
	btcTradeLoggingProviderNode := exchange.NewProviderNode(btcTradeLoggingProvider)
	coingeckoLoggingProviderNode := exchange.NewProviderNode(coingeckoLoggingProvider)

	coinbaseLoggingProviderNode.SetNext(binanceLoggingProviderNode)
	binanceLoggingProviderNode.SetNext(btcTradeLoggingProviderNode)
	btcTradeLoggingProviderNode.SetNext(coingeckoLoggingProviderNode)

	return coinbaseLoggingProviderNode
}

func createServices(persistence *persistence) *services {
	cfg := config.Get()
	BTCUAHPair := &model.CurrencyPair{
		BaseCurrency:  model.CurrencyFactory{}.CreateCurrency(cfg.BaseCurrency),
		QuoteCurrency: model.CurrencyFactory{}.CreateCurrency(cfg.QuoteCurrency),
	}
	newExchangeService := exchangeServicePkg.NewExchangeService(persistence.Providers)
	newNewsletterService := newsletterServicePkg.NewNewsletterService(
		persistence.Providers,
		persistence.Storage,
		persistence.Sender,
		BTCUAHPair)
	newSubscriptionService := subscriptionServicePkg.NewSubscriptionService(persistence.Storage)

	return &services{
		Exchange:     newExchangeService,
		Newsletter:   newNewsletterService,
		Subscription: newSubscriptionService,
	}
}

func BuildApp(app *fiber.App) {
	initLogger()
	newPersistence := createPersistence()
	newServices := createServices(newPersistence)
	jsonPresenter := presenter.NewJSONPresenter()

	exchangeHandler := http.NewExchangeHandler(
		newServices.Exchange,
		jsonPresenter,
		appLogger)
	newsletterHandler := http.NewNewsletterHandler(
		newServices.Newsletter,
		jsonPresenter,
		appLogger,
	)
	subscriptionHandler := http.NewSubscriptionHandler(
		newServices.Subscription,
		jsonPresenter,
		appLogger,
	)

	middleware.FiberMiddleware(app)
	route.RegisterExchangeHandler(app, exchangeHandler)
	route.RegisterNewsletterHandler(app, newsletterHandler)
	route.RegisterSubscriptionHandler(app, subscriptionHandler)
	route.RegisterSwaggerHandler(app)
}
