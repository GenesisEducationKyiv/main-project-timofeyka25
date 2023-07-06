package main

import (
	"flag"
	"genesis-test/src/app/route"
	"genesis-test/src/app/utils"
	"genesis-test/src/config"
	"log"

	"github.com/joho/godotenv"

	"github.com/gofiber/fiber/v2"
)

//	@title		Genesis Mailer
// 	@version 1.0
//	@host		localhost:8000
//	@BasePath	/api

func main() {
	if err := loadEnv(); err != nil {
		log.Fatalf("failed to load env file")
	}
	app := fiber.New(config.FiberConfig())

	route.InitRoutes(app)
	utils.StartServerWithGracefulShutdown(app)
}

func loadEnv() error {
	testFlag := flag.Bool("test", false, "")
	flag.Parse()
	var envFile string
	switch *testFlag {
	case true:
		envFile = "test.env"
	default:
		envFile = ".env"
	}

	return godotenv.Load(envFile)
}
