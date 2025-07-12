BIN_DIR := bin
API_SRC := cmd/api/main.go
CLI_SRC := cmd/cli/main.go
API_BIN := $(BIN_DIR)/redis-document-api
CLI_BIN := $(BIN_DIR)/redis-document-cli

.PHONY: all build clean

all: build

build: $(API_BIN) $(CLI_BIN)

$(BIN_DIR):
	@mkdir -p $(BIN_DIR)

$(API_BIN): $(API_SRC) | $(BIN_DIR)
	go build -o $(API_BIN) ./cmd/api

$(CLI_BIN): $(CLI_SRC) | $(BIN_DIR)
	go build -o $(CLI_BIN) ./cmd/cli

clean:
	rm -rf $(BIN_DIR)
