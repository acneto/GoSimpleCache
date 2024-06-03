BINARY_CACHE=bin/simple_cache
SRC_CACHE=cache/cmd/main.go

BINARY_API=bin/api
SRC_API=api/cmd/main.go

all: build

build:
	@echo "Building the projects..."
	go build -o $(BINARY_CACHE) $(SRC_CACHE)
	go build -o $(BINARY_API) $(SRC_API)

run-cache-engine: build
	@echo "Running cache-engine..."
	./$(BINARY_CACHE)

run-web-api: build
	@echo "Running the web-api..."
	./$(BINARY_API)

test:
	go test ./...