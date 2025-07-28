package handler

import (
	"api-service/internal/client"
	"bytes"
	"io"
	"log/slog"

	"github.com/gofiber/fiber/v3"
)

func RegisterAuthRoutes(app *fiber.App) {
	app.Post("/api/register", proxyRegister)
	app.Post("/api/login", proxyLogin)
}

func proxyRegister(c fiber.Ctx) error {
	// Проксируем запрос в auth-service
	resp, err := client.Post("http://auth-service:5000/register", bytes.NewReader(c.Body()))
	if err != nil {
		slog.Error("Failed to connect to auth-service", "error", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Service unavailable")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to read response body")
	}
	c.Status(resp.StatusCode)
	return c.Send(body)
}

func proxyLogin(c fiber.Ctx) error {
	resp, err := client.Post("http://auth-service:5000/login", bytes.NewReader(c.Body()))
	if err != nil {
		slog.Error("Failed to connect to auth-service", "error", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Service unavailable")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to read response body")
	}
	c.Status(resp.StatusCode)
	return c.Send(body)
}
