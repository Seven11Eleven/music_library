package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	log "github.com/sirupsen/logrus"
	"os"
)

func MiddlewaresSetup(server *fiber.App) {
	log.Info("Setting up middlewares...")

	file := createFileForLogs()
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
	})
	log.SetOutput(file)

	server.Use(
		cors.New(cors.Config{
			AllowMethods: "POST, GET, DELETE, PUT",
		}),
		logger.New(logger.Config{
			Output: file,
		}),
		limiter.New(limiter.Config{
			Max: 1000,
		}),
	)

	log.Debug("Middlewares have been set up.")
}

func createFileForLogs() *os.File {
	log.Info("Creating log file...")

	file, err := os.OpenFile("/var/log/app/muslib.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}

	log.Info("Log file created successfully.")
	return file
}
