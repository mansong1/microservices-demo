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
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestHTTPRequestMetrics(t *testing.T) {
	// Just test that the function doesn't panic
	recordHTTPRequest("GET", "/", "200", 100*time.Millisecond)

	// This test mainly ensures the function doesn't panic
	// In a real scenario, you'd use more sophisticated metrics testing
}

func TestHandlerResponseTimeMetrics(t *testing.T) {
	// Record handler response time
	recordHandlerResponseTime("home", "GET", "200", 50*time.Millisecond)

	// For histogram metrics, we can check that it doesn't panic and has been recorded
	// by using a simple metric collection rather than trying to get the exact count
	gatherer := prometheus.DefaultGatherer
	metricFamilies, err := gatherer.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	// Look for our handler response time metric
	found := false
	for _, mf := range metricFamilies {
		if mf.GetName() == "frontend_handler_response_time_seconds" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected handler response time metric to be registered")
	}
}

func TestCartOperationMetrics(t *testing.T) {
	// Record a cart operation
	recordCartOperation("add", "success")

	// Check if metric was recorded
	count := testutil.ToFloat64(cartOperationsTotal.WithLabelValues("add", "success"))
	if count == 0 {
		t.Error("Expected cart operation metric to be recorded")
	}
}

func TestGRPCRequestMetrics(t *testing.T) {
	// Record a gRPC request
	recordGRPCRequest("productcatalog", "GetProduct", "success", 25*time.Millisecond)

	// Check if counter metric was recorded
	count := testutil.ToFloat64(grpcRequestsTotal.WithLabelValues("productcatalog", "GetProduct", "success"))
	if count == 0 {
		t.Error("Expected gRPC request metric to be recorded")
	}

	// For histogram metrics, check that the metric family is registered
	gatherer := prometheus.DefaultGatherer
	metricFamilies, err := gatherer.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	// Look for our gRPC request duration metric
	found := false
	for _, mf := range metricFamilies {
		if mf.GetName() == "frontend_grpc_request_duration_seconds" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected gRPC request duration metric to be registered")
	}
}

func TestErrorMetrics(t *testing.T) {
	// Record an error
	recordError("grpc_error", "homeHandler")

	// Check if metric was recorded
	count := testutil.ToFloat64(errorsTotal.WithLabelValues("grpc_error", "homeHandler"))
	if count == 0 {
		t.Error("Expected error metric to be recorded")
	}
}

func TestProductViewMetrics(t *testing.T) {
	// Get initial count
	initialCount := testutil.ToFloat64(productViewsTotal)

	// Increment product views
	productViewsTotal.Inc()

	// Check if metric was incremented
	newCount := testutil.ToFloat64(productViewsTotal)
	if newCount != initialCount+1 {
		t.Errorf("Expected product views to increase by 1, got %f -> %f", initialCount, newCount)
	}
}

func TestOrderMetrics(t *testing.T) {
	// Get initial count
	initialCount := testutil.ToFloat64(ordersTotal.WithLabelValues("success"))

	// Record an order
	ordersTotal.WithLabelValues("success").Inc()
	orderValue.Observe(99.99)

	// Check if counter was incremented
	newCount := testutil.ToFloat64(ordersTotal.WithLabelValues("success"))
	if newCount != initialCount+1 {
		t.Errorf("Expected orders total to increase by 1, got %f -> %f", initialCount, newCount)
	}

	// For histogram metrics, check that the metric family is registered
	gatherer := prometheus.DefaultGatherer
	metricFamilies, err := gatherer.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	// Look for our order value metric
	found := false
	for _, mf := range metricFamilies {
		if mf.GetName() == "frontend_order_value_usd" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected order value histogram metric to be registered")
	}
}

func TestCurrencyConversionMetrics(t *testing.T) {
	// Get initial count
	initialCount := testutil.ToFloat64(currencyConversionsTotal.WithLabelValues("USD", "EUR"))

	// Record a currency conversion
	currencyConversionsTotal.WithLabelValues("USD", "EUR").Inc()

	// Check if metric was incremented
	newCount := testutil.ToFloat64(currencyConversionsTotal.WithLabelValues("USD", "EUR"))
	if newCount != initialCount+1 {
		t.Errorf("Expected currency conversions to increase by 1, got %f -> %f", initialCount, newCount)
	}
}

func TestActiveSessionsTracking(t *testing.T) {
	// Test session tracking
	originalValue := testutil.ToFloat64(activeSessionsTotal)

	// Simulate session creation
	activeSessionsTotal.Inc()

	newValue := testutil.ToFloat64(activeSessionsTotal)
	if newValue != originalValue+1 {
		t.Errorf("Expected active sessions to increase by 1, got %f -> %f", originalValue, newValue)
	}

	// Simulate session destruction
	activeSessionsTotal.Dec()

	finalValue := testutil.ToFloat64(activeSessionsTotal)
	if finalValue != originalValue {
		t.Errorf("Expected active sessions to return to original value %f, got %f", originalValue, finalValue)
	}
}

func TestMetricsIntegration(t *testing.T) {
	// Test that all metric recording functions don't panic
	tests := []struct {
		name string
		fn   func()
	}{
		{"recordHTTPRequest", func() { recordHTTPRequest("POST", "/cart", "200", 10*time.Millisecond) }},
		{"recordHandlerResponseTime", func() { recordHandlerResponseTime("cart", "POST", "200", 5*time.Millisecond) }},
		{"recordCartOperation", func() { recordCartOperation("add", "success") }},
		{"recordGRPCRequest", func() { recordGRPCRequest("cart", "AddItem", "success", 15*time.Millisecond) }},
		{"recordError", func() { recordError("validation_error", "cartHandler") }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic
			tt.fn()
		})
	}
}

func TestRecommendationMetrics(t *testing.T) {
	// Get initial count
	initialCount := testutil.ToFloat64(recommendationsServedTotal)

	// Record recommendations served
	recommendationsServedTotal.Inc()

	// Check if metric was incremented
	newCount := testutil.ToFloat64(recommendationsServedTotal)
	if newCount != initialCount+1 {
		t.Errorf("Expected recommendations served to increase by 1, got %f -> %f", initialCount, newCount)
	}
}

func TestConcurrentMetrics(t *testing.T) {
	// Test concurrent access to metrics
	done := make(chan bool)

	// Start multiple goroutines that record metrics
	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()

			// Record various metrics
			productViewsTotal.Inc()
			ordersTotal.WithLabelValues("success").Inc()
			activeSessionsTotal.Inc()
			currencyConversionsTotal.WithLabelValues("USD", "EUR").Inc()
			recommendationsServedTotal.Inc()
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify that metrics were recorded (exact values depend on previous tests)
	if testutil.ToFloat64(productViewsTotal) == 0 {
		t.Error("Expected product views to be recorded")
	}
}
