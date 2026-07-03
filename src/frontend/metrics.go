// Copyright 2026 Google LLC
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
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "frontend_http_requests_total",
		Help: "Total number of HTTP requests handled by the frontend, by route, method and status code.",
	}, []string{"route", "method", "status"})

	httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "frontend_http_request_duration_seconds",
		Help:    "Latency of HTTP requests handled by the frontend, in seconds.",
		Buckets: prometheus.DefBuckets,
	}, []string{"route", "method"})

	cartItemsAdded = promauto.NewCounter(prometheus.CounterOpts{
		Name: "frontend_cart_items_added_total",
		Help: "Total number of items added to carts through the frontend.",
	})

	ordersPlaced = promauto.NewCounter(prometheus.CounterOpts{
		Name: "frontend_orders_placed_total",
		Help: "Total number of orders successfully placed through the frontend.",
	})

	orderValueUSD = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "frontend_order_value_usd",
		Help:    "Total value (in USD major units) of orders placed through the frontend.",
		Buckets: []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	})
)

// metricsMiddleware records request counts and latency per route template so
// that high-cardinality path parameters (like product IDs) don't blow up
// metric cardinality.
func metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)

		route := r.URL.Path
		if m := mux.CurrentRoute(r); m != nil {
			if tmpl, err := m.GetPathTemplate(); err == nil {
				route = tmpl
			}
		}
		httpRequestsTotal.WithLabelValues(route, r.Method, strconv.Itoa(rec.status)).Inc()
		httpRequestDuration.WithLabelValues(route, r.Method).Observe(time.Since(start).Seconds())
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}
