package main

import (
	"api-service/internal/handler"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	app := fiber.New()
	handler.RegisterTaskRoutes(app)
	log.Printf("Starting API service on port %s", port)
	log.Fatal(app.Listen(":" + port))

}
