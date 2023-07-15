package main

import (
	_ "github.com/joho/godotenv/autoload"
	rabbitmq "logger-service/src/rabbitmq"
)

func main() {
	consumer := rabbitmq.NewRabbitMQConsumer()
	consumer.LogBindingMessages()
}
