# Cinema Food & Beverage Recommender
# Apple-grade developer build automation

BINARY_NAME=cinema-recommender
BIN_DIR=bin
PORT?=8080

.PHONY: all build test run clean help

all: test build

build:
	@echo "🍎 Compiling standalone binary..."
	go build -ldflags="-s -w" -o $(BIN_DIR)/$(BINARY_NAME).exe ./cmd/api
	@echo "✅ Build complete: $(BIN_DIR)/$(BINARY_NAME).exe"

test:
	@echo "🧪 Running full test suite..."
	go test -v -race ./...

run: build
	@echo "🚀 Launching HTTP server on port $(PORT)..."
	./$(BIN_DIR)/$(BINARY_NAME).exe -server -port :$(PORT)

cli-kandy: build
	@echo "🍿 Evaluating deal for Kandy..."
	./$(BIN_DIR)/$(BINARY_NAME).exe -city Kandy -budget 2500 -party 2

cli-promos: build
	@echo "👑 Listing active Scope Privilege promotions..."
	./$(BIN_DIR)/$(BINARY_NAME).exe -list-promos

clean:
	@echo "🧹 Cleaning artifacts..."
	rm -rf $(BIN_DIR)

help:
	@echo "Available commands:"
	@echo "  make build      - Compile binary into $(BIN_DIR)/"
	@echo "  make test       - Run all unit and policy tests"
	@echo "  make run        - Compile and start HTTP server"
	@echo "  make cli-kandy  - Run sample CLI calculation"
	@echo "  make cli-promos - List promotional campaigns"
	@echo "  make clean      - Remove build artifacts"
