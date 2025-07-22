# Go CACM Service

A Go implementation of the SISKEUDES data synchronization service, migrated from the C# SituwassaSync project. This service extracts financial data from SISKEUDES SQL Server databases, packages it into Protocol Buffers, and uploads to a central processing service.

## Features

- Connects to multiple SISKEUDES SQL Server databases
- Extracts revenue, expense, and financing data
- Packages data into Protocol Buffer format
- Uploads binary files to processing service
- Handles budget revision logic (KdPosting)

## Prerequisites

- Go 1.21 or higher
- Protocol Buffer compiler (protoc)
- SQL Server access with appropriate credentials

## Installation

1. Clone the repository:
```bash
git clone https://github.com/septiannugraha/go-cacm-service.git
cd go-cacm-service
```

2. Install dependencies:
```bash
make deps
```

3. Generate Protocol Buffer code:
```bash
make proto
```

4. Build the application:
```bash
make build
```

## Configuration

Create a `situwassa.conf` file in the project root:

```json
{
  "server": "localhost",
  "integrated_security": false,
  "user_id": "your_username",
  "password": "your_password",
  "databases": ["PEMDA_2024", "PEMDA_2025"]
}
```

## Usage

Run the application:
```bash
make run
# or
./bin/cacm-service
```

## Project Structure

```
go-cacm-service/
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── database/            # Database connection and queries
│   │   ├── mssql.go        # SQL Server client
│   │   └── queries.go      # SQL query constants
│   ├── models/             # Data models
│   │   └── models.go       # Go structs for data
│   ├── packager/           # Protocol Buffer packaging
│   │   └── packager.go     # Converts data to protobuf
│   └── uploader/           # HTTP upload functionality
│       └── uploader.go     # Uploads files to service
├── proto/
│   └── query_result.proto  # Protocol Buffer definitions
├── pb/                     # Generated protobuf code (after make proto)
├── go.mod                  # Go module definition
├── Makefile               # Build commands
└── README.md              # This file
```

## Development

### Adding New Report Types

1. Add the model struct in `internal/models/models.go`
2. Create a packaging method in `internal/packager/packager.go`
3. Add the query logic in `internal/database/mssql.go`
4. Update the main loop in `cmd/main.go`

### Running Tests

```bash
go test ./...
```

### Generating Protocol Buffers

After modifying `proto/query_result.proto`:
```bash
make proto
```

## Migration Notes from C#

Key differences from the C# implementation:
- Uses native Go SQL drivers instead of Entity Framework
- Protocol Buffer handling via google.golang.org/protobuf
- Simplified configuration (no DI container)
- Same business logic for KdPosting handling

## License

[Your License Here]