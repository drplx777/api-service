package handler

import (
	"api-service/internal/client"
	"bytes"
	"io"

	"github.com/gofiber/fiber/v3"
)

func RegisterTaskRoutes(app *fiber.App) {
	app.Get("/list", listTasks)
	app.Post("/create", createTask)
	app.Get("/task/by_id/:id", getTaskByID)
	app.Put("/update/:id", updateTask)
	app.Delete("/delete/:id", deleteTask)
}

func createTask(c fiber.Ctx) error {
	reader := bytes.NewReader(c.Body())
	resp, err := client.Post("/create", reader)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to create task")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to read response body")
	}
	c.Status(resp.StatusCode)
	return c.Send(body)
}

func listTasks(c fiber.Ctx) error {
	resp, err := client.Get("/list")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to list tasks")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to read response body")
	}
	c.Status(resp.StatusCode)
	return c.Send(body)
}
func getTaskByID(c fiber.Ctx) error {
	id := c.Params("id")
	resp, err := client.Get("/task/by_id?id=" + id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get task")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to read response body")
	}
	c.Status(resp.StatusCode)
	return c.Send(body)
}
func updateTask(c fiber.Ctx) error {
	reader := bytes.NewReader(c.Body())
	id := c.Params("id")
	resp, err := client.Put("/update?id="+id, reader)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to update task")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to read response body")
	}
	c.Status(resp.StatusCode)
	return c.Send(body)
}
func deleteTask(c fiber.Ctx) error {
	id := c.Params("id")
	resp, err := client.Delete("/delete?id=" + id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete task")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to read response body")
	}
	c.Status(resp.StatusCode)
	return c.Send(body)
}
