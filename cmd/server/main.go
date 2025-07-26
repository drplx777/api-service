package main

import (
	"api-service/internal/handler"
	"api-service/internal/middleware"
	"log"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
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
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://my-samovar-to-do-list.duckdns.org/", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Authorization"},
	}))

	app.Use(middleware.SlogLogger())
	handler.RegisterTaskRoutes(app)

	slog.Info("Service started", "port", 8080)
	slog.Warn("Low disk space", "disk", "/dev/sda1", "free_percent", 5)
	slog.Error("Database connection failed", "error", "timeout")
	log.Fatal(app.Listen(":" + port))

}
