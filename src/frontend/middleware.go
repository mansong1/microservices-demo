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
	"strconv"
	"strings"
	"time"
	"os"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ctxKeyLog struct{}
type ctxKeyRequestID struct{}

type logHandler struct {
	log  *logrus.Logger
	next http.Handler
}

type metricsHandler struct {
	next http.Handler
}

type responseRecorder struct {
	b      int
	status int
	w      http.ResponseWriter
}

func (r *responseRecorder) Header() http.Header { return r.w.Header() }

func (r *responseRecorder) Write(p []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	n, err := r.w.Write(p)
	r.b += n
	return n, err
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
	r.w.WriteHeader(statusCode)
}

func (lh *logHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID, _ := uuid.NewRandom()
	ctx = context.WithValue(ctx, ctxKeyRequestID{}, requestID.String())

	start := time.Now()
	rr := &responseRecorder{w: w}
	log := lh.log.WithFields(logrus.Fields{
		"http.req.path":   r.URL.Path,
		"http.req.method": r.Method,
		"http.req.id":     requestID.String(),
	})
	if v, ok := r.Context().Value(ctxKeySessionID{}).(string); ok {
		log = log.WithField("session", v)
	}
	log.Debug("request started")
	defer func() {
		log.WithFields(logrus.Fields{
			"http.resp.took_ms": int64(time.Since(start) / time.Millisecond),
			"http.resp.status":  rr.status,
			"http.resp.bytes":   rr.b}).Debugf("request complete")
	}()

	ctx = context.WithValue(ctx, ctxKeyLog{}, log)
	r = r.WithContext(ctx)
	lh.next.ServeHTTP(rr, r)
}

func (mh *metricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	rr := &responseRecorder{w: w}
	
	// Extract path without query parameters and remove any path variables
	path := r.URL.Path
	handlerName := path
	
	// Normalize paths with IDs to avoid high cardinality
	if strings.Contains(path, "/product/") {
		path = "/product/{id}"
		handlerName = "product"
	} else if strings.Contains(path, "/product-meta/") {
		path = "/product-meta/{ids}"
		handlerName = "product-meta"
	} else if path == "/" {
		handlerName = "home"
	} else if path == "/cart" {
		handlerName = "cart"
	} else if path == "/cart/checkout" {
		handlerName = "checkout"
	} else if path == "/set_currency" {
		handlerName = "set-currency"
	} else if strings.HasPrefix(path, "/static/") {
		handlerName = "static"
	} else if path == "/_healthz" {
		handlerName = "health"
	} else {
		// Extract handler name from path for other routes
		parts := strings.Split(strings.Trim(path, "/"), "/")
		if len(parts) > 0 && parts[0] != "" {
			handlerName = parts[0]
		} else {
			handlerName = "unknown"
		}
	}
	
	// Call the next handler
	mh.next.ServeHTTP(rr, r)
	
	// Record metrics
	duration := time.Since(start)
	statusCode := strconv.Itoa(rr.status)
	recordHTTPRequest(r.Method, path, statusCode, duration)
	recordHandlerResponseTime(handlerName, r.Method, statusCode, duration)
}

func ensureSessionID(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sessionID string
		var newSession bool
		c, err := r.Cookie(cookieSessionID)
		if err == http.ErrNoCookie {
			newSession = true
			if os.Getenv("ENABLE_SINGLE_SHARED_SESSION") == "true" {
				// Hard coded user id, shared across sessions
				sessionID = "12345678-1234-1234-1234-123456789123"
			} else {
				u, _ := uuid.NewRandom()
				sessionID = u.String()
			}
			http.SetCookie(w, &http.Cookie{
				Name:   cookieSessionID,
				Value:  sessionID,
				MaxAge: cookieMaxAge,
			})
		} else if err != nil {
			return
		} else {
			sessionID = c.Value
		}
		
		// Record session metrics for new sessions
		if newSession {
			activeSessionsTotal.Inc()
		}
		
		ctx := context.WithValue(r.Context(), ctxKeySessionID{}, sessionID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}
