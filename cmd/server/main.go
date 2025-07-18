package main

import (
	"api-service/internal/handler"
	"api-service/internal/middleware"
	"log"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v3"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	slog.SetDefault(logger)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	app := fiber.New()
	app.Use(middleware.SlogLogger())
	handler.RegisterTaskRoutes(app)

	slog.Info("Service started", "port", 8080)
	slog.Warn("Low disk space", "disk", "/dev/sda1", "free_percent", 5)
	slog.Error("Database connection failed", "error", "timeout")
	log.Fatal(app.Listen(":" + port))

}
