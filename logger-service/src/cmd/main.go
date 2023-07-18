package main

import (
	"logger-service/src/rabbitmq"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	consumer := rabbitmq.NewRabbitMQConsumer()
	consumer.LogBindingMessages()
}
