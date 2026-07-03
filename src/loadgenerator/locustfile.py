#!/usr/bin/python
#
# Copyright 2018 Google LLC
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

import logging
import os
import random
from locust import FastHttpUser, TaskSet, between, events
from prometheus_client import Counter, Histogram, start_http_server
from faker import Faker
import datetime
fake = Faker()

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger("loadgenerator")

requests_total = Counter(
    'loadgenerator_requests_total',
    'Total number of requests issued by simulated users, by request name and outcome.',
    ['name', 'outcome'])

request_duration_seconds = Histogram(
    'loadgenerator_request_duration_seconds',
    'Latency of requests issued by simulated users, in seconds, by request name.',
    ['name'])

users_total = Counter(
    'loadgenerator_users_started_total',
    'Total number of simulated users started.')


@events.init.add_listener
def on_locust_init(environment, **kwargs):
    port = int(os.environ.get('METRICS_PORT', '9090'))
    start_http_server(port)
    logger.info("serving Prometheus metrics on :%d/metrics", port)


@events.request.add_listener
def on_request(request_type, name, response_time, response_length, exception, **kwargs):
    outcome = 'error' if exception else 'success'
    requests_total.labels(name=name, outcome=outcome).inc()
    request_duration_seconds.labels(name=name).observe(response_time / 1000.0)

products = [
    '0PUK6V6EV0',
    '1YMWWN1N4O',
    '2ZYFJ3GM2N',
    '66VCHSJNUP',
    '6E92ZMYYFZ',
    '9SIQT8TOJO',
    'L9ECAV7KIM',
    'LS4PSXUNUM',
    'OLJCESPC7Z']

def index(l):
    l.client.get("/")

def setCurrency(l):
    currencies = ['EUR', 'USD', 'JPY', 'CAD', 'GBP', 'TRY']
    l.client.post("/setCurrency",
        {'currency_code': random.choice(currencies)})

def browseProduct(l):
    l.client.get("/product/" + random.choice(products))

def viewCart(l):
    l.client.get("/cart")

def addToCart(l):
    product = random.choice(products)
    quantity = random.randint(1,10)
    l.client.get("/product/" + product)
    l.client.post("/cart", {
        'product_id': product,
        'quantity': quantity})
    logger.info("simulated user added product=%s quantity=%d to cart", product, quantity)

def empty_cart(l):
    l.client.post('/cart/empty')
    logger.info("simulated user emptied cart")

def checkout(l):
    addToCart(l)
    current_year = datetime.datetime.now().year+1
    l.client.post("/cart/checkout", {
        'email': fake.email(),
        'street_address': fake.street_address(),
        'zip_code': fake.zipcode(),
        'city': fake.city(),
        'state': fake.state_abbr(),
        'country': fake.country(),
        'credit_card_number': fake.credit_card_number(card_type="visa"),
        'credit_card_expiration_month': random.randint(1, 12),
        'credit_card_expiration_year': random.randint(current_year, current_year + 70),
        'credit_card_cvv': f"{random.randint(100, 999)}",
    })
    logger.info("simulated user completed checkout")
    
def logout(l):
    l.client.get('/logout')  


class UserBehavior(TaskSet):

    def on_start(self):
        users_total.inc()
        index(self)

    tasks = {index: 1,
        setCurrency: 2,
        browseProduct: 10,
        addToCart: 2,
        viewCart: 3,
        checkout: 1}

class WebsiteUser(FastHttpUser):
    tasks = [UserBehavior]
    wait_time = between(1, 10)
