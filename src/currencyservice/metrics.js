/*
 * Copyright 2026 Google LLC.
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

const http = require('http');
const client = require('prom-client');

const register = new client.Registry();
client.collectDefaultMetrics({ register });

const rpcRequestsTotal = new client.Counter({
  name: 'currencyservice_rpc_requests_total',
  help: 'Total number of gRPC requests processed, by method and status.',
  labelNames: ['method', 'status'],
  registers: [register]
});

const rpcDuration = new client.Histogram({
  name: 'currencyservice_rpc_duration_seconds',
  help: 'Latency of gRPC requests in seconds, by method.',
  labelNames: ['method'],
  buckets: [0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5],
  registers: [register]
});

const conversionsTotal = new client.Counter({
  name: 'currencyservice_conversions_total',
  help: 'Total number of currency conversions, by source and target currency.',
  labelNames: ['from_currency', 'to_currency'],
  registers: [register]
});

/**
 * Wraps a gRPC handler function so that every call is timed and counted,
 * regardless of whether it succeeds or fails.
 */
function instrument (method, fn) {
  return function (call, callback) {
    const end = rpcDuration.startTimer({ method });
    fn(call, function (err, result) {
      end();
      rpcRequestsTotal.inc({ method, status: err ? 'error' : 'success' });
      callback(err, result);
    });
  };
}

function startMetricsServer (logger) {
  const port = process.env.METRICS_PORT || 9090;
  const server = http.createServer(async (req, res) => {
    if (req.url === '/metrics') {
      res.setHeader('Content-Type', register.contentType);
      res.end(await register.metrics());
      return;
    }
    res.writeHead(404);
    res.end();
  });
  server.listen(port, () => {
    logger.info(`serving Prometheus metrics on :${port}/metrics`);
  });
}

module.exports = { register, rpcRequestsTotal, rpcDuration, conversionsTotal, instrument, startMetricsServer };
