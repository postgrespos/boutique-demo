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

using Prometheus;

namespace cartservice.services
{
    /// <summary>
    /// Custom Prometheus metrics for CartService, on top of the automatic
    /// HTTP/gRPC request metrics registered by prometheus-net.AspNetCore.
    /// </summary>
    public static class CartServiceMetrics
    {
        public static readonly Counter ItemsAddedTotal = Prometheus.Metrics.CreateCounter(
            "cartservice_items_added_total",
            "Total number of items added to carts.");

        public static readonly Counter CartsEmptiedTotal = Prometheus.Metrics.CreateCounter(
            "cartservice_carts_emptied_total",
            "Total number of times a cart was emptied.");

        public static readonly Histogram CartSizeOnRead = Prometheus.Metrics.CreateHistogram(
            "cartservice_cart_size_on_read",
            "Number of distinct line items in a cart at the time it was read.",
            new HistogramConfiguration
            {
                Buckets = Histogram.LinearBuckets(start: 0, width: 1, count: 10)
            });

        public static readonly Counter CartStoreErrorsTotal = Prometheus.Metrics.CreateCounter(
            "cartservice_cart_store_errors_total",
            "Total number of cart store operations that failed, by operation.",
            new CounterConfiguration
            {
                LabelNames = new[] { "operation" }
            });
    }
}
