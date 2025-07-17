package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func setupStubDBServer(status int, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Write([]byte(body))
	}))

}

func TestProxyCreate(t *testing.T) {
	stub := setupStubDBServer(http.StatusCreated, `{"id":"1","title":"T","description":"D"}`)
	defer stub.Close()
	os.Setenv("DB_SERVICE_URL", stub.URL)

	app := fiber.New()
	RegisterTaskRoutes(app)
	req := httptest.NewRequest("POST", "/create", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}
}

func TestListTasks(t *testing.T) {
	stub := setupStubDBServer(http.StatusOK, `[{"id":"1","title":"T","description":"D","created_at":"","updated_at":"","done_at":null,"completed":false}]`)
	defer stub.Close()
	os.Setenv("DB_SERVICE_URL", stub.URL)
	app := fiber.New()
	RegisterTaskRoutes(app)
	req, _ := http.NewRequest("GET", "/list", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to list tasks: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	data, _ := io.ReadAll(resp.Body)
	expected := `[{"id":"1","title":"T","description":"D","created_at":"","updated_at":"","done_at":null,"completed":false}]`
	if string(data) != expected {
		t.Errorf("Expected body %s, got %s", expected, string(data))
	}
}

func TestProxyGetByID(t *testing.T) {
	stub := setupStubDBServer(http.StatusOK, `{"id":"1","title":"T","description":"D"}`)
	defer stub.Close()
	os.Setenv("DB_SERVICE_URL", stub.URL)

	app := fiber.New()
	RegisterTaskRoutes(app)
	req := httptest.NewRequest("GET", "/task/by_id/1", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to get task by ID: %v", err)

	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestProxyUpdate(t *testing.T) {
	stub := setupStubDBServer(http.StatusOK, `{"id":"1","title":"T","description":"D"}`)
	defer stub.Close()
	os.Setenv("DB_SERVICE_URL", stub.URL)

	app := fiber.New()
	RegisterTaskRoutes(app)
	req := httptest.NewRequest("PUT", "/update/1", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to update task: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestProxyDelete(t *testing.T) {
	stub := setupStubDBServer(http.StatusOK, `{"message":"Task deleted"}`)
	defer stub.Close()
	os.Setenv("DB_SERVICE_URL", stub.URL)

	app := fiber.New()
	RegisterTaskRoutes(app)
	req := httptest.NewRequest("DELETE", "/delete/1", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to delete task: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestProxyDone(t *testing.T) {
	stub := setupStubDBServer(http.StatusOK, `{"id":"1","title":"T","description":"D"}`)
	defer stub.Close()
	os.Setenv("DB_SERVICE_URL", stub.URL)

	app := fiber.New()
	RegisterTaskRoutes(app)
	req := httptest.NewRequest("PUT", "/done/1", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to mark task as done: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}
