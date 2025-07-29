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

	handler.RegisterAuthRoutes(app)

	app.Use(middleware.AuthMiddleware())

	handler.RegisterTaskRoutes(app)

	slog.Info("Service started", "port", port)
	log.Fatal(app.Listen(":" + port))
}
