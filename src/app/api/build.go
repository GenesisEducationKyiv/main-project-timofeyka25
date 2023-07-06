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
	exchangerChain := getExchangeChains()

	return &application.Persistence{
		Sender:    newsletterSender,
		Storage:   csvStorage,
		Providers: exchangerChain,
	}
}

func getExchangeChains() application.ExchangeChain {
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

func CreateServices(persistence *application.Persistence) *domain.Services {
	cfg := config.Get()
	BTCUAHPair := &model.CurrencyPair{
		BaseCurrency:  model.CurrencyFactory{}.CreateCurrency(cfg.BaseCurrency),
		QuoteCurrency: model.CurrencyFactory{}.CreateCurrency(cfg.QuoteCurrency),
	}
	newExchangeService := exchangeService.NewExchangeService(
		BTCUAHPair,
		persistence.Providers)
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
