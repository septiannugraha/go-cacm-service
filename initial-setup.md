# Initial Setup Instructions

Since the pb package is generated locally, you need to follow these steps:

## 1. First, create the pb directory and a placeholder file:

```bash
cd /mnt/c/Users/villager/code/go-cacm-service
mkdir -p pb
echo "package pb" > pb/doc.go
```

## 2. Clean up go.mod to remove the current dependencies:

```bash
rm go.mod go.sum
```

## 3. Re-initialize the module:

```bash
go mod init github.com/septiannugraha/go-cacm-service
```

## 4. Add required dependencies manually:

```bash
go get google.golang.org/protobuf
go get github.com/go-sql-driver/mysql
go get github.com/lib/pq
go get github.com/microsoft/go-mssqldb
```

## 5. Install protoc-gen tools:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## 6. Generate the protobuf files:

```bash
protoc --go_out=pb --go_opt=paths=source_relative \
       --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
       proto/query_result.proto
```

## 7. Now restore the original packager.go file with the pb imports

## 8. Finally, run go mod tidy:

```bash
go mod tidy
```

The issue occurs because the `pb` package is a local generated package, not a remote module. By creating the pb directory and generating the files first, Go will recognize it as a local package.