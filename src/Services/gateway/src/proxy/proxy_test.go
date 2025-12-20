package proxy

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestProxyHandler_Routing tests proxy routing logic
func TestProxyHandler_Routing(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedTarget string
		expectedCode   int
	}{
		{
			name:           "route_to_contest_service",
			path:           "/api/v1/contest/123",
			expectedTarget: "contest",
			expectedCode:   http.StatusOK,
		},
		{
			name:           "route_to_submission_service",
			path:           "/api/v1/submission/submit",
			expectedTarget: "submission-judge",
			expectedCode:   http.StatusOK,
		},
		{
			name:           "route_to_problem_service",
			path:           "/api/v1/problem/list",
			expectedTarget: "problem",
			expectedCode:   http.StatusOK,
		},
		{
			name:           "invalid_route",
			path:           "/api/v1/invalid",
			expectedTarget: "",
			expectedCode:   http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test request
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			// Create proxy handler
			handler := NewProxyHandler()
			handler.ServeHTTP(w, req)

			// Verify response code
			if w.Code != tt.expectedCode {
				t.Errorf("Expected status %d, got %d", tt.expectedCode, w.Code)
			}
		})
	}
}

// TestProxyHandler_LoadBalancing tests load balancing across multiple instances
func TestProxyHandler_LoadBalancing(t *testing.T) {
	// Setup multiple backend servers
	backend1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("backend1"))
	}))
	defer backend1.Close()

	backend2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("backend2"))
	}))
	defer backend2.Close()

	// Configure proxy with multiple backends
	proxy := NewLoadBalancingProxy([]string{backend1.URL, backend2.URL})

	// Make multiple requests
	responses := make(map[string]int)
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/contest", nil)
		w := httptest.NewRecorder()
		proxy.ServeHTTP(w, req)
		
		body := w.Body.String()
		responses[body]++
	}

	// Verify load distribution (should hit both backends)
	if len(responses) != 2 {
		t.Error("Load balancing failed, not all backends were hit")
	}
}

// TestProxyHandler_HeaderForwarding tests header propagation
func TestProxyHandler_HeaderForwarding(t *testing.T) {
	// Setup backend server that echoes headers
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Echo Authorization header
		if auth := r.Header.Get("Authorization"); auth != "" {
			w.Header().Set("X-Forwarded-Auth", auth)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	proxy := NewProxyWithBackend(backend.URL)

	// Create request with auth header
	req := httptest.NewRequest(http.MethodGet, "/api/v1/contest", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()

	proxy.ServeHTTP(w, req)

	// Verify header was forwarded
	if w.Header().Get("X-Forwarded-Auth") != "Bearer test-token" {
		t.Error("Authorization header was not properly forwarded")
	}
}

// TestProxyHandler_CircuitBreaker tests circuit breaker functionality
func TestProxyHandler_CircuitBreaker(t *testing.T) {
	failCount := 0
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		failCount++
		if failCount <= 5 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	proxy := NewProxyWithCircuitBreaker(backend.URL, 3) // Open after 3 failures

	// Make requests until circuit opens
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/test", nil)
		w := httptest.NewRecorder()
		proxy.ServeHTTP(w, req)
	}

	// Next request should be blocked by circuit breaker
	req := httptest.NewRequest(http.MethodGet, "/api/v1/test", nil)
	w := httptest.NewRecorder()
	proxy.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Error("Circuit breaker should have blocked the request")
	}
}

// TestProxyHandler_Timeout tests request timeout handling
func TestProxyHandler_Timeout(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // Simulate slow backend
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	proxy := NewProxyWithTimeout(backend.URL, 500*time.Millisecond)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/test", nil)
	w := httptest.NewRecorder()
	proxy.ServeHTTP(w, req)

	if w.Code != http.StatusGatewayTimeout {
		t.Errorf("Expected timeout status, got %d", w.Code)
	}
}

// TestProxyHandler_Retry tests retry logic
func TestProxyHandler_Retry(t *testing.T) {
	attemptCount := 0
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		if attemptCount < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	proxy := NewProxyWithRetry(backend.URL, 3) // Retry up to 3 times

	req := httptest.NewRequest(http.MethodGet, "/api/v1/test", nil)
	w := httptest.NewRecorder()
	proxy.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("Retry logic failed")
	}

	if attemptCount != 3 {
		t.Errorf("Expected 3 attempts, got %d", attemptCount)
	}
}
