package config

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func FiberConfig() fiber.Config {
	readTimeout, _ := strconv.Atoi(Get().ServerReadTimeout)

	return fiber.Config{
		ReadTimeout: time.Second * time.Duration(readTimeout),
	}
}
