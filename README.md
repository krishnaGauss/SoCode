
# SoCode - Distributed Log Aggregation Microservice

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)](https://github.com/krishnaGauss/SoCode)

SoCode is a high-performance distributed log aggregation microservice built with Go. It provides centralized collection, storage, and querying of logs from multiple services, making it easier to monitor and troubleshoot distributed systems.

## 🚀 Features

- **Distributed Log Aggregation**: Centralized collection and storage of logs from multiple services
- **Multi-Protocol Support**: HTTP REST API and gRPC endpoints for log ingestion
- **High Performance**: Built with Go for optimal concurrency and performance
- **Persistent Storage**: PostgreSQL backend for reliable data persistence
- **Fast Caching**: Redis integration for high-speed data retrieval
- **Flexible Querying**: Search and filter logs by service, level, time range, and more
- **Real-time & Historical Analysis**: Support for both live monitoring and historical log analysis
- **Scalable Architecture**: Designed for horizontal scaling across multiple instances

## 🎯 What SoCode Does

### **1. Log Collection & Ingestion**
- Receives log messages from various applications/services via HTTP REST API or gRPC
- Each log entry contains information like timestamp, severity level, message, service name, and metadata
- Supports high-throughput ingestion from multiple concurrent sources

### **2. Data Storage**
- Stores all collected logs in PostgreSQL database for persistence and querying
- Uses Redis as a caching layer for fast retrieval of recent or frequently accessed logs

### **3. Log Querying & Retrieval**
- Provides APIs to search and filter logs by various criteria (service, log level, time range, etc.)
- Supports both real-time and historical log analysis

## 🌐 Real-World Use Case Example

Imagine you have a web application with multiple services:
- Frontend web server
- Authentication service
- Payment processing service
- Database service

Instead of checking logs on each individual server, SoCode allows you to:

1. **Send all logs to one place**: Each service sends its logs to SoCode
2. **Centralized monitoring**: View all logs from a single dashboard
3. **Easy troubleshooting**: Search for errors across all services at once
4. **Historical analysis**: Query past logs to identify patterns or issues

### Example Scenario:
```
Frontend Service → SoCode ← Auth Service
                      ↓
Payment Service → PostgreSQL ← Database Service
                      ↓
                   Redis Cache
```

When a user reports a payment failure, instead of checking 4 different servers, you can:
- Query SoCode for all logs related to that user's session
- See the complete flow: frontend → auth → payment → database
- Identify exactly where the failure occurred
- Analyze patterns if it's a recurring issue

## 📋 Prerequisites

Before installing SoCode, ensure you have:

- **Go 1.21+** - [Download here](https://golang.org/dl/)
- **PostgreSQL 13+** - [Installation guide](https://www.postgresql.org/download/)
- **Redis 6+** - [Installation guide](https://redis.io/download)
- **Docker** (optional) - [Get Docker](https://docs.docker.com/get-docker/)
- **Git** - [Install Git](https://git-scm.com/downloads)

## 🛠️ Installation

### Method 1: Quick Start with Docker (Recommended)

1. **Clone the repository:**
   ```bash
   git clone https://github.com/krishnaGauss/SoCode.git
   cd SoCode
   ```

2. **Start services with Docker Compose:**
   ```bash
   docker-compose up -d
   ```

3. **Verify installation:**
   ```bash
   curl http://localhost:8080/health
   ```

### Method 2: Manual Installation

#### Step 1: Clone and Setup

```bash
# Clone repository
git clone https://github.com/krishnaGauss/SoCode.git
cd SoCode

# Install Go dependencies
go mod download
```

#### Step 2: Database Setup

**PostgreSQL Setup:**
```bash
# Connect to PostgreSQL
psql -U postgres

# Create database and user
CREATE DATABASE logs;
CREATE USER socode WITH PASSWORD 'your_secure_password';
GRANT ALL PRIVILEGES ON DATABASE logs TO socode;
\q
```

**Redis Setup:**
```bash
# Start Redis server
redis-server

# Verify Redis is running
redis-cli ping
# Should return: PONG
```

#### Step 3: Environment Configuration

Create a `.env` file in the project root:

```bash
# Server Configuration
SERVER_HOST=localhost
SERVER_PORT=8080
GRPC_PORT=9090

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=logs
DB_USER=socode
DB_PASSWORD=your_secure_password
DB_SSL_MODE=disable

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Application Settings
LOG_LEVEL=info
LOG_FORMAT=json
```

#### Step 4: Build and Run

```bash
# Build the application
go build -o socode main.go

# Run the application
./socode
```

### Method 3: Development Setup

For development with hot reload:

```bash
# Install air for hot reload
go install github.com/cosmtrek/air@latest

# Start development server
air

# Or use make commands
make dev
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `SERVER_HOST` | HTTP server bind address | `localhost` | No |
| `SERVER_PORT` | HTTP server port | `8080` | No |
| `GRPC_PORT` | gRPC server port | `9090` | No |
| `DB_HOST` | PostgreSQL host | `localhost` | Yes |
| `DB_PORT` | PostgreSQL port | `5432` | No |
| `DB_NAME` | Database name | `logs` | Yes |
| `DB_USER` | Database username | `postgres` | Yes |
| `DB_PASSWORD` | Database password | - | Yes |
| `DB_SSL_MODE` | PostgreSQL SSL mode | `disable` | No |
| `REDIS_HOST` | Redis host | `localhost` | Yes |
| `REDIS_PORT` | Redis port | `6379` | No |
| `REDIS_PASSWORD` | Redis password | - | No |
| `REDIS_DB` | Redis database number | `0` | No |
| `LOG_LEVEL` | Application log level | `info` | No |
| `LOG_FORMAT` | Log format (json/text) | `json` | No |

### Configuration File

Alternative YAML configuration (`config.yaml`):

```yaml
server:
  host: "localhost"
  port: 8080
  grpc_port: 9090

database:
  host: "localhost"
  port: 5432
  name: "logs"
  username: "socode"
  password: "your_secure_password"
  ssl_mode: "disable"

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0

logging:
  level: "info"
  format: "json"
```

## 🔌 API Usage

### HTTP REST API

#### Health Check
```bash
curl http://localhost:8080/health
```

#### Log Ingestion
```bash
curl -X POST http://localhost:8080/api/v1/logs \
  -H "Content-Type: application/json" \
  -d '{
    "timestamp": "2025-07-05T10:30:00Z",
    "level": "info",
    "message": "User login successful",
    "service": "auth-service",
    "metadata": {
      "user_id": "12345",
      "ip_address": "192.168.1.100",
      "request_id": "req-abc123"
    }
  }'
```

#### Query Logs
```bash
# Get recent logs
curl "http://localhost:8080/api/v1/logs?limit=100"

# Filter by service and level
curl "http://localhost:8080/api/v1/logs?service=auth-service&level=error"

# Time range query
curl "http://localhost:8080/api/v1/logs?from=2025-07-05T00:00:00Z&to=2025-07-05T23:59:59Z"

# Search by message content
curl "http://localhost:8080/api/v1/logs?search=payment%20failed"

# Combined filters
curl "http://localhost:8080/api/v1/logs?service=payment-service&level=error&limit=50"
```

### gRPC API

Generate client code:
```bash
# Install protoc tools
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Generate code
protoc --go_out=. --go-grpc_out=. proto/*.proto
```

## 🚨 Troubleshooting

### Common Issues

#### 1. Database Connection Failed

**Error:**
```
Failed to initialize PostgreSQL: pq: password authentication failed for user "postgres"
```

**Solutions:**
- Verify PostgreSQL is running: `pg_isready -U postgres`
- Check username/password in `.env` file
- Ensure database exists: `psql -U postgres -l`
- Reset PostgreSQL password if needed:
  ```bash
  sudo -u postgres psql
  ALTER USER postgres PASSWORD 'newpassword';
  ```

#### 2. SSL Mode Error

**Error:**
```
pq: unsupported sslmode "%!s(MISSING)"
```

**Solutions:**
- Set `DB_SSL_MODE=disable` for local development
- Use `DB_SSL_MODE=require` for production
- Valid SSL modes: `disable`, `require`, `verify-ca`, `verify-full`

#### 3. Redis Connection Error

**Error:**
```
dial tcp 127.0.0.1:6379: connect: connection refused
```

**Solutions:**
- Start Redis server: `redis-server`
- Check Redis status: `redis-cli ping`
- Verify Redis port: `netstat -tlnp | grep :6379`
- Check Redis configuration: `redis-cli config get port`

#### 4. Port Already in Use

**Error:**
```
bind: address already in use
```

**Solutions:**
- Find process using port: `lsof -i :8080`
- Kill process: `kill -9 <PID>`
- Change port in configuration
- Use different port: `SERVER_PORT=8081`

#### 5. Go Module Issues

**Error:**
```
go: cannot find main module
```

**Solutions:**
- Initialize Go module: `go mod init github.com/krishnaGauss/SoCode`
- Download dependencies: `go mod download`
- Clean module cache: `go clean -modcache`

#### 6. Permission Denied

**Error:**
```
permission denied while trying to connect to the Docker daemon
```

**Solutions:**
- Add user to docker group: `sudo usermod -aG docker $USER`
- Restart terminal or run: `newgrp docker`
- Use sudo: `sudo docker-compose up -d`

#### 7. High Memory Usage

**Issue:** SoCode consuming too much memory

**Solutions:**
- Configure Redis memory limit: `redis-cli config set maxmemory 256mb`
- Implement log retention policies
- Add database indexes for better query performance:
  ```sql
  CREATE INDEX idx_logs_timestamp ON logs(timestamp);
  CREATE INDEX idx_logs_service ON logs(service);
  CREATE INDEX idx_logs_level ON logs(level);
  ```

## 🏗️ Development

### Project Structure

```
SoCode/
├── cmd/                    # Application entry points
├── internal/               # Private application code
│   ├── config/            # Configuration management
│   ├── handler/           # HTTP/gRPC handlers
│   ├── service/           # Business logic
│   └── storage/           # Data access layer
├── pkg/                   # Public packages
├── proto/                 # Protocol buffer definitions
├── migrations/            # Database migrations
├── scripts/               # Build and deployment scripts
├── docker-compose.yml     # Docker services
├── Dockerfile            # Container build file
├── Makefile              # Build automation
└── README.md             # This file
```

## 📊 Monitoring

### Metrics Endpoints

- **Health Check**: `GET /health`


### Available Metrics

- `socode_logs_ingested_total` - Total logs processed
- `socode_requests_duration_seconds` - Request duration histogram
- `socode_database_connections_active` - Active DB connections
- `socode_redis_operations_total` - Redis operations counter
- `socode_log_queries_total` - Total log queries processed

## 🔒 Security

### Production Checklist

- [ ] Change default passwords
- [ ] Enable SSL/TLS (`DB_SSL_MODE=require`)
- [ ] Set up firewall rules
- [ ] Configure Redis authentication
- [ ] Use environment variables for secrets
- [ ] Enable audit logging
- [ ] Regular security updates
- [ ] Implement rate limiting for log ingestion

## 🚀 Deployment

### Docker Compose Production

```yaml
version: '3.8'
services:
  socode:
    image: socode:latest
    ports:
      - "8080:8080"
      - "9090:9090"
    environment:
      - DB_HOST=postgres
      - REDIS_HOST=redis
    depends_on:
      - postgres
      - redis
    restart: unless-stopped

  postgres:
    image: postgres:15
    environment:
      - POSTGRES_DB=logs
      - POSTGRES_USER=socode
      - POSTGRES_PASSWORD=secure_password
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:
```


## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

**Made with ❤️ by [krishnaGauss](https://github.com/krishnaGauss)**

⭐ Star this repo if you find it helpful!
