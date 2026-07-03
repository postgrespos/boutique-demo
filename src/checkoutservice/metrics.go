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
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	ordersTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "checkoutservice_orders_total",
		Help: "Total number of PlaceOrder requests processed, by outcome.",
	}, []string{"status"})

	orderValueUSD = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "checkoutservice_order_value_usd",
		Help:    "Total value (in USD major units) of placed orders.",
		Buckets: []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	})

	orderItemCount = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "checkoutservice_order_items",
		Help:    "Number of line items per placed order.",
		Buckets: []float64{1, 2, 3, 5, 8, 13, 21},
	})

	placeOrderDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "checkoutservice_place_order_duration_seconds",
		Help:    "Latency of PlaceOrder requests in seconds.",
		Buckets: prometheus.DefBuckets,
	})

	dependencyRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "checkoutservice_dependency_requests_total",
		Help: "Total number of outbound requests to downstream services, by dependency and outcome.",
	}, []string{"dependency", "status"})

	dependencyDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "checkoutservice_dependency_duration_seconds",
		Help:    "Latency of outbound requests to downstream services in seconds.",
		Buckets: prometheus.DefBuckets,
	}, []string{"dependency"})
)

// startMetricsServer exposes Prometheus metrics on /metrics on a dedicated
// port so it can be scraped independently of the gRPC server.
func startMetricsServer() {
	port := os.Getenv("METRICS_PORT")
	if port == "" {
		port = "9090"
	}
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Infof("serving Prometheus metrics on :%s/metrics", port)
		if err := http.ListenAndServe(":"+port, mux); err != nil {
			log.Warnf("metrics server stopped: %+v", err)
		}
	}()
}
