package main

import (
	"flag"
	"fmt"
	"mixtureai/config"
	"mixtureai/log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	configFile := flag.String("config", "config.yaml", "the path of config.yaml")
	flag.Parse()
	config.Load(*configFile)

	app := fiber.New(fiber.Config{
		BodyLimit: 100 * 1024 * 1024, // 限制100MB
	})

	app.Use(logger.New(logger.Config{
		// For more options, see the Config section
		Format: "${pid} ${ip} ${status} - ${method} ${path}",
		Output: log.FiberWriter{},
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("MixtureAI")
	})

	app.Listen(fmt.Sprintf(":%v", config.C.Port))
}
