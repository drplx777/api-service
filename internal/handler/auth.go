package handler

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
)

func RegisterAuthRoutes(app *fiber.App) {
	app.Post("/api/register", proxyRegister)
	app.Post("/api/login", proxyLogin)
	app.Get("/api/validate", proxyValidateToken)
	app.Get("/api/getuserbyJWT", proxyGetUser)
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
	// Берём заголовок Authorization: Bearer <token>
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Authorization header is missing")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusBadRequest).SendString("Malformed Authorization header")
	}
	token := parts[1]

	// Формируем URL с query-параметром token
	authServiceURL := "http://auth-service:5000/validate"
	u, err := url.Parse(authServiceURL)
	if err != nil {
		slog.Error("Failed to parse auth-service URL", "error", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal server error")
	}
	q := u.Query()
	q.Set("token", token)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		slog.Error("Failed to create request", "error", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal server error")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Failed to connect to auth-service", "error", err)
		return c.Status(fiber.StatusServiceUnavailable).SendString("Service unavailable")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to read response body")
	}

	c.Status(resp.StatusCode)
	return c.Send(body)
}
func proxyGetUser(c fiber.Ctx) error {
	// Берём заголовок Authorization: Bearer <token>
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Authorization header is missing")
	}

	// Создаём новый запрос к auth-service
	authServiceURL := "http://auth-service:5000/userbyJWT"
	req, err := http.NewRequest("GET", authServiceURL, nil)
	if err != nil {
		slog.Error("Failed to create request", "error", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Internal server error")
	}

	// Передаём оригинальный заголовок авторизации
	req.Header.Set("Authorization", authHeader)

	// Выполняем запрос
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Failed to connect to auth-service", "error", err)
		return c.Status(fiber.StatusServiceUnavailable).SendString("Service unavailable")
	}
	defer resp.Body.Close()

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Failed to read response body", "error", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to read response")
	}

	// Проксируем статус и тело
	c.Status(resp.StatusCode)
	// Если контент‑тайп JSON, вернём JSON
	if strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
		return c.Send(body)
	}
	// Иначе просто строка
	return c.SendString(string(body))
}
