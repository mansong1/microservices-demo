# Prometheus Metrics Implementation - Frontend Service

This document describes the comprehensive Prometheus metrics implementation added to the frontend service of the microservices-demo project.

## Overview

The frontend service now exposes detailed metrics for monitoring and observability through a Prometheus endpoint at `/metrics`. The implementation includes HTTP request metrics, business metrics, gRPC metrics, and handler response time histograms.

## Metrics Endpoint

- **URL**: `http://localhost:8080/metrics`
- **Format**: Prometheus exposition format
- **Update**: Real-time metrics collection
- **Port**: Same port as the main application (8080)

## Implemented Metrics

### 1. HTTP Request Metrics

#### `frontend_http_requests_total` (Counter)
- **Description**: Total number of HTTP requests received
- **Labels**: 
  - `method`: HTTP method (GET, POST, etc.)
  - `path`: Request path (normalized)
  - `status`: HTTP status code
- **Example**: `frontend_http_requests_total{method="GET",path="/",status="200"} 42`

#### `frontend_http_request_duration_seconds` (Histogram)
- **Description**: HTTP request duration in seconds
- **Labels**: `method`, `path`, `status`
- **Buckets**: `0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10`

#### `http_server_request_duration_seconds` (Histogram)
- **Description**: Duration of HTTP server requests in seconds (OpenTelemetry-style naming)
- **Labels**: 
  - `method`: HTTP method (GET, POST, etc.)
  - `route`: Request route/path
  - `status_code`: HTTP status code
- **Buckets**: `0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10`
- **Example**: `http_server_request_duration_seconds_bucket{method="GET",route="/",status_code="200",le="0.1"} 15`

### 2. Handler Response Time Histogram (As Requested)

#### `frontend_handler_response_time_seconds` (Histogram)
- **Description**: Response time for handlers in seconds
- **Labels**: 
  - `handler`: Handler identifier (home, product, cart, checkout, etc.)
  - `method`: HTTP method
  - `status`: HTTP status code
- **Buckets**: Web-optimized (0.001 to 10 seconds)
- **Example**: `frontend_handler_response_time_seconds_bucket{handler="home",method="GET",status="200",le="0.1"} 15`

### 3. Business Metrics

#### `frontend_cart_operations_total` (Counter)
- **Description**: Total number of cart operations performed
- **Labels**: 
  - `operation`: add, empty
  - `status`: success, error

#### `frontend_orders_total` (Counter)
- **Description**: Total number of orders placed
- **Labels**: `status`: success, error

#### `frontend_order_value_usd` (Histogram)
- **Description**: Order value in USD
- **Buckets**: `1, 5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000`

#### `frontend_product_views_total` (Counter)
- **Description**: Total number of product page views

#### `frontend_currency_conversions_total` (Counter)
- **Description**: Total number of currency conversions
- **Labels**: 
  - `from_currency`: Source currency code
  - `to_currency`: Target currency code

#### `frontend_recommendations_served_total` (Counter)
- **Description**: Total number of product recommendations served

### 4. Session Metrics

#### `frontend_active_sessions_total` (Gauge)
- **Description**: Total number of active user sessions

### 5. gRPC Metrics

#### `frontend_grpc_requests_total` (Counter)
- **Description**: Total number of gRPC requests made to backend services
- **Labels**:
  - `service`: Backend service name (ProductCatalogService, CartService, etc.)
  - `method`: gRPC method name
  - `status`: success, error

#### `frontend_grpc_request_duration_seconds` (Histogram)
- **Description**: gRPC request duration in seconds
- **Labels**: `service`, `method`, `status`
- **Buckets**: `0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10`

### 6. Error Metrics

#### `frontend_errors_total` (Counter)
- **Description**: Total number of errors by type
- **Labels**: `error_type`

## Instrumentation Details

### HTTP Middleware
- Automatic metrics collection for all HTTP requests
- Path normalization to prevent cardinality explosion
- Handler identification for response time tracking

### Business Logic Instrumentation
- Cart operations in `addToCartHandler` and `emptyCartHandler`
- Product views in `productHandler`
- Order tracking in `placeOrderHandler`
- Currency conversions in `setCurrencyHandler`
- Session tracking in `ensureSessionID` middleware

### gRPC Client Instrumentation
All backend service calls are instrumented in `rpc.go`:
- ProductCatalogService.GetProduct
- CartService.GetCart, EmptyCart, AddItem
- CurrencyService.Convert
- ShippingService.GetQuote
- RecommendationService.ListRecommendations
- AdService.GetAds

## File Changes

### New Files
- `src/frontend/metrics.go`: Metrics definitions and helper functions

