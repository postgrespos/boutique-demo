#!/usr/bin/python
#
# Copyright 2026 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import functools
import os
import time

from prometheus_client import Counter, Histogram, start_http_server

rpc_requests_total = Counter(
    'emailservice_rpc_requests_total',
    'Total number of gRPC requests processed, by method and status.',
    ['method', 'status'])

rpc_duration_seconds = Histogram(
    'emailservice_rpc_duration_seconds',
    'Latency of gRPC requests in seconds, by method.',
    ['method'])

emails_sent_total = Counter(
    'emailservice_emails_sent_total',
    'Total number of order confirmation emails sent.')


def start_metrics_server(logger):
  port = int(os.environ.get('METRICS_PORT', '9090'))
  start_http_server(port)
  logger.info('serving Prometheus metrics on :{}/metrics'.format(port))


def instrument(method_name):
  """Decorator that records request count and latency for a gRPC handler."""
  def decorator(func):
    @functools.wraps(func)
    def wrapper(self, request, context, *args, **kwargs):
      start = time.time()
      status = 'success'
      try:
        return func(self, request, context, *args, **kwargs)
      except Exception:
        status = 'error'
        raise
      finally:
        rpc_duration_seconds.labels(method=method_name).observe(time.time() - start)
        rpc_requests_total.labels(method=method_name, status=status).inc()
    return wrapper
  return decorator
