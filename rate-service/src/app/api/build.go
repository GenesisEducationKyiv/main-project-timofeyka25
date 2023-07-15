package api

import (
	"genesis-test/src/app/application"
	exchangeService "genesis-test/src/app/application/exchange"
	newsletterService "genesis-test/src/app/application/newsletter"
	subscriptionService "genesis-test/src/app/application/subscription"
	"genesis-test/src/app/domain"
	"genesis-test/src/app/domain/model"
	"genesis-test/src/app/persistence/exchange"
	"genesis-test/src/app/persistence/newsletter"
	"genesis-test/src/app/persistence/storage"
	"genesis-test/src/config"
	"genesis-test/src/logger"
	"genesis-test/src/pkg/mailer"
)

func CreatePersistence() *application.Persistence {
	cfg := config.Get()
	smtpMailer := mailer.NewSMTPMailer(
		cfg.SMTPServer,
		cfg.SMTPPort,
		cfg.SMTPUsername,
		cfg.SMTPPassword)
	csvStorage := storage.NewCsvRepository(cfg.StorageFile)
	newsletterSender := newsletter.NewNewsletterSender(smtpMailer)
	exchangerProvider := getExchangeProviders()

	return &application.Persistence{
		Sender:    newsletterSender,
		Storage:   csvStorage,
		Providers: exchangerProvider,
	}
}

func getExchangeProviders() application.ExchangeProvider {
	newLogger := logger.NewZapLogger(config.Get().LogPath)
	exchangeLogger := exchange.NewExchangeLogger(newLogger)

	coinbaseProvider := exchange.CoinbaseFactory{}.CreateCoinbaseFactory()
	binanceProvider := exchange.BinanceFactory{}.CreateBinanceFactory()
	btcTradeProvider := exchange.BTCTradeUAFactory{}.CreateBTCTradeUAFactory()
	coingeckoProvider := exchange.CoingeckoFactory{}.CreateCoingeckoFactory()

	coinbaseLoggingProvider := exchange.NewLoggingWrapper(coinbaseProvider, exchangeLogger)
	binanceLoggingProvider := exchange.NewLoggingWrapper(binanceProvider, exchangeLogger)
	btcTradeLoggingProvider := exchange.NewLoggingWrapper(btcTradeProvider, exchangeLogger)
	coingeckoLoggingProvider := exchange.NewLoggingWrapper(coingeckoProvider, exchangeLogger)

	coinbaseLoggingProviderNode := exchange.NewProviderNode(coinbaseLoggingProvider)
	binanceLoggingProviderNode := exchange.NewProviderNode(binanceLoggingProvider)
	btcTradeLoggingProviderNode := exchange.NewProviderNode(btcTradeLoggingProvider)
	coingeckoLoggingProviderNode := exchange.NewProviderNode(coingeckoLoggingProvider)

	coinbaseLoggingProviderNode.SetNext(binanceLoggingProviderNode)
	binanceLoggingProviderNode.SetNext(btcTradeLoggingProviderNode)
	btcTradeLoggingProviderNode.SetNext(coingeckoLoggingProviderNode)

	return coinbaseLoggingProviderNode
}

func CreateServices(persistence *application.Persistence) *domain.Services {
	cfg := config.Get()
	BTCUAHPair := &model.CurrencyPair{
		BaseCurrency:  model.CurrencyFactory{}.CreateCurrency(cfg.BaseCurrency),
		QuoteCurrency: model.CurrencyFactory{}.CreateCurrency(cfg.QuoteCurrency),
	}
	newExchangeService := exchangeService.NewExchangeService(persistence.Providers)
	newNewsletterService := newsletterService.NewNewsletterService(
		persistence.Providers,
		persistence.Storage,
		persistence.Sender,
		BTCUAHPair)
	newSubscriptionService := subscriptionService.NewSubscriptionService(persistence.Storage)

	return &domain.Services{
		Exchange:     newExchangeService,
		Newsletter:   newNewsletterService,
		Subscription: newSubscriptionService,
	}
}
