# Valkey Document Data Search (Go)

This project provides a Go-based CLI and API for generating, storing, indexing, and searching synthetic customer and event data in Valkey (the open-source Redis fork—yes, I'm on team fork, but Redis folks are cool too!) using ValkeyJSON and ValkeySearch modules.

## Features
- Generate random customer and event data
- Store data in Valkey as JSON documents
- Create ValkeySearch indexes for fast querying
- Search for customers/events by identifiers
- Print random customer/event JSON for testing
- REST API and CLI interface

## Prerequisites
- Go 1.18+
- Valkey server with ValkeyJSON and ValkeySearch modules enabled (Redis also works, but hey, forks are spicy!)

## Project Structure
- `cmd/cli/main.go` — CLI entry point
- `cmd/api/main.go` — API server entry point
- `internal/faker/` — Random data generation library
- `internal/valkeyutil/` — Valkey and ValkeySearch utilities

---

# CLI

## Building and Running (Recommended)

This project uses a Makefile to build both the CLI and API binaries and place them in the `bin/` directory.

### Build both binaries
```sh
make
```

This will generate:
- `bin/redis-document-cli` — CLI binary
- `bin/redis-document-api` — API server binary

### Clean binaries
```sh
make clean
```

### Running the CLI
Set the Valkey/Redis URL with the `REDIS_URL` environment variable (optional, defaults to `redis://localhost:6379/0`):

```sh
export REDIS_URL=redis://localhost:6379/0
```

Run the CLI:
```sh
./bin/redis-document-cli <command> [args]
```

## CLI Commands

### Generate Customers
Generate and store N customers in Valkey/Redis:
```sh
./bin/redis-document-cli generate_customers 1000
```

### Generate Events
Generate and store N events in Valkey/Redis:
```sh
./bin/redis-document-cli generate_events 1000
```

### Example Records Stored in Redis

Below are examples of the JSON structure for both customer and event records as stored in Redis. These are helpful for understanding the data model and for generating your own test data.

#### Example: Customer Record

```json
{
  "query_time_ms": 704,
  "result": {
    "confidenceScore": 0.9336259178648032,
    "createdAt": "2011-04-05T02:47:52Z",
    "customerId": "1026be83-1ee2-405f-8ec2-e96ff1a0447f",
    "deleted": 0,
    "identifiers": {
      "session_ids": [
        "1d586573-48bd-40d6-be99-0f55067e6aac",
        "c233b9da-5c9a-45fe-89e7-f995db960f59"
      ],
      "visitor_ids": [
        "06ced6eb-58a5-4f29-bba0-b28f5bda3e5b",
        "12e6ec8e-8435-4e22-827b-b0d366695054"
      ]
    },
    "merged": 1,
    "personalData": {
      "company": "Jurispect",
      "inferred_location": "Colorado Springs, Chad",
      "name": "Berneice Nikolaus",
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

#### Example: Event Record

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
      "cmec_contact_call_id": "call_EkcgN",
      "cmec_contact_chat_id": "chat_idXJv",
      "cmec_contact_external_id": "ext_wHzAe",
      "cmec_contact_form2lead_id": "f2l_KkarS",
      "cmec_contact_tickets_id": "ticket_fhChr",
      "cmec_visitor_id": "RCu"
    },
    "source": "rGnzBVZY",
    "timestamp": "2025-07-08T18:30:19Z",
    "visitor_data": {
      "behavior": {
        "interactions": [
          "scroll",
          "click_cta",
          "hover"
        ],
        "pages_viewed": 2,
        "time_on_site": 354
      },
      "device_info": {
        "device_type": "mobile",
        "ip_address": "192.168.133.97",
        "user_agent": "KsAslrZeXe"
      },
      "page_url": "https://sitels9t9.com/pagehe3q0mbdv3",
      "referrer": "/internal/path",
      "session_id": "qSo",
      "utm_params": {
        "utm_campaign": "campMhsER",
        "utm_medium": "medPstqW",
        "utm_source": "srcpFwi"
      },
      "visitor_id": "RCu"
    }
  }
}
```

### Create Indexes
Create RediSearch indexes for customers and events:
```sh
./bin/redis-document-cli create_indexes
```

### Search Customers
Search for customers by identifiers (e.g., email, phone, visitor_id):
```sh
./bin/redis-document-cli search_customers email=foo@bar.com phone=123456789
```

### Search Events
Search for events by identifiers (e.g., visitor_id, call_id, chat_id):
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

Or set environment variables for Redis and port:

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
          {"path": "$.primaryIdentifiers.cmec_visitor_id", "alias": "visitor_id", "type": "TEXT"}
        ]
      },
      {
        "name": "eventIdx",
        "fields": [
          {"path": "$.identifiers.cmec_visitor_id", "alias": "visitor_id", "type": "TEXT"},
          {"path": "$.identifiers.cmec_contact_call_id", "alias": "call_id", "type": "TEXT"},
          {"path": "$.identifiers.cmec_contact_chat_id", "alias": "chat_id", "type": "TEXT"},
          {"path": "$.identifiers.cmec_contact_external_id", "alias": "external_id", "type": "TEXT"},
          {"path": "$.identifiers.cmec_contact_form2lead_id", "alias": "form2lead_id", "type": "TEXT"},
          {"path": "$.identifiers.cmec_contact_tickets_id", "alias": "tickets_id", "type": "TEXT"}
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
   - Before running any performance tests, ensure your Redis database is populated with test data using the `generate_*` commands (e.g., `generate_customers`, `generate_events`).
   - See the [CLI Usage](#cli-usage) section below for instructions on generating and populating data.
2. **Extract sample data:**
   - Use the provided script to extract a 5% random sample of customer and event records (with their indexed fields) into CSV files. This step reads from the populated Redis database.
3. **Run performance tests:**
   - Only after extracting the sample CSVs can you run the k6 performance tests, as these files are required for generating the HTTP requests.

### 1. Extract Sample Data for Testing
If your database is already populated, use the provided script to extract a 5% random sample of customer and event records (with their indexed fields) into CSV files:

```sh
cd perf
./sample_redis_to_csv.sh customer
./sample_redis_to_csv.sh event
```

This will generate `customer_sample.csv` and `event_sample.csv` in the `perf/` directory, ready for performance testing.

**How the Sample CSVs Are Used for Testing**

The extracted sample CSV files serve as test data for generating a wide range of HTTP requests to the API. During performance testing, the k6 script will:

- Select one of the indexed fields (columns) from the CSV for each request.
- Use the field values from each record to construct search queries against the API.

This means the total number of possible different requests performed is equal to:

```
number of indexed fields × number of records in the CSV
```

This approach ensures the benchmark covers all indexed fields and a diverse set of real data from your database.


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

## License
MIT
