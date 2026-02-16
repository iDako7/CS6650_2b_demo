# Load Testing Results — CS 6650 Homework 5

## Test Environment

| Parameter | Value |
|---|---|
| **Target** | ECS Fargate (us-west-2) / Local server |
| **ECS Task** | 256 CPU units (0.25 vCPU), 512 MB memory |
| **Server** | Go + Gin, in-memory store with `sync.RWMutex`, 3 seeded products |
| **Test Client** | Locust 2.43.2, running locally on macOS |
| **Date** | 2026-02-15 |

---

## How to Run Tests

### Step 1: Start the target server

**Option A — Test against local server:**

```bash
cd src/
go run .
# Server starts at http://localhost:8080
```

**Option B — Test against ECS (deployed via Terraform):**

```bash
# Get the public IP (run from terraform/ directory)
aws ec2 describe-network-interfaces \
  --network-interface-ids $(
    aws ecs describe-tasks \
      --cluster $(terraform output -raw ecs_cluster_name) \
      --tasks $(
        aws ecs list-tasks \
          --cluster $(terraform output -raw ecs_cluster_name) \
          --service-name $(terraform output -raw ecs_service_name) \
          --query 'taskArns[0]' --output text
      ) \
      --query "tasks[0].attachments[0].details[?name=='networkInterfaceId'].value" \
      --output text
  ) \
  --query 'NetworkInterfaces[0].Association.PublicIp' \
  --output text

# Note the IP, e.g. 44.242.237.19
```

### Step 2: Launch Locust with web UI

Run one of the commands below from the **repo root** (`CS6650_2b_demo/`).
Then open **http://localhost:8089** in your browser.

#### Against local server (`localhost:8080`)

```bash
# Experiment A: Read-Only
locust -f load_tests/exp_a_read_only.py --host=http://localhost:8080

# Experiment B: Write-Only
locust -f load_tests/exp_b_write_only.py --host=http://localhost:8080

# Experiment C: Mixed (HttpUser)
locust -f load_tests/exp_c_mixed.py --host=http://localhost:8080

# Experiment D: Mixed (FastHttpUser)
locust -f load_tests/exp_d_mixed_fast.py --host=http://localhost:8080

# Experiment E: Stress Test
locust -f load_tests/exp_e_stress.py --host=http://localhost:8080
```

#### Against ECS deployed server (replace `<PUBLIC-IP>`)

```bash
# Experiment A: Read-Only
locust -f load_tests/exp_a_read_only.py --host=http://<PUBLIC-IP>:8080

# Experiment B: Write-Only
locust -f load_tests/exp_b_write_only.py --host=http://<PUBLIC-IP>:8080

# Experiment C: Mixed (HttpUser)
locust -f load_tests/exp_c_mixed.py --host=http://<PUBLIC-IP>:8080

# Experiment D: Mixed (FastHttpUser)
locust -f load_tests/exp_d_mixed_fast.py --host=http://<PUBLIC-IP>:8080

# Experiment E: Stress Test
locust -f load_tests/exp_e_stress.py --host=http://<PUBLIC-IP>:8080
```

### Step 3: Configure in the Locust web UI

After opening **http://localhost:8089**, fill in the fields:

| Experiment | Number of users | Spawn rate | Run time |
|---|---|---|---|
| **A: Read-Only** | 200 | 10 | 60s (or leave blank for manual stop) |
| **B: Write-Only** | 200 | 10 | 60s |
| **C: Mixed HttpUser** | 200 | 10 | 60s |
| **D: Mixed FastHttpUser** | 200 | 10 | 60s |
| **E: Stress Test** | 500 | 20 | 60s |

Click **Start swarming**, wait at least 60 seconds, then go to the **Charts** tab and take screenshots.

**Note:** You must stop Locust (`Ctrl+C` in terminal) between experiments — each experiment uses a different locust file.

---

## Experiment Descriptions

| Experiment | What It Tests | Locust File |
|---|---|---|
| **A: Read-Only** | Baseline GET performance, no write lock contention (RLock only) | `exp_a_read_only.py` |
| **B: Write-Only** | Write lock overhead (exclusive Lock on every request) | `exp_b_write_only.py` |
| **C: Mixed (HttpUser)** | Real-world 90% read / 10% write pattern | `exp_c_mixed.py` |
| **D: Mixed (FastHttpUser)** | Same as C but with connection pooling — is the bottleneck client or server? | `exp_d_mixed_fast.py` |
| **E: Stress Test** | Find the breaking point — ramp to 500+ users | `exp_e_stress.py` |

---

## Results Summary

> Fill in with your own numbers from the Locust web UI after running each experiment.
> The tables below are pre-filled with reference data from a CLI run.

### Overall Comparison Table

