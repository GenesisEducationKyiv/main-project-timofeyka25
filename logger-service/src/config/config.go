package config

import (
	"os"
	"sync"
)

type Config struct {
	RabbitMQHost      string
	RabbitMQPort      string
	RabbitMQUser      string
	RabbitMQPassword  string
	RabbitMQExchange  string
	DefaultQueueName  string
	DefaultBindingKey string
}

func Get() *Config {
	var cfg Config
	var once sync.Once
	once.Do(func() {
		cfg = Config{
			RabbitMQHost:      os.Getenv("RABBIT_MQ_HOST"),
			RabbitMQPort:      os.Getenv("RABBIT_MQ_PORT"),
			RabbitMQUser:      os.Getenv("RABBIT_MQ_USER"),
			RabbitMQPassword:  os.Getenv("RABBIT_MQ_PASSWORD"),
			RabbitMQExchange:  os.Getenv("RABBIT_MQ_EXCHANGE"),
			DefaultQueueName:  os.Getenv("DEFAULT_QUEUE_NAME"),
			DefaultBindingKey: os.Getenv("DEFAULT_BINDING_KEY"),
		}
	})
	return &cfg
}
