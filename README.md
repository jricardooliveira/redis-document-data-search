# Redis Document Data Search (Go)

This project provides a Go-based CLI for generating, storing, indexing, and searching synthetic customer and event data in Redis using RedisJSON and RediSearch. It is a full migration of the original Python script (`proj.py`).

## Features
- Generate random customer and event data
- Store data in Redis as JSON documents
- Create RediSearch indexes for fast querying
- Search for customers/events by identifiers
- Print random customer/event JSON for testing

## Prerequisites
- Go 1.18+
- Redis server with RedisJSON and RediSearch modules enabled

## Build
```sh
cd redis-document-data-search
go build -o redisdocsearch main.go
```

## Usage
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

## Redis URL
The CLI uses `redis://localhost:6379/0` by default. Edit `main.go` to change the Redis URL if needed.

## Project Structure
- `main.go` - CLI entry point
- `faker/` - Random data generation library
- `redisutil/` - Redis and RediSearch utilities

## License
MIT
