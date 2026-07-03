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
	"context"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var (
	grpcRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "productcatalogservice_grpc_requests_total",
		Help: "Total number of gRPC requests processed, by method and status code.",
	}, []string{"method", "code"})

	grpcRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "productcatalogservice_grpc_request_duration_seconds",
		Help:    "Latency of gRPC requests in seconds, by method.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method"})

	productNotFoundTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "productcatalogservice_product_not_found_total",
		Help: "Total number of GetProduct lookups for an unknown product ID.",
	})

	searchResultsCount = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "productcatalogservice_search_results",
		Help:    "Number of results returned per SearchProducts call.",
		Buckets: []float64{0, 1, 2, 3, 5, 8, 13, 21},
	})
)

// metricsUnaryInterceptor records request counts and latency for every
// unary gRPC method on this service.
func metricsUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	grpcRequestDuration.WithLabelValues(info.FullMethod).Observe(time.Since(start).Seconds())
	grpcRequestsTotal.WithLabelValues(info.FullMethod, status.Code(err).String()).Inc()
	return resp, err
}

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