### Modified Files
- `src/frontend/main.go`: Added metrics endpoint and middleware
- `src/frontend/middleware.go`: Added metrics collection middleware
- `src/frontend/handlers.go`: Instrumented business logic handlers
- `src/frontend/rpc.go`: Instrumented gRPC client calls
- `src/frontend/go.mod`: Added Prometheus client dependency

### Kubernetes Configuration
- `kubernetes-manifests/frontend.yaml`: Added metrics port and service

## Testing

### Prerequisites
```bash
cd /Users/mansong/Projects/microservices-demo/src/frontend
go mod tidy
go build -o frontend .
```

### Environment Variables (for testing)
```bash
export PORT=8080
export PRODUCT_CATALOG_SERVICE_ADDR="localhost:3550"
export CURRENCY_SERVICE_ADDR="localhost:7000"
export CART_SERVICE_ADDR="localhost:7070"
export RECOMMENDATION_SERVICE_ADDR="localhost:8080"
export SHIPPING_SERVICE_ADDR="localhost:50051"
export CHECKOUT_SERVICE_ADDR="localhost:5050"
export AD_SERVICE_ADDR="localhost:9555"
export SHOPPING_ASSISTANT_SERVICE_ADDR="localhost:80"
```

### Running the Service
```bash
./frontend
```

### Testing Metrics
```bash
# Check metrics endpoint
curl http://localhost:8080/metrics

# Trigger some requests to generate data
curl http://localhost:8080/

# Check specific metrics
curl -s http://localhost:8080/metrics | grep "frontend_http_requests_total"
curl -s http://localhost:8080/metrics | grep "handler_response_time_seconds"
```

## Example Metrics Output

```
# HELP frontend_handler_response_time_seconds Response time for handlers in seconds
# TYPE frontend_handler_response_time_seconds histogram
frontend_handler_response_time_seconds_bucket{handler="home",method="GET",status="500",le="0.001"} 0
frontend_handler_response_time_seconds_bucket{handler="home",method="GET",status="500",le="0.005"} 0
frontend_handler_response_time_seconds_bucket{handler="home",method="GET",status="500",le="0.01"} 1
frontend_handler_response_time_seconds_count{handler="home",method="GET",status="500"} 1
frontend_handler_response_time_seconds_sum{handler="home",method="GET",status="500"} 0.006543

# HELP frontend_http_requests_total Total number of HTTP requests received
# TYPE frontend_http_requests_total counter
frontend_http_requests_total{method="GET",path="/",status="500"} 1
frontend_http_requests_total{method="GET",path="/metrics",status="200"} 5

# HELP http_server_request_duration_seconds Duration of HTTP server requests in seconds
# TYPE http_server_request_duration_seconds histogram
http_server_request_duration_seconds_bucket{method="GET",route="/",status_code="500",le="0.01"} 1
http_server_request_duration_seconds_count{method="GET",route="/",status_code="500"} 1
http_server_request_duration_seconds_sum{method="GET",route="/",status_code="500"} 0.003763

# HELP frontend_active_sessions_total Total number of active user sessions
# TYPE frontend_active_sessions_total gauge
frontend_active_sessions_total 19
```

## Prometheus Configuration

For Kubernetes deployment, the metrics are exposed via a dedicated service on the same port as the main application:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: frontend-metrics
  annotations:
    prometheus.io/scrape: "true"
    prometheus.io/port: "8080"
    prometheus.io/path: "/metrics"
spec:
  selector:
    app: frontend
  ports:
  - name: metrics
    port: 8080
    targetPort: 8080
```

## Monitoring Queries

### Sample PromQL Queries

```promql
# Request rate by endpoint
rate(frontend_http_requests_total[5m])

# Average response time by handler
rate(frontend_handler_response_time_seconds_sum[5m]) / rate(frontend_handler_response_time_seconds_count[5m])

# Average HTTP server request duration (OpenTelemetry-style)
rate(http_server_request_duration_seconds_sum[5m]) / rate(http_server_request_duration_seconds_count[5m])

# Error rate
rate(frontend_http_requests_total{status=~"5.."}[5m]) / rate(frontend_http_requests_total[5m])

# 95th percentile response time
histogram_quantile(0.95, rate(frontend_handler_response_time_seconds_bucket[5m]))

# 95th percentile HTTP server request duration
histogram_quantile(0.95, rate(http_server_request_duration_seconds_bucket[5m]))

# gRPC success rate
rate(frontend_grpc_requests_total{status="success"}[5m]) / rate(frontend_grpc_requests_total[5m])
```

## Dependencies

- `github.com/prometheus/client_golang v1.20.5`

## Notes

- Metrics collection has minimal performance impact
- Path normalization prevents cardinality explosion
- All metrics include appropriate labels for filtering and aggregation
- Handler response time histogram specifically addresses the requirement for "response time for handlers in seconds"
- Buckets are optimized for web application response times
