package logger

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"genesis-test/src/config"
	"io"
	"log"
	"net/url"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type customWriter struct {
	io.Writer
}

func (cw customWriter) Close() error {
	return nil
}

func (cw customWriter) Sync() error {
	return nil
}

const (
	debugLogKey       = "debug"
	infoLogKey        = "info"
	errorLogKey       = "error"
	rabbitmqWriterKey = "rabbitmqwriter"
)

type rabbitMQConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Exchange string
}

func configureRabbitMQConfig() rabbitMQConfig {
	envCfg := config.GetRabbit()
	return rabbitMQConfig{
		Host:     envCfg.RabbitMQHost,
		Port:     envCfg.RabbitMQPort,
		Username: envCfg.RabbitMQUser,
		Password: envCfg.RabbitMQPassword,
		Exchange: envCfg.RabbitMQExchange,
	}
}

type ZapRabbitMQLogger struct {
	buffer     *bytes.Buffer
	writer     *bufio.Writer
	logger     *ZapLogger
	channel    *amqp.Channel
	connection *amqp.Connection
	config     rabbitMQConfig
}

func NewZapRabbitMQLogger() *ZapRabbitMQLogger {
	logPath := fmt.Sprintf("%s:", rabbitmqWriterKey)
	cfg := configureRabbitMQConfig()
	connection := dialWithRabbitMQ(&cfg)
	channel := configureRabbitMQChannel(&cfg, connection)
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)

	registerLoggerSink(writer)
	logger := &ZapRabbitMQLogger{
		buffer:     &buffer,
		writer:     writer,
		logger:     NewZapLogger(logPath),
		config:     cfg,
		channel:    channel,
		connection: connection,
	}

	return logger
}

func (l *ZapRabbitMQLogger) Debug(msg string) {
	l.logger.Debug(msg)
	l.publishLog(debugLogKey)
}

func (l *ZapRabbitMQLogger) Info(msg string) {
	l.logger.Info(msg)
	l.publishLog(infoLogKey)
}

func (l *ZapRabbitMQLogger) Error(msg string) {
	l.logger.Error(msg)
	l.publishLog(errorLogKey)
}

func (l *ZapRabbitMQLogger) Close() {
	err := l.channel.Close()
	log.Println(err.Error(), "failed to close a channel")

	err = l.connection.Close()
	log.Println(err.Error(), "failed to close connection to RabbitMQ")
}

func dialWithRabbitMQ(cfg *rabbitMQConfig) *amqp.Connection {
	rabbitURL := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)

	connection, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatal(err, "failed connect to RabbitMQ")
		return nil
	}

	return connection
}

func configureRabbitMQChannel(cfg *rabbitMQConfig, connection *amqp.Connection) *amqp.Channel {
	channel, err := connection.Channel()
	if err != nil {
		log.Fatal(err, "failed to configure RabbitMQ channel")
		return nil
	}

	if err := channel.ExchangeDeclare(
		cfg.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		log.Fatal(err, "failed to declare exchange")
	}
	return channel
}

func registerLoggerSink(writer *bufio.Writer) {
	err := zap.RegisterSink(rabbitmqWriterKey, func(u *url.URL) (zap.Sink, error) {
		return customWriter{writer}, nil
	})
	if err != nil {
		log.Fatal(err, "failed to sink logger")
	}
}

func (l *ZapRabbitMQLogger) publishLog(key string) {
	_ = l.writer.Flush()
	err := l.channel.PublishWithContext(context.Background(),
		l.config.Exchange,
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        l.buffer.Bytes(),
		},
	)
	if err != nil {
		log.Fatal(err.Error(), "failed to publish a message")
	}
	l.buffer.Reset()
}
