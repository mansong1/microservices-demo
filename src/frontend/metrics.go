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
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP request metrics
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "frontend_http_requests_total",
			Help: "Total number of HTTP requests received",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "frontend_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// OpenTelemetry-style HTTP server request duration metric
	httpServerRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_server_request_duration_seconds",
			Help:    "Duration of HTTP server requests in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"method", "route", "status_code"},
	)

	// Handler response time histogram with web-optimized buckets
	handlerResponseTime = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "frontend_handler_response_time_seconds",
			Help:    "Response time for handlers in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"handler", "method", "status"},
	)

	// Business logic metrics
	cartOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "frontend_cart_operations_total",
			Help: "Total number of cart operations",
		},
		[]string{"operation", "status"},
	)

	productViewsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "frontend_product_views_total",
			Help: "Total number of product page views",
		},
	)

	ordersTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "frontend_orders_total",
			Help: "Total number of orders placed",
		},
		[]string{"status"},
	)

	orderValue = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "frontend_order_value_usd",
			Help:    "Order value in USD",
			Buckets: []float64{1, 5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000},
		},
	)

	// External service call metrics
	grpcRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "frontend_grpc_requests_total",
			Help: "Total number of gRPC requests to backend services",
		},
		[]string{"service", "method", "status"},
	)

	grpcRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "frontend_grpc_request_duration_seconds",
			Help:    "gRPC request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service", "method"},
	)

	// Session metrics
	activeSessionsTotal = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "frontend_active_sessions_total",
			Help: "Total number of active user sessions",
		},
	)

	// Currency conversion metrics
	currencyConversionsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "frontend_currency_conversions_total",
			Help: "Total number of currency conversions",
		},
		[]string{"from_currency", "to_currency"},
	)

	// Recommendation metrics
	recommendationsServedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "frontend_recommendations_served_total",
			Help: "Total number of product recommendations served",
		},
	)

	// Error metrics
	errorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "frontend_errors_total",
			Help: "Total number of errors by type",
		},
		[]string{"error_type", "handler"},
	)
)

// Helper functions for recording metrics
func recordHTTPRequest(method, path, status string, duration time.Duration) {
	httpRequestsTotal.WithLabelValues(method, path, status).Inc()
	httpRequestDuration.WithLabelValues(method, path).Observe(duration.Seconds())

	// Also record using OpenTelemetry-style metric name
	// Use path as route and status as status_code to match OTel conventions
	httpServerRequestDuration.WithLabelValues(method, path, status).Observe(duration.Seconds())
}

func recordHandlerResponseTime(handler, method, status string, duration time.Duration) {
	handlerResponseTime.WithLabelValues(handler, method, status).Observe(duration.Seconds())
}

func recordCartOperation(operation, status string) {
	cartOperationsTotal.WithLabelValues(operation, status).Inc()
}

func recordGRPCRequest(service, method, status string, duration time.Duration) {
	grpcRequestsTotal.WithLabelValues(service, method, status).Inc()
	grpcRequestDuration.WithLabelValues(service, method).Observe(duration.Seconds())
}

func recordError(errorType, handler string) {
	errorsTotal.WithLabelValues(errorType, handler).Inc()
}
