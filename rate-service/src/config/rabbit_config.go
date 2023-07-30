package config

import (
	"os"
	"sync"
)

type RabbitConfig struct {
	RabbitMQHost     string
	RabbitMQPort     string
	RabbitMQUser     string
	RabbitMQPassword string
	RabbitMQExchange string
}

func GetRabbit() *RabbitConfig {
	var cfg RabbitConfig
	var once sync.Once
	once.Do(func() {
		cfg = RabbitConfig{
			RabbitMQHost:     os.Getenv("RABBIT_MQ_HOST"),
			RabbitMQPort:     os.Getenv("RABBIT_MQ_PORT"),
			RabbitMQUser:     os.Getenv("RABBIT_MQ_USER"),
			RabbitMQPassword: os.Getenv("RABBIT_MQ_PASSWORD"),
			RabbitMQExchange: os.Getenv("RABBIT_MQ_EXCHANGE"),
		}
	})
	return &cfg
}
