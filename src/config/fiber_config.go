package config

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
	"time"
)

func FiberConfig() fiber.Config {
	readTimeout, _ := strconv.Atoi(cfg.ServerReadTimeout)

	return fiber.Config{
		ReadTimeout: time.Second * time.Duration(readTimeout),
	}
}
