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

## Build the CLI

```sh
cd redis-document-data-search/cmd/cli
go build -o redisdocsearch
```

## CLI Usage

Set the Redis URL with the `REDIS_URL` environment variable (optional, defaults to `redis://localhost:6379/0`):

```sh
export REDIS_URL=redis://localhost:6379/0
```

Run the CLI with one of the following commands:

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

### API Endpoints

- **Store Customers:**  
  `POST /store_customers?count=1000`
- **Store Events:**  
  `POST /store_events?count=1000`
- **Create Indexes:**  
  `POST /create_indexes`
- **Search Customers:**  
  `GET /search_customers?email=foo@bar.com`
- **Search Events:**  
  `GET /search_events?visitor_id=123`
- **Random Customer:**  
  `GET /random_customer`
- **Random Event:**  
  `GET /random_event`
- **Health Check:**  
  `GET /healthz`

## License

MIT