| Experiment | Users | User Class | Workload | Total Requests | Failures | RPS | Avg (ms) | Median (ms) | p95 (ms) | p99 (ms) | Max (ms) |
|---|---|---|---|---|---|---|---|---|---|---|---|
| **A: Read-Only** | 200 | HttpUser | 100% GET | 5,030 | 0 (0%) | 84.54 | 32 | 29 | 63 | 130 | 1,129 |
| **B: Write-Only** | 200 | HttpUser | 100% POST | 5,052 | 0 (0%) | 84.48 | 31 | 29 | 61 | 110 | 251 |
| **C: Mixed** | 200 | HttpUser | 90/10 | 5,042 | 0 (0%) | 84.40 | 31 | 29 | 65 | 120 | 257 |
| **D: Mixed Fast** | 200 | FastHttpUser | 90/10 | 5,079 | 0 (0%) | 85.00 | 28 | 27 | 53 | 82 | 251 |
| **E: Stress** | 500 | HttpUser | 90/10 | 11,957 | 0 (0%) | 200.05 | 31 | 28 | 68 | 120 | 259 |

### Per-Endpoint Breakdown

| Experiment | Endpoint | Requests | Avg (ms) | Median (ms) | p95 (ms) | p99 (ms) |
|---|---|---|---|---|---|---|
| **A** | GET /products/[id] | 5,030 | 32 | 29 | 63 | 130 |
| **B** | POST /products/[id]/details | 5,052 | 31 | 29 | 61 | 110 |
| **C** | GET /products/[id] | 4,555 | 31 | 29 | 65 | 120 |
| **C** | POST /products/[id]/details | 487 | 32 | 29 | 71 | 160 |
| **D** | GET /products/[id] | 4,568 | 28 | 27 | 54 | 82 |
| **D** | POST /products/[id]/details | 511 | 28 | 27 | 47 | 73 |
| **E** | GET /products/[id] | 10,752 | 31 | 28 | 68 | 120 |
| **E** | POST /products/[id]/details | 1,205 | 30 | 28 | 63 | 120 |

---

## Analysis

### Q1: "Which operations will be most common in a real-world scenario?"

**GETs dominate.** In a real e-commerce system, customers browse products far more often than admins add or update product details. Experiment C confirms the 90/10 split is realistic — 4,555 GETs vs 487 POSTs. This read-heavy pattern directly influences the choice of concurrency control.

### Q2: "How does that impact the data structure?"

**A read-heavy workload favors `sync.RWMutex` over `sync.Mutex`.**

Comparing Experiments A (read-only) and B (write-only):

| Metric | A: Read-Only | B: Write-Only |
|---|---|---|
| Avg latency | 32 ms | 31 ms |
| p95 latency | 63 ms | 61 ms |
| RPS | 84.54 | 84.48 |

In this test, read and write performance are nearly identical. This is because the in-memory hashmap operations are extremely fast (sub-microsecond), so the lock contention overhead is negligible compared to the ~28 ms network round-trip between the Locust client (local machine) and the ECS server (us-west-2). The network latency dominates, masking any difference in lock behavior.

However, `sync.RWMutex` is still the correct choice: under higher concurrency or with slower operations (e.g., database queries), concurrent reads via `RLock()` would significantly outperform serialized access via `Mutex.Lock()`.

### Q3: "Why did many people not see a difference between HttpUser and FastHttpUser?"

Comparing Experiments C (HttpUser) and D (FastHttpUser):

| Metric | C: HttpUser | D: FastHttpUser |
|---|---|---|
| Avg latency | 31 ms | 28 ms |
| Median latency | 29 ms | 27 ms |
| p95 latency | 65 ms | 53 ms |
| p99 latency | 120 ms | 82 ms |
| RPS | 84.40 | 85.00 |

FastHttpUser shows a **small improvement** (~10% lower at p95, ~32% lower at p99), but overall RPS is nearly identical. This confirms that **the server (and network) is the bottleneck, not the test client**. FastHttpUser's connection pooling saves a few milliseconds of connection setup overhead, which shows up at higher percentiles, but the server's response time is the dominant factor. When the server is fast enough (e.g., local testing with no network hop), the difference would be more pronounced.

### Q4: Stress Test — Where is the breaking point?

Experiment E ramped to **500 users** (2.5x the other experiments):

| Metric | C: 200 users | E: 500 users |
|---|---|---|
| RPS | 84.40 | 200.05 |
| Avg latency | 31 ms | 31 ms |
| p95 latency | 65 ms | 68 ms |
| Failures | 0 (0%) | 0 (0%) |

**The server handled 500 users with zero failures and nearly identical latency.** RPS scaled linearly from ~84 to ~200 (proportional to user count with `wait_time=between(1,3)`). The server did not reach its breaking point at 500 concurrent users — the in-memory Go server with Gin is extremely efficient for this workload. The bottleneck is the Locust client's `wait_time` (1-3s between requests per user), not the server capacity.

To find the true breaking point, you would need to either:
- Remove or reduce `wait_time` to generate more requests per user
- Increase to 1000+ users
- Run Locust in distributed mode across multiple machines

---

## Key Takeaways

1. **Network latency dominates**: The ~28 ms round-trip to us-west-2 dwarfs any server-side processing differences (in-memory operations take microseconds).
2. **RWMutex is the right choice**: Even though read vs write performance looks similar here (masked by network latency), RWMutex allows concurrent reads which matters under real load with faster network paths.
3. **Server is the bottleneck, not the client**: HttpUser vs FastHttpUser shows minimal difference, confirming the bottleneck is server/network, not client connection overhead.
4. **Server is robust**: Zero failures across all experiments, even at 500 concurrent users. The Go + Gin + in-memory hashmap combination is very efficient for this API.
