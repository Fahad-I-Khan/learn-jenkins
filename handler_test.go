package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHealthHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	HealthHandler(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json response: %v", err)
	}

	if resp["status"] != "ok" {
		t.Fatalf("expected status=ok, got %s", resp["status"])
	}

	if resp["message"] != "Testing CI/CD Pipeline" {
		t.Fatalf("unexpected message: %s", resp["message"])
	}
}
