"""
Experiment E: Stress Test
What we're testing: Find the breaking point — ramp to 500+ users and watch for failures.
Config: 90% GET + 10% POST, ramp to 500 users, spawn rate 20/s
"""

import random
from locust import HttpUser, task, between


class StressUser(HttpUser):
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
                "sku": f"STRESS-{random.randint(1000, 9999)}",
                "manufacturer": "Stress Test Corp",
                "category_id": random.randint(1, 50),
                "weight": random.randint(0, 5000),
                "some_other_id": random.randint(1, 100),
            },
            name="/products/[id]/details",
        )
