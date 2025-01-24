package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func generateUniqueUsername() string {
	return fmt.Sprintf("test_user_%d", time.Now().UnixNano())
}

func TestRegister(t *testing.T) {
	username := generateUniqueUsername()
	payload := map[string]string{
		"username": username,
		"password": "test_pass1",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal payload: %v", err)
	}

	resp, err := http.Post("http://localhost:8080/auth/register", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode, "Expected status code 201 Created")
}

func TestLogin(t *testing.T) {
	username := generateUniqueUsername()
	payload := map[string]string{
		"username": username,
		"password": "test_pass1",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal payload: %v", err)
	}

	resp, err := http.Post("http://localhost:8080/auth/register", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode, "Expected status code 201 Created")

	resp, err = http.Post("http://localhost:8080/auth/login", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200 OK")
}
