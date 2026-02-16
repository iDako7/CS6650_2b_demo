"""
Experiment C: Mixed Realistic (HttpUser)
What we're testing: Real-world 90% read / 10% write pattern with RWMutex contention.
Config: 90% GET + 10% POST, ramp 50 -> 200 users, spawn rate 10/s
"""

import random
from locust import HttpUser, task, between


class MixedUser(HttpUser):
    wait_time = between(1, 3)

    @task(9)
    def get_product(self):
        product_id = random.randint(1, 3)
        self.client.get(f"/products/{product_id}", name="/products/[id]")

    @task(1)
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
