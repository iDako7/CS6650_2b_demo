"""
Experiment B: Write-Only
What we're testing: Write lock overhead (exclusive Mutex lock on every request).
Config: 100% POST, ramp 50 -> 200 users, spawn rate 10/s
"""

import random
from locust import HttpUser, task, between


class WriteOnlyUser(HttpUser):
    wait_time = between(1, 3)

    @task
    def update_product(self):
        product_id = random.randint(1, 3)
        self.client.post(
            f"/products/{product_id}/details",
            json={
                "product_id": product_id,
                "sku": f"LOAD-TEST-{random.randint(1000, 9999)}",
                "manufacturer": "Load Test Corp",
                "category_id": random.randint(1, 50),
                "weight": random.randint(0, 5000),
                "some_other_id": random.randint(1, 100),
            },
            name="/products/[id]/details",
        )
