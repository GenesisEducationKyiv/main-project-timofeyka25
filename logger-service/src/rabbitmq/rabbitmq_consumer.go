package rabbit_mq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"logger-service/src/config"
)

type rabbitMQConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Exchange string
}

type RabbitMQConsumer struct {
	config     rabbitMQConfig
	channel    *amqp.Channel
	connection *amqp.Connection
	queueName  string
	bindingKey string
}

func NewRabbitMQConsumer() *RabbitMQConsumer {
	cfg := configureRabbitMQConfig()
	conn := dialWithRabbitMQ(&cfg)
	channel := configureRabbitMQChannel(&cfg, conn)

	consumer := &RabbitMQConsumer{
		config:     cfg,
		channel:    channel,
		connection: conn,
		queueName:  config.Get().DefaultQueueName,
		bindingKey: config.Get().DefaultBindingKey,
	}
	consumer.configureQueue()

	return consumer
}

func (r *RabbitMQConsumer) LogBindingMessages() {
	messages := r.consumeMessages()
	r.logMessages(messages)
}

func configureRabbitMQConfig() rabbitMQConfig {
	envCfg := config.Get()
	return rabbitMQConfig{
		Host:     envCfg.RabbitMQHost,
		Port:     envCfg.RabbitMQPort,
		Username: envCfg.RabbitMQUser,
		Password: envCfg.RabbitMQPassword,
		Exchange: envCfg.RabbitMQExchange,
	}
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

func configureRabbitMQChannel(cfg *rabbitMQConfig, conn *amqp.Connection) *amqp.Channel {
	channel, err := conn.Channel()
	if err != nil {
		log.Fatalln(err.Error(), "failed to create channel")
	}

	if err = channel.ExchangeDeclare(
		cfg.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		log.Fatalln(err.Error(), "failed to declare exchange")
	}

	return channel
}

func (r *RabbitMQConsumer) configureQueue() {
	if _, err := r.channel.QueueDeclare(
		r.queueName,
		false,
		false,
		true,
		false,
		nil,
	); err != nil {
		log.Fatalln(err.Error(), "failed to declare queue")
	}

	if err := r.channel.QueueBind(
		r.queueName,
		r.bindingKey,
		r.config.Exchange,
		false,
		nil,
	); err != nil {
		log.Fatalln(err.Error(), "failed to bind queue")
	}
}

func (r *RabbitMQConsumer) consumeMessages() <-chan amqp.Delivery {
	messages, err := r.channel.Consume(
		r.queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalln(err.Error(), "failed to consume channel")
	}

	return messages
}

func (r *RabbitMQConsumer) logMessages(messages <-chan amqp.Delivery) {
	var forever chan struct{}

	go func() {
		for d := range messages {
			log.Printf("%s", d.Body)
		}
	}()

	log.Printf("Start listening log queue")
	<-forever
}
