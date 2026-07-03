/*
 * Copyright 2026, Google LLC.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package hipstershop;

import io.prometheus.metrics.core.metrics.Counter;
import io.prometheus.metrics.core.metrics.Histogram;
import io.prometheus.metrics.exporter.httpserver.HTTPServer;
import io.prometheus.metrics.instrumentation.jvm.JvmMetrics;
import java.io.IOException;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;

/** Registers and exposes Prometheus metrics for AdService on a dedicated HTTP port. */
final class AdServiceMetrics {

  private static final Logger logger = LogManager.getLogger(AdServiceMetrics.class);

  static final Counter GRPC_REQUESTS_TOTAL =
      Counter.builder()
          .name("adservice_grpc_requests_total")
          .help("Total number of gRPC requests processed, by method and status.")
          .labelNames("method", "status")
          .register();

  static final Histogram GRPC_REQUEST_DURATION_SECONDS =
      Histogram.builder()
          .name("adservice_grpc_request_duration_seconds")
          .help("Latency of gRPC requests in seconds, by method.")
          .labelNames("method")
          .register();

  static final Histogram ADS_SERVED =
      Histogram.builder()
          .name("adservice_ads_served")
          .help("Number of ads served per GetAds call.")
          .classicUpperBounds(0, 1, 2, 3, 5)
          .register();

  static final Counter ADS_FALLBACK_TOTAL =
      Counter.builder()
          .name("adservice_ads_fallback_total")
          .help("Total number of GetAds calls that fell back to random ads.")
          .register();

  private static HTTPServer server;

  private AdServiceMetrics() {}

  static void start() {
    int port = Integer.parseInt(System.getenv().getOrDefault("METRICS_PORT", "9090"));
    try {
      JvmMetrics.builder().register();
      server = HTTPServer.builder().port(port).buildAndStart();
      logger.info("serving Prometheus metrics on :" + port + "/metrics");
    } catch (IOException e) {
      logger.warn("failed to start Prometheus metrics server: " + e.getMessage());
    }
  }
}
