package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/soundprediction/predicato/pkg/server/dto"
)

func TestGenerateProcessID(t *testing.T) {
	// Generate multiple process IDs and ensure they're unique
	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id := generateProcessID()
		if id == "" {
			t.Error("generateProcessID returned empty string")
		}
		if ids[id] {
			t.Errorf("generateProcessID returned duplicate ID: %s", id)
		}
		ids[id] = true

		// Check format
		if len(id) < 5 || id[:5] != "proc_" {
			t.Errorf("generateProcessID returned invalid format: %s", id)
		}
	}
}

func TestWriteErrorJSON(t *testing.T) {
	w := httptest.NewRecorder()

	writeErrorJSON(w, http.StatusBadRequest, "invalid_request", "test error message")

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, res.StatusCode)
	}

	if res.Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", res.Header.Get("Content-Type"))
	}

	var response dto.ErrorResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Error != "invalid_request" {
		t.Errorf("expected error 'invalid_request', got %s", response.Error)
	}
	if response.Message != "test error message" {
		t.Errorf("expected message 'test error message', got %s", response.Message)
	}
}

func TestAddMessagesValidation(t *testing.T) {
	handler := NewIngestHandler(nil)

	tests := []struct {
		name           string
		body           interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "invalid JSON",
			body:           "not json",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid_request",
		},
		{
			name: "missing group_id",
			body: dto.AddMessagesRequest{
				GroupID:  "",
				Messages: []dto.Message{{Role: "user", Content: "test"}},
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid_request",
		},
		{
			name: "empty messages",
			body: dto.AddMessagesRequest{
				GroupID:  "test-group",
				Messages: []dto.Message{},
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid_request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			var err error
			if s, ok := tt.body.(string); ok {
				body = []byte(s)
			} else {
				body, err = json.Marshal(tt.body)
				if err != nil {
					t.Fatalf("failed to marshal body: %v", err)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/ingest/messages", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.AddMessages(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, res.StatusCode)
			}

			var response dto.ErrorResponse
			if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if response.Error != tt.expectedError {
				t.Errorf("expected error %s, got %s", tt.expectedError, response.Error)
			}
		})
	}
}

func TestAddEntityNodeValidation(t *testing.T) {
	handler := NewIngestHandler(nil)

	tests := []struct {
		name           string
		body           interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "invalid JSON",
			body:           "not json",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid_request",
		},
		{
			name: "missing group_id",
			body: dto.AddEntityNodeRequest{
				GroupID: "",
				Name:    "test-entity",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid_request",
		},
		{
			name: "missing name",
			body: dto.AddEntityNodeRequest{
				GroupID: "test-group",
				Name:    "",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid_request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			var err error
			if s, ok := tt.body.(string); ok {
				body = []byte(s)
			} else {
				body, err = json.Marshal(tt.body)
				if err != nil {
					t.Fatalf("failed to marshal body: %v", err)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/ingest/entity", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.AddEntityNode(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, res.StatusCode)
			}

			var response dto.ErrorResponse
			if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if response.Error != tt.expectedError {
				t.Errorf("expected error %s, got %s", tt.expectedError, response.Error)
			}
		})
	}
}

func TestClearDataValidation(t *testing.T) {
	handler := NewIngestHandler(nil)

	tests := []struct {
		name           string
		body           interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "invalid JSON",
			body:           "not json",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid_request",
		},
		{
			name: "empty group_ids",
			body: dto.ClearDataRequest{
				GroupIDs: []string{},
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid_request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			var err error
			if s, ok := tt.body.(string); ok {
				body = []byte(s)
			} else {
				body, err = json.Marshal(tt.body)
				if err != nil {
					t.Fatalf("failed to marshal body: %v", err)
				}
			}

			req := httptest.NewRequest(http.MethodDelete, "/ingest/clear", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.ClearData(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, res.StatusCode)
			}

			var response dto.ErrorResponse
			if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if response.Error != tt.expectedError {
				t.Errorf("expected error %s, got %s", tt.expectedError, response.Error)
			}
		})
	}
}
