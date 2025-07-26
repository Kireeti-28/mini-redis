package main

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/kireeti-28/mini-redis/internal/storage"
)

func setupTestServer(t *testing.T) (*httptest.Server, func()) {
	file, err := os.CreateTemp("", "test_api_*.log")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	storage, err := storage.NewStorage(file.Name())
	if err != nil {
		t.Fatalf("failed to create storage: %v", err)
	}

	cfg := &apiConfig{
		store:  storage,
		logger: slog.Default(),
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /kv/{key}", cfg.getHandler)
	mux.HandleFunc("POST /kv/{key}", cfg.setHandler)
	mux.HandleFunc("DELETE /kv/{key}", cfg.deleteHandler)

	server := httptest.NewServer(mux)

	cleanup := func() {
		server.Close()
		file.Close()
		storage.Close()

		os.Remove(file.Name())
	}

	return server, cleanup
}

func TestStorageApiSetAndGet(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	resp, err := http.Post(server.URL+"/kv/name", "text/plain", strings.NewReader("void"))
	if err != nil {
		t.Fatalf("failed to POST: %v", err)
	}

	resp.Body.Close()

	resp, err = http.Get(server.URL + "/kv/name")
	if err != nil {
		t.Fatalf("failed to GET: %v", err)
	}
	data, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(data), "void") {
		t.Errorf("expected to contain %s but got %s", "void", string(data))
	}

	resp.Body.Close()
}

func TestStorageApiGetNonExist(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	resp, err := http.Get(server.URL + "/kv/name")
	if err != nil {
		t.Errorf("failed to GET: %v", err)
	}

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status code %v but got %v", http.StatusInternalServerError, resp.StatusCode)
	}

	resp.Body.Close()
}

func TestStorageApiDelete(t *testing.T) {
	server, cleanup := setupTestServer(t)
	defer cleanup()

	resp, err := http.Post(server.URL+"/kv/name", "text/plain", strings.NewReader("void"))
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("failed to POST: %v", err)
	}

	req, err := http.NewRequest(http.MethodDelete, server.URL+"/kv/name", nil)
	if err != nil {
		t.Fatalf("failed to make request %v", err)
	}

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to DELETE: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status code %v but got %v", http.StatusNoContent, resp.StatusCode)
	}

	resp, err = http.Get(server.URL + "/kv/name")
	if err != nil {
		t.Fatalf("failed to GET: %v", err)
	}

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status code %v but got %v", http.StatusInternalServerError, resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "key name does not exist") {
		t.Errorf("expected to contain 'key name does not exist' but got '%v'", string(body))
	}

	resp.Body.Close()
}
