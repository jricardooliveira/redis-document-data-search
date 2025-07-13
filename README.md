# Valkey Document Data Search (Go)

This project provides a Go-based CLI and API for generating, storing, indexing, and searching synthetic customer and event data in Valkey/Redis using ValkeyJSON and ValkeySearch modules.

---

## Table of Contents
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Project Structure](#project-structure)
- [Installation & Setup](#installation--setup)
- [Usage](#usage)
  - [CLI Usage](#cli-usage)
  - [API Usage](#api-usage)
- [System Resource Monitoring](#system-resource-monitoring)
- [Best Practices](#best-practices)
- [Performance Testing](#performance-testing)
- [Performance Results](#performance-results-mac-m3-32gb-ram)
- [License](#license)

---

## Features
- Generate random customer and event data
- Store data in Valkey/Redis as JSON documents
- Create ValkeySearch indexes for fast querying
- Search for customers/events by identifiers
- Print random customer/event JSON for testing
- REST API and CLI interface

## Prerequisites
- Go 1.18+
- Valkey/Redis server with ValkeyJSON and ValkeySearch modules enabled (Redis also works)

## Project Structure
- `cmd/cli/main.go` — CLI entry point
- `cmd/api/main.go` — API server entry point
- `internal/faker/` — Random data generation library
- `internal/valkeyutil/` — Valkey/Redis and ValkeySearch utilities
- `internal/monitor/` — Platform-specific resource limit logging utilities
- `scripts/monitor_resources.sh` — Live system resource monitoring script

---

## Installation & Setup

### Building Binaries
This project uses a Makefile to build both the CLI and API binaries and place them in the `bin/` directory.

```sh
make
```

This will generate:
- `bin/redis-document-cli` — CLI binary
- `bin/redis-document-api` — API server binary

### Clean Binaries
```sh
make clean
```

### Environment Variables
Set the Valkey/Redis URL with the `REDIS_URL` environment variable (optional, defaults to `redis://localhost:6379/0`):

```sh
export REDIS_URL=redis://localhost:6379/0
```

---

## Usage

## CLI Usage

Run the CLI:
```sh
./bin/redis-document-cli <command> [args]
```

#### CLI Commands
- **Generate Customers:**
  ```sh
  ./bin/redis-document-cli generate_customers 1000
  ```
- **Generate Events:**
  ```sh
  ./bin/redis-document-cli generate_events 1000
  ```
- **Create Indexes:**
  ```sh
  ./bin/redis-document-cli create_indexes
  ```
- **Search Customers:**
  ```sh
  ./bin/redis-document-cli search_customers email=foo@bar.com phone=123456789
  ```
- **Search Events:**
  ```sh
  ./bin/redis-document-cli search_events visitor_id=123 call_id=abc
  ```
- **Print Random Customer/Event:**
  ```sh
  ./bin/redis-document-cli customer
  ./bin/redis-document-cli event
  ```

## Example Records Stored in Valkey/Redis

- **Customer Record:**
  ```json
  {
    "query_time_ms": 704,
    "result": {
      "confidenceScore": 0.93,
      "createdAt": "2011-04-05T02:47:52Z",
      "customerId": "1026be83-1ee2-405f-8ec2-e96ff1a0447f",
      "deleted": 0,
      "identifiers": {
        "session_ids": ["..."],
        "visitor_id": "visitor_6z2kl"
      },
      "job": {
        "company": "Acme Inc.",
        "department": "Accounting",
        "title": "Engineer"
      },
      "primaryIdentifiers": {
        "email": "elroygleichner@denesik.net",
        "phone": "7885540765"
      },
      "updatedAt": "1907-04-25T05:53:38Z"
    }
  }
  ```
- **Event Record:**
  ```json
  {
    "query_time_ms": 1068,
    "result": {
      "data": {
        "cookie": "cookie_yTYxZnVY",
        "email": "FnAjVGtPjV@example.com",
        "phone": "mGJzxUDIAX"
      },
      "event_id": "evt_GxHerj",
      "event_type": "visitor_event",
      "identifiers": {
        "visitor_id": "visitor_6z2kl",
        "call_id": "call_6z2kl",
        "chat_id": "chat_6z2kl"
      },
      "timestamp": "2022-01-01T12:00:00Z"
    }
  }
  ```

> For more detailed sample records, see [`perf/customer_sample.csv`](perf/customer_sample.csv) and [`perf/event_sample.csv`](perf/event_sample.csv) after running the sample extraction script.


---

## API Usage

### API Endpoint Summary

| Method | Path                        | Description                              |
|--------|-----------------------------|------------------------------------------|
| POST   | /generate_customers         | Generate customers (count param)         |
| POST   | /generate_events            | Generate events (count param)            |
| POST   | /create_indexes             | Create RediSearch indexes                |
| GET    | /search_customers           | Search customers by identifiers          |
| GET    | /search_events              | Search events by identifiers             |
| GET    | /random_event               | Get a random event                       |
| GET    | /random_customer            | Get a random customer                    |
| GET    | /healthz                    | Health check endpoint                    |

#### Running the API Server
After building with `make`, run the API server binary:
```sh
./bin/redis-document-api
```
Or set environment variables for Valkey/Redis and port:
```sh
REDIS_URL=redis://localhost:6379/0 API_PORT=8080 ./bin/redis-document-api
```

#### API Endpoints
- **Generate Customers:**
  - `POST /generate_customers?count=1000`
- **Generate Events:**
  - `POST /generate_events?count=1000`
- **Create Indexes:**
  - `POST /create_indexes`
- **Search Customers:**
  - `GET /search_customers?email=foo@bar.com`
- **Search Events:**
  - `GET /search_events?visitor_id=123&call_id=abc`
- **Get Random Event:**
  - `GET /random_event`
- **Get Random Customer:**
  - `GET /random_customer`
- **Health Check:**
  - `GET /healthz`

See the original README for detailed request/response examples.

---

## System Resource Monitoring & Best Practices

To monitor open files, sockets, and resource usage while running the API or CLI (especially during load testing), use:
```sh
./scripts/monitor_resources.sh
```
This script prints:
- Current and max open files (system-wide and per-process)
- Open sockets to Redis
- Sockets in TIME_WAIT state

### Platform-Specific Resource Limit Logging

## Best Practices

### Redis Client Usage (Best Practice)
The API server uses a singleton Redis client with connection pooling for all HTTP requests.
**Do not create a new Redis client per request**—this prevents connection leaks and resource exhaustion.

The application logs OS resource limits (open files, processes) at startup using platform-specific code:
- On Linux: logs both open files and process limits.
- On macOS: logs open files (process limits are not available).

### Redis Client Usage (Best Practice)
The API server uses a singleton Redis client with connection pooling for all HTTP requests.
**Do not create a new Redis client per request**—this prevents connection leaks and resource exhaustion.

---
```sh
./bin/redis-document-cli search_events visitor_id=123 call_id=abc
```

### Print Random Customer/Event
Print a random customer or event as JSON (for testing):
```sh
./bin/redis-document-cli customer
./bin/redis-document-cli event
```

---

# API

## Running the API Server

After building with `make`, run the API server binary:

```sh
./bin/redis-document-api
```

Or set environment variables for Valkey/Redis and port:

```sh
REDIS_URL=redis://localhost:6379/0 API_PORT=8080 ./bin/redis-document-api
```

## API Endpoints

### 1. Generate Customers
- **Method:** `POST`
- **Path:** `/generate_customers`
- **Query Parameters:**
  - `count` (optional, default: `1000`): Number of customers to generate and store.
- **Note:** Generation is parallelized using half of available CPU cores for fast bulk data creation.
- **Example:**
  ```sh
  curl -X POST "http://localhost:8080/generate_customers?count=10000"
  ```
- **Response:**
  ```json
  { "status": "ok", "stored": 10000, "query_time_ms": 1234 }
  ```

### 2. Generate Events
- **Method:** `POST`
- **Path:** `/generate_events`
- **Query Parameters:**
  - `count` (optional, default: `1000`): Number of events to generate and store.
- **Note:** Generation is parallelized using half of available CPU cores for fast bulk data creation.
- **Example:**
  ```sh
  curl -X POST "http://localhost:8080/generate_events?count=10000"
  ```
- **Response:**
  ```json
  { "status": "ok", "stored": 10000, "query_time_ms": 1234 }
  ```

### 3. Create Indexes
- **Method:** `POST`
- **Path:** `/create_indexes`
- **Response:**
  ```json
  { "status": "ok", "query_time_ms": 7 }
  ```

> **Note:** All search and indexing features require Valkey/Redis to have the RedisJSON and RediSearch modules enabled. These modules are supported in both (but Valkey is the cool new kid on the block!).

### 4. Search Customers
- **Method:** `GET`
- **Path:** `/search_customers`
- **Query Parameters:**
  - Any combination of customer identifiers (e.g., `email`, `phone`, `visitor_id`).
  - `limit` (optional, default: `10`): Max results.
  - `offset` (optional, default: `0`): Offset for pagination.
- **Example:**
  ```sh
  curl "http://localhost:8080/search_customers?email=foo@bar.com"
  ```
- **Response:**
  ```json
  { "results": [ ... ], "query_time_ms": 1234 }
  ```

### 5. Search Events
- **Method:** `GET`
- **Path:** `/search_events`
- **Query Parameters:**
  - Any combination of event identifiers (e.g., `visitor_id`, `call_id`, `chat_id`).
  - `limit` (optional, default: `10`): Max results.
  - `offset` (optional, default: `0`): Offset for pagination.
- **Example:**
  ```sh
  curl "http://localhost:8080/search_events?call_id=call_6z2kl"
  ```
- **Response:**
  ```json
  { "results": [ ... ], "query_time_ms": 1234 }
  ```

### 6. Get Random Event
- **Method:** `GET`
- **Path:** `/random_event`
- **Response:**
  A random event as JSON, or
  ```json
  { "error": "no events found", "query_time_ms": 7 }
  ```

### 7. Get Random Customer
- **Method:** `GET`
- **Path:** `/random_customer`
- **Response:**
  A random customer as JSON, or
  ```json
  { "error": "no customers found", "query_time_ms": 7 }
  ```

### 8. Health Check
- **Method:** `GET`
- **Path:** `/healthz`
- **Response:**
  ```json
  {
    "status": "ok",
    "redis_url": "redis://localhost:6379/0",
    "db_index": "0",
    "customer_count": 10000,
    "event_count": 10000,
    "redis_memory": {
      "used_memory_bytes": 12345678,
      "used_memory_human": "12.3M"
    },
    "indexes": [
      {
        "name": "customerIdx",
        "fields": [
          {"path": "$.primaryIdentifiers.email", "alias": "email", "type": "TEXT"},
          {"path": "$.primaryIdentifiers.phone", "alias": "phone", "type": "TEXT"},
          {"path": "$.primaryIdentifiers.visitor_id", "alias": "visitor_id", "type": "TEXT"}
        ]
      },
      {
        "name": "eventIdx",
        "fields": [
          {"path": "$.identifiers.visitor_id", "alias": "visitor_id", "type": "TEXT"},
          {"path": "$.identifiers.call_id", "alias": "call_id", "type": "TEXT"},
          {"path": "$.identifiers.chat_id", "alias": "chat_id", "type": "TEXT"},
          {"path": "$.identifiers.external_id", "alias": "external_id", "type": "TEXT"},
          {"path": "$.identifiers.lead_id", "alias": "lead_id", "type": "TEXT"},
          {"path": "$.identifiers.tickets_id", "alias": "tickets_id", "type": "TEXT"}
        ]
      }
    ],
    "query_time_ms": 7
  }
  ```

---

## Performance Testing

You can easily benchmark API search performance using real data from your Valkey/Redis database!

**Prerequisites:**

1. **Populate the database:**
   - Before running any performance tests, ensure your Valkey/Redis database is populated with test data using the `generate_*` commands (e.g., `generate_customers`, `generate_events`).
   - See the [CLI Usage](#cli-usage) section below for instructions on generating and populating data.
2. **Extract sample data:**
   - Use the provided script to extract a 5% random sample of customer and event records (with their indexed fields) into CSV files. This step reads from the populated Valkey/Redis database.
3. **Run performance tests:**
   - Only after extracting the sample CSVs can you run the k6 performance tests, as these files are required for generating the HTTP requests.

### 1. Extract Sample Data for Testing
If your database is already populated, use the provided script to extract a 5% random sample of customer and event records (with their indexed fields) into CSV files:

### Generating Sample Data for Performance Testing

To benchmark API search performance with realistic data, you should first extract a random sample of customer and event documents from your Valkey/Redis database into CSV files.

> **Note:** The `sample_to_csv` CLI command relies on the API server being already running and accessible (default: http://localhost:8080). Start the API with:
> ```sh
> ./bin/redis-document-api
> ```

The recommended way to extract samples is using the built-in CLI command:


```sh
bin/redis-document-cli sample_to_csv --percent 5 --output perf/customer_sample.csv --type customer
```
Example output:
```
Sampled 5000 records out of 100000. Output written to perf/customer_sample.csv
```

```sh
bin/redis-document-cli sample_to_csv --percent 5 --output perf/event_sample.csv --type event
```
Example output:
```
Sampled 5000 records out of 100000. Output written to perf/event_sample.csv
```

- `--percent 5` controls what fraction of your data is sampled (e.g., 5%).
- `--output` specifies the CSV file to write.
- `--type` must be either `customer` or `event`.

This creates `perf/customer_sample.csv` and `perf/event_sample.csv` ready for use with performance testing tools like k6. You can adjust the sample size as needed.

**Alternative:** If you prefer, you can also use the provided shell script:

```sh
cd perf
./sample_redis_to_csv.sh customer
./sample_redis_to_csv.sh event
```

This will generate `customer_sample.csv` and `event_sample.csv` in the `perf/` directory.

**How the Sample CSVs Are Used for Testing**

The extracted sample CSV files serve as test data for generating a wide range of HTTP requests to the API. During performance testing, the k6 script will:

- Select one of the indexed fields (columns) from the CSV for each request.
- Use the field values from each record to construct search queries against the API.

This means the total number of possible different requests performed is equal to:

```
number of indexed fields × number of records in the CSV
```

This approach ensures the benchmark covers all indexed fields and a diverse set of real data from your Valkey/Redis database.


### 2. Install k6 (if needed)
k6 is a modern open-source load testing tool. Install it with:

- **macOS (Homebrew):**
  ```sh
  brew install k6
  ```
- **Linux (Debian/Ubuntu):**
  ```sh
  sudo apt update && sudo apt install -y gnupg ca-certificates
  curl -s https://dl.k6.io/key.gpg | sudo apt-key add -
  echo "deb https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
  sudo apt update && sudo apt install k6
  ```
- Or see: https://k6.io/docs/getting-started/installation/

### 3. Run the Performance Test
With your API running and sample CSVs in place, run:

```sh
cd perf
k6 run k6_search_from_csv.js
```

#### Example: Staged Ramp-Up with Web Dashboard

To simulate a gradual increase in load (e.g., 10 users for 10 seconds, then add 10 more every 30 seconds until reaching 50 users), you must configure the stages in your k6 script (not via CLI):

At the top of your `perf/k6_search_from_csv.js` file, set:

```js
export const options = {
  stages: [
    { duration: '10s', target: 10 },
    { duration: '30s', target: 20 },
    { duration: '30s', target: 30 },
    { duration: '30s', target: 40 },
    { duration: '30s', target: 50 },
    // Optionally hold at 50:
    // { duration: '30s', target: 50 },
  ],
};
```

Then, run your test with the web dashboard enabled:

```sh
cd perf
k6 run --out=web-dashboard k6_search_from_csv.js
```

- This will:
  - Start with 10 users for 10 seconds.
  - Ramp up to 20 users over the next 30 seconds, then to 30, 40, and finally 50 users, each over 30 seconds.
  - The `--out=web-dashboard` flag launches a local web dashboard (usually at http://localhost:5665/) so you can monitor the test in real time.

You can adjust the stage durations and user counts as needed for your own testing scenarios.

> **Note:** The `--stages` CLI flag is not supported in most k6 versions. Always configure ramp-up stages in your script's `options` block.

The test will:
- Randomly alternate between `/search_customers` and `/search_events` API endpoints
- Pick real, indexed field values from your sampled CSVs for each request
- Simulate concurrent users and measure latency, throughput, and errors

You’ll get a detailed report on how your search endpoints perform with realistic data and queries!

---

## Performance Results (Mac M3, 32GB RAM)

The following results were obtained running the k6 performance test on a Mac M3 with 32GB of memory:

- **Test Data Volume:**
  - 50,000 customers
  - 2,000,000 events
- **Memory Used by Redis/Valkey:** ~7.4 GB
- **Throughput:**
  - ~2,522,939 HTTP requests processed in 1m40s (~25,000 requests/sec sustained)
- **Latency:**
  - Average: 1.32ms
  - Median: 1.29ms
  - 90th percentile: 2.17ms
  - 95th percentile: 2.49ms
  - Max: 1.83s (rare outliers)
- **Success Rate:**
  - 99.99% of checks succeeded
  - Only 71 failed event searches out of 1,261,996 (0.0056%)
  - 0.00% HTTP request failures (71 out of 2,522,939)
- **Network:**
  - Data received: 5.1 GB
  - Data sent: 271 MB

These results demonstrate that the system can efficiently handle high-throughput, low-latency document search workloads at scale, with minimal errors and consistent performance on modern Apple silicon hardware.

---

## License
MIT
