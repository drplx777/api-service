package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

var dbServiceURL = func() string {
	if envURL := os.Getenv("DB_SERVICE_URL"); envURL != "" {
		return envURL
	}

	return "http://db-service:8000"

}

func Post(path string, payload interface{}) (*http.Response, error) {
	var body *bytes.Reader
	if payload != nil {
		buf, _ := json.Marshal(payload)
		body = bytes.NewReader(buf)
	} else {
		body = bytes.NewReader([]byte{})
	}
	return http.Post(dbServiceURL()+path, "application/json", body)
}

func Get(path string) (*http.Response, error) {
	return http.Get(dbServiceURL() + path)
}

func Delete(path string) (*http.Response, error) {
	req, _ := http.NewRequest(http.MethodDelete, dbServiceURL()+path, nil)
	return http.DefaultClient.Do(req)
}

func Put(path string, payload interface{}) (*http.Response, error) {
	var body *bytes.Reader
	if payload != nil {
		buf, _ := json.Marshal(payload)
		body = bytes.NewReader(buf)
	} else {
		body = bytes.NewReader(nil)
	}
	req, err := http.NewRequest(http.MethodPut, dbServiceURL()+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return http.DefaultClient.Do(req)
}
