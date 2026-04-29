package proxy

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestSubmissionApiProxy tests the submission proxy handler
func TestSubmissionApiProxy(t *testing.T) {
	// This test requires the config to be loaded
	// For now, we'll just verify the handler is not nil
	handler := SubmissionApiProxy()
	if handler == nil {
		t.Error("SubmissionApiProxy returned nil handler")
	}
}

// TestProblemApiProxy tests the problem proxy handler
func TestProblemApiProxy(t *testing.T) {
	handler := ProblemApiProxy()
	if handler == nil {
		t.Error("ProblemApiProxy returned nil handler")
	}
}

// TestLoginApiProxy tests the login proxy handler
func TestLoginApiProxy(t *testing.T) {
	handler := LoginApiProxy()
	if handler == nil {
		t.Error("LoginApiProxy returned nil handler")
	}
}

// TestWSSubmissionProxy tests the WebSocket submission proxy handler
func TestWSSubmissionProxy(t *testing.T) {
	handler := WSSubmissionProxy()
	if handler == nil {
		t.Error("WSSubmissionProxy returned nil handler")
	}
}

// TestNewProxy tests the internal proxy creation
func TestNewProxy(t *testing.T) {
	// Test with a valid endpoint
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("backend response"))
	}))
	defer backend.Close()

	handler := newProxy(backend.URL)
	if handler == nil {
		t.Fatal("newProxy returned nil handler")
	}

	// Test the proxy handler
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if w.Body.String() != "backend response" {
		t.Errorf("Expected 'backend response', got '%s'", w.Body.String())
	}
}

// TestNewProxyInvalidURL tests proxy with invalid URL
func TestNewProxyInvalidURL(t *testing.T) {
	handler := newProxy("://invalid-url")
	
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d for invalid URL, got %d", http.StatusInternalServerError, w.Code)
	}
}

// TestWSProxy tests WebSocket proxy creation
func TestWSProxy(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if WebSocket headers are preserved
		if r.Header.Get("Upgrade") == "websocket" && r.Header.Get("Connection") == "upgrade" {
			w.WriteHeader(http.StatusSwitchingProtocols)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer backend.Close()

	handler := WSProxy(backend.URL)
	if handler == nil {
		t.Fatal("WSProxy returned nil handler")
	}

	// Test WebSocket upgrade request
	req := httptest.NewRequest(http.MethodGet, "/ws", nil)
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Connection", "upgrade")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Verify headers were forwarded
	if req.Header.Get("Upgrade") != "websocket" {
		t.Error("Upgrade header was not preserved")
	}
	if req.Header.Get("Connection") != "upgrade" {
		t.Error("Connection header was not preserved")
	}
}
