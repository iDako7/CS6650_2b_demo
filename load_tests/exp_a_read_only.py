"""
Experiment A: Read-Only
What we're testing: Baseline GET performance with NO lock contention (RLock only).
Config: 100% GET, ramp 50 -> 200 users, spawn rate 10/s
"""

import random
from locust import HttpUser, task, between


class ReadOnlyUser(HttpUser):
    wait_time = between(1, 3)

    @task
    def get_product(self):
        product_id = random.randint(1, 3)
        self.client.get(f"/products/{product_id}", name="/products/[id]")
