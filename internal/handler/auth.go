package handler

import (
	"api-service/internal/client"
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	Role  int    `json:"roleID"`
}

func RegisterAuthRoutes(app *fiber.App) {
	app.Post("/api/register", registerHandler)
	app.Post("/api/login", loginHandler)
}

func registerHandler(c fiber.Ctx) error {
	resp, err := client.Post("/user/register", bytes.NewReader(c.Body()))
	if err != nil {
		slog.Error("Registration failed", "error", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Registration failed")
	}
	defer resp.Body.Close()

	return forwardResponse(c, resp)
}

func loginHandler(c fiber.Ctx) error {
	var req AuthRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	jsonBody, _ := json.Marshal(req)
	resp, err := client.Post("/user/login", bytes.NewReader(jsonBody))
	if err != nil {
		slog.Error("Login failed", "error", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Login failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return c.Status(resp.StatusCode).Send(body)
	}

	var user UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		slog.Error("Failed to parse user data", "error", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to parse user data")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		slog.Error("Token generation failed", "error", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Token generation failed")
	}

	return c.JSON(fiber.Map{
		"token": tokenString,
		"user":  user,
	})
}

func forwardResponse(c fiber.Ctx, resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Failed to read response body", "error", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to read response body")
	}
	c.Status(resp.StatusCode)
	return c.Send(body)
}
