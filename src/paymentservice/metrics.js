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

const http = require('http');
const client = require('prom-client');

const register = new client.Registry();
client.collectDefaultMetrics({ register });

const rpcRequestsTotal = new client.Counter({
  name: 'paymentservice_rpc_requests_total',
  help: 'Total number of gRPC requests processed, by method and status.',
  labelNames: ['method', 'status'],
  registers: [register]
});

const rpcDuration = new client.Histogram({
  name: 'paymentservice_rpc_duration_seconds',
  help: 'Latency of gRPC requests in seconds, by method.',
  labelNames: ['method'],
  buckets: [0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5],
  registers: [register]
});

const chargesTotal = new client.Counter({
  name: 'paymentservice_charges_total',
  help: 'Total number of successfully processed charges, by card type.',
  labelNames: ['card_type'],
  registers: [register]
});

const chargeAmountUSD = new client.Histogram({
  name: 'paymentservice_charge_amount_usd',
  help: 'Amount (in USD major units) charged per transaction.',
  buckets: [5, 10, 25, 50, 100, 250, 500, 1000],
  registers: [register]
});

const chargeErrorsTotal = new client.Counter({
  name: 'paymentservice_charge_errors_total',
  help: 'Total number of failed charge attempts, by error type.',
  labelNames: ['error_type'],
  registers: [register]
});

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

module.exports = { register, rpcRequestsTotal, rpcDuration, chargesTotal, chargeAmountUSD, chargeErrorsTotal, startMetricsServer };
