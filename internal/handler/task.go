package handler

import (
	"api-service/internal/client"
	"bytes"
	"io"
	"time"

	"github.com/gofiber/fiber/v3"
)

func RegisterTaskRoutes(app *fiber.App) {
	app.Get("/list", listTasks)
	app.Post("/create", createTask)
	app.Get("/task/by_id/:id", getTaskByID)
	app.Put("/update/:id", updateTask)
	app.Delete("/delete/:id", deleteTask)
	app.Put("/done/:id", doneTask)
	app.Get("/tasklist", mockTasks) //ВРЕМЕННО

}

// WARNING ВРЕМЕННО
type Task struct {
	ID            string     `json:"id"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	Status        string     `json:"status"`
	ReporterID    string     `json:"reporterId"`
	AssignerID    *string    `json:"assignerId,omitempty"`
	ReviewerID    *string    `json:"reviewerId,omitempty"`
	ApproverID    string     `json:"approverId"`
	ApproveStatus string     `json:"approveStatus"`
	CreatedAt     time.Time  `json:"createdAt"`
	StartedAt     *time.Time `json:"startedAt,omitempty"`
	CompletedAt   *time.Time `json:"completedAt,omitempty"`
	Deadline      string     `json:"deadline"`
	DashboardID   string     `json:"dashboardId"`
	BlockedBy     []string   `json:"blockedBy"`
}

// WARNING ВРЕМЕННО
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
func doneTask(c fiber.Ctx) error {
	id := c.Params("id")
	resp, err := client.Put("/done?id="+id, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to mark task as done")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to read response body")
	}
	c.Status(resp.StatusCode)
	return c.Send(body)
}

// WARNING ВРЕМЕННО
// WARNING ВРЕМЕННО
// WARNING ВРЕМЕННО
// WARNING ВРЕМЕННО
// WARNING ВРЕМЕННО
func mockTasks(c fiber.Ctx) error {
	assigner2 := "user-2"
	assigner3 := "user-3"
	assigner4 := "user-4"

	reviewer1 := "user-1"
	reviewer4 := "user-4"

	t1CreatedAt, _ := time.Parse(time.RFC3339, "2025-07-20T09:00:00Z")
	t1StartedAt, _ := time.Parse(time.RFC3339, "2025-07-21T10:00:00Z")

	t2CreatedAt, _ := time.Parse(time.RFC3339, "2025-07-19T14:00:00Z")

	t3CreatedAt, _ := time.Parse(time.RFC3339, "2025-07-15T08:00:00Z")
	t3StartedAt, _ := time.Parse(time.RFC3339, "2025-07-16T09:30:00Z")
	t3CompletedAt, _ := time.Parse(time.RFC3339, "2025-07-20T17:00:00Z")

	t4CreatedAt, _ := time.Parse(time.RFC3339, "2025-07-21T10:00:00Z")

	t5CreatedAt, _ := time.Parse(time.RFC3339, "2025-07-22T12:00:00Z")
	t5StartedAt, _ := time.Parse(time.RFC3339, "2025-07-23T13:00:00Z")

	t6CreatedAt, _ := time.Parse(time.RFC3339, "2025-07-15T10:00:00Z")
	t6CompletedAt, _ := time.Parse(time.RFC3339, "2025-07-18T18:00:00Z")

	tasks := []Task{
		{
			ID:            "1",
			Title:         "Добавить авторизацию",
			Description:   "Регистрация, логин, защита роутов",
			Status:        "in-progress",
			ReporterID:    "user-1",
			AssignerID:    &assigner2,
			ReviewerID:    &reviewer1,
			ApproverID:    "user-1",
			ApproveStatus: "approved",
			CreatedAt:     t1CreatedAt,
			StartedAt:     &t1StartedAt,
			CompletedAt:   nil,
			Deadline:      "2025-07-30",
			DashboardID:   "dash-1",
			BlockedBy:     []string{},
		},
		{
			ID:            "2",
			Title:         "Создать UI главного дашборда",
			Description:   "Компоненты, таблицы, графики",
			Status:        "to-do",
			ReporterID:    "user-2",
			AssignerID:    nil,
			ReviewerID:    nil,
			ApproverID:    "user-2",
			ApproveStatus: "need-approval",
			CreatedAt:     t2CreatedAt,
			StartedAt:     nil,
			CompletedAt:   nil,
			Deadline:      "2025-08-01",
			DashboardID:   "dash-1",
			BlockedBy:     []string{"1"},
		},
		{
			ID:            "3",
			Title:         "Настроить CI/CD",
			Description:   "GitHub Actions + Vercel",
			Status:        "done",
			ReporterID:    "user-1",
			AssignerID:    &assigner3,
			ReviewerID:    &reviewer1,
			ApproverID:    "user-4",
			ApproveStatus: "approved",
			CreatedAt:     t3CreatedAt,
			StartedAt:     &t3StartedAt,
			CompletedAt:   &t3CompletedAt,
			Deadline:      "2025-07-25",
			DashboardID:   "dash-2",
			BlockedBy:     []string{},
		},
		{
			ID:            "4",
			Title:         "Рефакторинг компонентов",
			Description:   "Разделение логики и UI",
			Status:        "blocked",
			ReporterID:    "user-2",
			AssignerID:    &assigner4,
			ReviewerID:    nil,
			ApproverID:    "user-3",
			ApproveStatus: "need-approval",
			CreatedAt:     t4CreatedAt,
			StartedAt:     nil,
			CompletedAt:   nil,
			Deadline:      "2025-07-29",
			DashboardID:   "dash-2",
			BlockedBy:     []string{"1", "3"},
		},
		{
			ID:            "5",
			Title:         "Интеграция с беком",
			Description:   "REST API, fetch, hook-логика",
			Status:        "review",
			ReporterID:    "user-3",
			AssignerID:    &assigner2,
			ReviewerID:    &reviewer4,
			ApproverID:    "user-3",
			ApproveStatus: "approval",
			CreatedAt:     t5CreatedAt,
			StartedAt:     &t5StartedAt,
			CompletedAt:   nil,
			Deadline:      "2025-08-02",
			DashboardID:   "dash-1",
			BlockedBy:     []string{},
		},
		{
			ID:            "6",
			Title:         "Создать страницу профиля",
			Description:   "Имя, логин, роль, смена пароля",
			Status:        "canceled",
			ReporterID:    "user-1",
			AssignerID:    nil,
			ReviewerID:    nil,
			ApproverID:    "user-1",
			ApproveStatus: "rejected",
			CreatedAt:     t6CreatedAt,
			StartedAt:     nil,
			CompletedAt:   &t6CompletedAt,
			Deadline:      "2025-08-05",
			DashboardID:   "dash-5",
			BlockedBy:     []string{},
		},
	}

	return c.JSON(tasks)
}

// WARNING ВРЕМЕННО// WARNING ВРЕМЕННО// WARNING ВРЕМЕННО// WARNING ВРЕМЕННО// WARNING ВРЕМЕННО
