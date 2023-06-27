package main

import (
	"genesis-test/src/app/route"
	"genesis-test/src/app/utils"
	"genesis-test/src/config"

	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
)

//	@title		Genesis Mailer
// 	@version 1.0
//	@host		localhost:8000
//	@BasePath	/api

func main() {
	app := fiber.New(config.FiberConfig())

	route.InitRoutes(app)
	utils.StartServerWithGracefulShutdown(app)
}
