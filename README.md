# SoCode - Distributed Log Aggregator

A high-performance distributed log aggregation system built with Go, featuring automated code review capabilities and deployment assistance.

## Features

- **Distributed Log Aggregation**: Collect, process, and store logs from multiple sources
- **High Performance**: Built with Go for optimal performance and concurrency
- **Multiple Protocols**: Support for HTTP REST API and gRPC
- **Persistent Storage**: PostgreSQL backend for reliable log storage
- **Caching Layer**: Redis integration for fast data retrieval
- **Code Review Automation**: Automated code review and deployment assistant
- **Scalable Architecture**: Designed for horizontal scaling

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Log Sources   │───▶│    SoCode       │───▶│   PostgreSQL    │
│                 │    │   (HTTP/gRPC)   │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │      Redis      │
                       │   (Cache Layer) │
                       └─────────────────┘
```

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 13+
- Redis 6+
- Docker (optional, for containerized deployment)

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/krishnaGauss/SoCode.git
cd SoCode
```

### 2. Environment Configuration

Create a `.env` file in the root directory:

```bash
cp .env.example .env
```

Edit the `.env` file with your configuration:

```env
# Server Configuration
SERVER_HOST=localhost
SERVER_PORT=8080
GRPC_PORT=9090

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=logs
DB_USER=postgres
DB_PASSWORD=postgres
DB_SSL_MODE=disable

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

### 3. Database Setup

#### Option A: Using Docker Compose (Recommended)

```bash
docker-compose up -d postgres redis
```

#### Option B: Manual Setup

**PostgreSQL:**
```sql
CREATE DATABASE logs;
CREATE USER postgres WITH PASSWORD 'postgres';
GRANT ALL PRIVILEGES ON DATABASE logs TO postgres;
```

**Redis:**
```bash
redis-server
```

### 4. Install Dependencies

```bash
go mod tidy
```

### 5. Run the Application

```bash
go run main.go
```

The application will start with:
- HTTP API server on `http://localhost:8080`
- gRPC server on `localhost:9090`

## Installation Options

### Development Setup

1. **Local Development:**
   ```bash
   # Install dependencies
   go mod download
   
   # Run with hot reload (requires air)
   go install github.com/cosmtrek/air@latest
   air
   ```

2. **Using Make:**
   ```bash
   make dev          # Start development server
   make build        # Build binary
   make test         # Run tests
   make clean        # Clean build artifacts
   ```

### Production Deployment

#### Docker Deployment

1. **Build Docker Image:**
   ```bash
   docker build -t socode:latest .
   ```

2. **Run with Docker Compose:**
   ```bash
   docker-compose up -d
   ```

#### Manual Deployment

1. **Build Binary:**
   ```bash
   CGO_ENABLED=0 GOOS=linux go build -o socode main.go
   ```

2. **Deploy:**
   ```bash
   # Copy binary to server
   scp socode user@server:/opt/socode/
   
   # Run with systemd (create socode.service)
   sudo systemctl start socode
   sudo systemctl enable socode
   ```

## API Documentation

### HTTP REST API

#### Health Check
```bash
GET /health
```

#### Log Ingestion
```bash
POST /api/v1/logs
Content-Type: application/json

{
  "timestamp": "2025-07-05T10:30:00Z",
  "level": "info",
  "message": "Application started",
  "service": "web-server",
  "metadata": {
    "user_id": "12345",
    "request_id": "req-abc123"
  }
}
```

#### Query Logs
```bash
GET /api/v1/logs?service=web-server&level=error&limit=100
```

### gRPC API

The gRPC service definitions are available in `proto/` directory. Generate client code:

```bash
# Install protoc-gen-go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Generate code
make proto
```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_HOST` | HTTP server host | `localhost` |
| `SERVER_PORT` | HTTP server port | `8080` |
| `GRPC_PORT` | gRPC server port | `9090` |
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_NAME` | Database name | `logs` |
| `DB_USER` | Database username | `postgres` |
| `DB_PASSWORD` | Database password | `postgres` |
| `DB_SSL_MODE` | SSL mode for PostgreSQL | `disable` |
| `REDIS_HOST` | Redis host | `localhost` |
| `REDIS_PORT` | Redis port | `6379` |
| `REDIS_PASSWORD` | Redis password | `` |
| `REDIS_DB` | Redis database number | `0` |

### Configuration File

Alternatively, use a YAML configuration file:

```yaml
# config.yaml
server:
  host: localhost
  port: 8080
  grpc_port: 9090

database:
  host: localhost
  port: 5432
  name: logs
  username: postgres
  password: postgres
  ssl_mode: disable

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0
```

## Development

### Project Structure

```
SoCode/
├── cmd/                    # Application entry points
├── internal/               # Private application code
│   ├── config/            # Configuration management
│   ├── handler/           # HTTP handlers
│   ├── service/           # Business logic
│   └── storage/           # Data access layer
├── pkg/                   # Public packages
├── proto/                 # Protocol buffer definitions
├── migrations/            # Database migrations
├── docker-compose.yml     # Docker services
├── Dockerfile            # Container build file
├── Makefile              # Build automation
└── README.md             # This file
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./internal/storage -v
```

### Code Style

This project follows Go best practices:
- Use `gofmt` for formatting
- Use `golint` for linting
- Follow effective Go guidelines

```bash
# Format code
go fmt ./...

# Lint code
golangci-lint run
```

## Monitoring

### Metrics

SoCode exposes Prometheus metrics on `/metrics` endpoint:

- `socode_logs_ingested_total` - Total number of logs ingested
- `socode_requests_duration_seconds` - Request duration histogram
- `socode_database_connections` - Active database connections

### Logging

Application logs are structured and can be configured via environment:

```bash
LOG_LEVEL=info
LOG_FORMAT=json
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Write tests for new features
- Follow Go conventions
- Update documentation
- Add proper error handling

## Troubleshooting

### Common Issues

1. **Database Connection Error:**
   ```
   Error: pq: password authentication failed
   ```
   **Solution:** Check your database credentials in `.env` file.

2. **Redis Connection Error:**
   ```
   Error: dial tcp: connection refused
   ```
   **Solution:** Ensure Redis is running and accessible.

3. **Port Already in Use:**
   ```
   Error: bind: address already in use
   ```
   **Solution:** Change ports in configuration or stop conflicting services.

### Debug Mode

Enable debug logging:

```bash
export LOG_LEVEL=debug
go run main.go
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- Create an issue on GitHub for bugs or feature requests
- Join our discussions for general questions
- Check the [documentation](docs/) for detailed guides

## Acknowledgments

- Built with [Go](https://golang.org/)
- Uses [PostgreSQL](https://www.postgresql.org/) for persistence
- Caching with [Redis](https://redis.io/)
- gRPC for high-performance communication

---

**Made with ❤️ by [krishnaGauss](https://github.com/krishnaGauss)**
Readme Generated by ChatGPT :)

