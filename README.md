# ChatPress Analytics Microservice

A Go-based microservice for handling analytics data in the ChatPress system.

## Features

- API usage cost tracking
- Monthly usage statistics
- Overall status metrics
- JWT authentication
- PostgreSQL integration

## API Endpoints

- `GET /analytics/api-usage-cost` - Get API usage and cost metrics
- `GET /analytics/monthly-usage` - Get monthly usage trends
- `GET /analytics/overall-status` - Get overall system status
- `GET /health` - Health check endpoint

## Environment Variables

- `DATABASE_URL` - PostgreSQL connection string
- `JWT_SECRET` - Secret key for JWT token validation
- `PORT` - Service port (default: 8083)

## Running the Service

1. Install dependencies:
```bash
go mod tidy
```

2. Set environment variables:
```bash
export DATABASE_URL="postgres://user:password@localhost/chatpress?sslmode=disable"
export JWT_SECRET="your-secret-key"
export PORT="8083"
```

3. Run the service:
```bash
go run cmd/main.go
```

## Docker

Build and run with Docker:

```bash
docker build -t chatpress-analytics .
docker run -p 8083:8083 chatpress-analytics
```

## Architecture

The service follows a clean architecture pattern:

- `cmd/` - Application entry point
- `api/v1/` - HTTP handlers and routes
- `db/` - Database connection and repositories
- `models/` - Data structures
- `services/` - Business logic
- `config/` - Configuration management