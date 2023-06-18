package utils

import (
	"genesis-test/src/config"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
)

func StartServerWithGracefulShutdown(app *fiber.App) {
	go func() {
		err := app.Listen(config.Get().ServerURL)
		if err != nil {
			log.Printf("Server unexpectedly stopped with error %s", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)
	<-stop

	if err := app.Shutdown(); err != nil {
		log.Printf("Error shutting down server %s", err)
	} else {
		log.Println("Server gracefully stopped")
	}
}
