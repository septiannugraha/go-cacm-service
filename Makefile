.PHONY: proto build run clean

# Generate Go code from proto files
proto:
	@echo "Generating Go code from proto files..."
	@mkdir -p pb
	@protoc --go_out=pb --go_opt=paths=source_relative \
		--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
		proto/query_result.proto

# Build the application
build: proto
	@echo "Building application..."
	@go build -o bin/cacm-service cmd/main.go

# Run the application
run: build
	@echo "Running application..."
	@./bin/cacm-service

# Clean generated files
clean:
	@echo "Cleaning generated files..."
	@rm -rf pb bin

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest