package client

import (
	"io"
	"net/http"
	"os"
)

var dbServiceURL = func() string {
	if envURL := os.Getenv("DB_SERVICE_URL"); envURL != "" {
		return envURL
	}
	return "http://db-service:8000"
}()

func Get(path string) (*http.Response, error) {
	return http.Get(dbServiceURL + path)
}

func Post(path string, body io.Reader) (*http.Response, error) {
	return http.Post(dbServiceURL+path, "application/json", body)
}

func Put(path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPut, dbServiceURL+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return http.DefaultClient.Do(req)
}

func Delete(path string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodDelete, dbServiceURL+path, nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}
