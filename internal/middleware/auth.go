package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func AuthMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		// Пропускаем публичные эндпоинты
		if c.Path() == "/api/login" || c.Path() == "/api/register" {
			return c.Next()
		}

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			slog.Warn("Authorization header missing")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header missing"})
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			slog.Warn("Invalid token format")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format"})
		}

		// Проверяем токен через auth-service
		req, err := http.NewRequest("GET", "http://auth-service:5000/validate?token="+token, nil)
		if err != nil {
			slog.Error("Failed to create validation request", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			slog.Error("Failed to connect to auth-service", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Authentication service unavailable"})
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return c.Status(resp.StatusCode).SendString("Invalid token")
		}

		var result struct {
			Valid    bool `json:"valid"`
			UserID   int  `json:"userID"`
			UserRole int  `json:"userRole"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			slog.Error("Failed to parse validation response", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
		}

		if !result.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		c.Locals("userID", result.UserID)
		c.Locals("userRole", result.UserRole)

		return c.Next()
	}
}
