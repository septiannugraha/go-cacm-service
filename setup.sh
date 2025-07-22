#!/bin/bash

echo "Setting up go-cacm-service..."

# Create pb directory
mkdir -p pb

# Create a dummy go file in pb to make it a valid package
cat > pb/doc.go << EOF
// Package pb contains generated protocol buffer code.
package pb
EOF

# Install protoc dependencies
echo "Installing protoc-gen-go..."
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Download dependencies first (without pb)
echo "Downloading dependencies..."
go mod download

# Generate protobuf files
echo "Generating protobuf files..."
protoc --go_out=pb --go_opt=paths=source_relative \
       --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
       proto/query_result.proto

# Now run go mod tidy
echo "Running go mod tidy..."
go mod tidy

echo "Setup complete!"