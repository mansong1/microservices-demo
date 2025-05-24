// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func TestLogHandler(t *testing.T) {
	// Create a test handler that the log handler will wrap
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify that context values are set
		if r.Context().Value(ctxKeyLog{}) == nil {
			t.Error("log context not set")
		}
		if r.Context().Value(ctxKeyRequestID{}) == nil {
			t.Error("request ID context not set")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	log := logrus.New()
	lh := &logHandler{
		log:  log,
		next: testHandler,
	}

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	lh.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("logHandler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}

	body := rr.Body.String()
	if body != "test response" {
		t.Errorf("logHandler returned wrong body: got %v want %v", body, "test response")
	}
}

func TestMetricsHandler(t *testing.T) {
	// Create a test handler that the metrics handler will wrap
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	mh := &metricsHandler{
		next: testHandler,
	}

	tests := []struct {
		name   string
		path   string
		method string
	}{
		{"home page", "/", "GET"},
		{"product page", "/product/123", "GET"},
		{"cart page", "/cart", "GET"},
		{"health check", "/_healthz", "GET"},
		{"static file", "/static/style.css", "GET"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rr := httptest.NewRecorder()

			mh.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("metricsHandler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
			}
		})
	}
}

func TestResponseRecorder(t *testing.T) {
	// Test responseRecorder functionality
	w := httptest.NewRecorder()
	rr := &responseRecorder{w: w}

	// Test Header method
	rr.Header().Set("Test-Header", "test-value")
	if rr.Header().Get("Test-Header") != "test-value" {
		t.Error("responseRecorder Header() not working correctly")
	}

	// Test Write method
	data := []byte("test data")
	n, err := rr.Write(data)
	if err != nil {
		t.Errorf("responseRecorder Write() error: %v", err)
	}
	if n != len(data) {
		t.Errorf("responseRecorder Write() returned wrong count: got %v want %v", n, len(data))
	}
	if rr.b != len(data) {
		t.Errorf("responseRecorder bytes count wrong: got %v want %v", rr.b, len(data))
	}
	if rr.status != http.StatusOK {
		t.Errorf("responseRecorder status should be 200 after Write(), got %v", rr.status)
	}

	// Test WriteHeader method
	rr2 := &responseRecorder{w: httptest.NewRecorder()}
	rr2.WriteHeader(http.StatusNotFound)
	if rr2.status != http.StatusNotFound {
		t.Errorf("responseRecorder WriteHeader() not working: got %v want %v", rr2.status, http.StatusNotFound)
	}
}

func TestEnsureSessionID(t *testing.T) {
	// Test handler that checks for session ID
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie(cookieSessionID)
		if err != nil {
			t.Error("session cookie not found")
			return
		}
		if sessionCookie.Value == "" {
			t.Error("session cookie has empty value")
			return
		}
		
		// Verify session ID is in context
		if r.Context().Value(ctxKeySessionID{}) == nil {
			t.Error("session ID not in context")
		}
		
		w.WriteHeader(http.StatusOK)
	})

	handler := ensureSessionID(testHandler)

	// Test with no existing session
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("ensureSessionID handler returned wrong status: got %v want %v", rr.Code, http.StatusOK)
	}

	// Check that Set-Cookie header was added
	cookies := rr.Header()["Set-Cookie"]
	if len(cookies) == 0 {
		t.Error("ensureSessionID should set a session cookie")
	}

	// Test with existing session
	req2 := httptest.NewRequest("GET", "/", nil)
	req2.AddCookie(&http.Cookie{
		Name:  cookieSessionID,
		Value: "existing-session-123",
	})
	rr2 := httptest.NewRecorder()

	testHandler2 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, _ := r.Cookie(cookieSessionID)
		if sessionCookie.Value != "existing-session-123" {
			t.Error("existing session ID should be preserved")
		}
		w.WriteHeader(http.StatusOK)
	})

	handler2 := ensureSessionID(testHandler2)
	handler2.ServeHTTP(rr2, req2)

	if rr2.Code != http.StatusOK {
		t.Errorf("ensureSessionID with existing session returned wrong status: got %v want %v", rr2.Code, http.StatusOK)
	}
}

func TestMiddlewareChain(t *testing.T) {
	// Test the complete middleware chain
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify all context values are set
		if r.Context().Value(ctxKeyLog{}) == nil {
			t.Error("log context not set in middleware chain")
		}
		if r.Context().Value(ctxKeyRequestID{}) == nil {
			t.Error("request ID context not set in middleware chain")
		}
		if r.Context().Value(ctxKeySessionID{}) == nil {
			t.Error("session ID context not set in middleware chain")
		}
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	// Build the middleware chain like in main.go
	var handler http.Handler = finalHandler
	log := logrus.New()
	handler = &logHandler{log: log, next: handler}
	handler = ensureSessionID(handler)
	handler = &metricsHandler{next: handler}

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("middleware chain returned wrong status: got %v want %v", rr.Code, http.StatusOK)
	}

	if rr.Body.String() != "success" {
		t.Errorf("middleware chain returned wrong body: got %v want %v", rr.Body.String(), "success")
	}
}

func TestContextKeys(t *testing.T) {
	// Test that context keys work correctly
	ctx := context.Background()
	
	// Test ctxKeyLog
	log := logrus.New()
	ctx = context.WithValue(ctx, ctxKeyLog{}, log)
	retrievedLog := ctx.Value(ctxKeyLog{})
	if retrievedLog != log {
		t.Error("ctxKeyLog not working correctly")
	}

	// Test ctxKeyRequestID
	requestID := "test-request-123"
	ctx = context.WithValue(ctx, ctxKeyRequestID{}, requestID)
	retrievedRequestID := ctx.Value(ctxKeyRequestID{})
	if retrievedRequestID != requestID {
		t.Error("ctxKeyRequestID not working correctly")
	}

	// Test ctxKeySessionID
	sessionID := "test-session-123"
	ctx = context.WithValue(ctx, ctxKeySessionID{}, sessionID)
	retrievedSessionID := ctx.Value(ctxKeySessionID{})
	if retrievedSessionID != sessionID {
		t.Error("ctxKeySessionID not working correctly")
	}
}
