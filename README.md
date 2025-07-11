# Redis Document Data Search (Go)

This project provides a Go-based CLI and API for generating, storing, indexing, and searching synthetic customer and event data in Redis using RedisJSON and RediSearch.

## Features
- Generate random customer and event data
- Store data in Redis as JSON documents
- Create RediSearch indexes for fast querying
- Search for customers/events by identifiers
- Print random customer/event JSON for testing
- REST API and CLI interface

## Prerequisites
- Go 1.18+
- Redis server with RedisJSON and RediSearch modules enabled

## Project Structure
- `cmd/cli/main.go` — CLI entry point
- `cmd/api/main.go` — API server entry point
- `internal/faker/` — Random data generation library
- `internal/redisutil/` — Redis and RediSearch utilities

---

# CLI

## Building the CLI

```sh
cd redis-document-data-search/cmd/cli
go build -o redisdocsearch
```

## Running the CLI

Set the Redis URL with the `REDIS_URL` environment variable (optional, defaults to `redis://localhost:6379/0`):

```sh
export REDIS_URL=redis://localhost:6379/0
```

## CLI Commands

### Store Customers
Generate and store N customers in Redis:
```sh
./redisdocsearch store_customers 1000
```

### Store Events
Generate and store N events in Redis:
```sh
./redisdocsearch store_events 1000
```

### Create Indexes
Create RediSearch indexes for customers and events:
```sh
./redisdocsearch create_indexes
```

### Search Customers
Search for customers by identifiers (e.g., email, phone, visitor_id):
```sh
./redisdocsearch search_customers email=foo@bar.com phone=123456789
```

### Search Events
Search for events by identifiers (e.g., visitor_id, call_id, chat_id):
```sh
./redisdocsearch search_events visitor_id=123 call_id=abc
```

### Print Random Customer/Event
Print a random customer or event as JSON (for testing):
```sh
./redisdocsearch customer
./redisdocsearch event
```

---

# API

## Running the API Server

Build and run the API server:

```sh
cd redis-document-data-search/cmd/api
go run main.go
```

Or set environment variables for Redis and port:

```sh
REDIS_URL=redis://localhost:6379/0 API_PORT=8080 go run main.go
```

## API Endpoints

### 1. Store Customers
- **Method:** `POST`
- **Path:** `/store_customers`
- **Query Parameters:**
  - `count` (optional, default: `1000`): Number of customers to generate and store.
- **Response:**
  ```json
  { "status": "ok", "stored": 1000 }
  ```
  or
  ```json
  { "error": "error message" }
  ```

### 2. Store Events
- **Method:** `POST`
- **Path:** `/store_events`
- **Query Parameters:**
  - `count` (optional, default: `1000`): Number of events to generate and store.
- **Response:**
  ```json
  { "status": "ok", "stored": 1000 }
  ```
  or
  ```json
  { "error": "error message" }
  ```

### 3. Create Indexes
- **Method:** `POST`
- **Path:** `/create_indexes`
- **Response:**
  ```json
  { "status": "ok" }
  ```
  or
  ```json
  { "customerIdx": "error (if any)", "eventIdx": "error (if any)" }
  ```

### 4. Search Customers
- **Method:** `GET`
- **Path:** `/search_customers`
- **Query Parameters:**
  - Any combination of customer identifiers (e.g., `email`, `phone`, `visitor_id`).
  - `limit` (optional, default: `10`): Max results.
  - `offset` (optional, default: `0`): Offset for pagination.
- **Response:**
  ```json
  { "results": [ ... ], "query_time_μs": 1234 }
  ```
  or
  ```json
  { "error": "error message", "query_time_μs": 1234 }
  ```

### 5. Search Events
- **Method:** `GET`
- **Path:** `/search_events`
- **Query Parameters:**
  - Any combination of event identifiers (e.g., `visitor_id`, `call_id`, `chat_id`).
  - `limit` (optional, default: `10`): Max results.
  - `offset` (optional, default: `0`): Offset for pagination.
- **Response:**
  ```json
  { "results": [ ... ], "query_time_μs": 1234 }
  ```
  or
  ```json
  { "error": "error message", "query_time_μs": 1234 }
  ```

### 6. Get Random Event
- **Method:** `GET`
- **Path:** `/random_event`
- **Response:**
  A random event as JSON, or
  ```json
  { "error": "no events found" }
  ```

### 7. Get Random Customer
- **Method:** `GET`
- **Path:** `/random_customer`
- **Response:**
  A random customer as JSON, or
  ```json
  { "error": "no customers found" }
  ```

### 8. Health Check
- **Method:** `GET`
- **Path:** `/healthz`
- **Response:**
  ```json
  { "status": "ok" }
  ```

---

## License
MIT
