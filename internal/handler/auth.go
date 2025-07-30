package handler

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

func RegisterAuthRoutes(app *fiber.App) {
	app.Post("/api/register", proxyRegister)
	app.Post("/api/login", proxyLogin)
	app.Get("/api/validate", proxyValidateToken)
}

func proxyRegister(c fiber.Ctx) error {
	// Используем прямой URL к auth-service
	authServiceURL := "http://auth-service:5000/register"

	// Создаем новый HTTP запрос
	req, err := http.NewRequest("POST", authServiceURL, bytes.NewReader(c.Body()))
	if err != nil {
		slog.Error("Failed to create request", "error", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal server error")
	}

	// Копируем заголовки
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Failed to connect to auth-service", "error", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Service unavailable")
	}
	defer resp.Body.Close()

	// Возвращаем ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to read response body")
	}

	c.Status(resp.StatusCode)
	return c.Send(body)
}

func proxyLogin(c fiber.Ctx) error {
	// Аналогично для входа
	authServiceURL := "http://auth-service:5000/login"

	req, err := http.NewRequest("POST", authServiceURL, bytes.NewReader(c.Body()))
	if err != nil {
		slog.Error("Failed to create request", "error", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal server error")
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
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
func proxyValidateToken(c fiber.Ctx) error {

	authServiceURL := "http://auth-service:5000/validate"
	req, err := http.NewRequest("GET", authServiceURL, nil)
	if err != nil {
		slog.Error("Failed to create request", "error", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal server error")
	}

	req.Header.Set("Authorization", c.Get("Authorization"))

	client := &http.Client{}
	resp, err := client.Do(req)
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
