package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	handler := NewHealthHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	handler.HealthCheck(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, res.StatusCode)
	}

	contentType := res.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("expected status healthy, got %v", response["status"])
	}

	if response["service"] != "predicato" {
		t.Errorf("expected service predicato, got %v", response["service"])
	}

	if _, ok := response["timestamp"]; !ok {
		t.Error("expected timestamp in response")
	}

	if _, ok := response["version"]; !ok {
		t.Error("expected version in response")
	}
}

func TestLivenessCheck(t *testing.T) {
	handler := NewHealthHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/live", nil)
	w := httptest.NewRecorder()

	handler.LivenessCheck(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, res.StatusCode)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response["status"] != "alive" {
		t.Errorf("expected status alive, got %v", response["status"])
	}
}

func TestReadinessCheckWithNilPredicato(t *testing.T) {
	handler := NewHealthHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	w := httptest.NewRecorder()

	handler.ReadinessCheck(w, req)

	res := w.Result()
	defer res.Body.Close()

	// With nil predicato, should return service unavailable
	if res.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("expected status %d, got %d", http.StatusServiceUnavailable, res.StatusCode)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response["status"] != "not_ready" {
		t.Errorf("expected status not_ready, got %v", response["status"])
	}

	checks, ok := response["checks"].(map[string]interface{})
	if !ok {
		t.Fatal("expected checks in response")
	}

	dbCheck, ok := checks["database"].(map[string]interface{})
	if !ok {
		t.Fatal("expected database check in response")
	}

	if dbCheck["status"] != "unhealthy" {
		t.Errorf("expected database status unhealthy, got %v", dbCheck["status"])
	}
}

func TestDetailedHealthCheckWithNilPredicato(t *testing.T) {
	handler := NewHealthHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/health/detailed", nil)
	w := httptest.NewRecorder()

	handler.DetailedHealthCheck(w, req)

	res := w.Result()
	defer res.Body.Close()

	// With nil predicato, should return service unavailable
	if res.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("expected status %d, got %d", http.StatusServiceUnavailable, res.StatusCode)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response["status"] != "unhealthy" {
		t.Errorf("expected status unhealthy, got %v", response["status"])
	}

	// Check build info is present
	if _, ok := response["build_info"]; !ok {
		t.Error("expected build_info in response")
	}

	// Check metrics is present
	metrics, ok := response["metrics"].(map[string]interface{})
	if !ok {
		t.Fatal("expected metrics in response")
	}

	if _, ok := metrics["response_time_ms"]; !ok {
		t.Error("expected response_time_ms in metrics")
	}
}

func TestGetSystemMetrics(t *testing.T) {
	handler := NewHealthHandler(nil)

	metrics := handler.getSystemMetrics()

	// Check that metrics are populated
	if metrics.MemoryUsage == "" {
		t.Error("expected memory_usage to be set")
	}

	if metrics.Goroutines < 1 {
		t.Errorf("expected at least 1 goroutine, got %d", metrics.Goroutines)
	}

	if metrics.StackUsage == "" {
		t.Error("expected stack_usage to be set")
	}
}

func TestWriteJSON(t *testing.T) {
	w := httptest.NewRecorder()

	data := map[string]interface{}{
		"key": "value",
		"num": 42,
	}

	writeJSON(w, http.StatusOK, data)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, res.StatusCode)
	}

	if res.Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", res.Header.Get("Content-Type"))
	}

	var response map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response["key"] != "value" {
		t.Errorf("expected key=value, got %v", response["key"])
	}

	if response["num"] != float64(42) { // JSON numbers decode as float64
		t.Errorf("expected num=42, got %v", response["num"])
	}
}
