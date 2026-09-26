# Auth Service

Auth Service written in Go.

## Features

- User registration
- User login
- JWT access tokens
- Refresh tokens
- Refresh token rotation
- Logout
- Role-based access control
- PostgreSQL
- Redis
- Swagger documentation
- Unit tests
- HTTP handler tests
- Docker support
- GitHub Actions CI

## Tech Stack

- Go
- Gin
- PostgreSQL
- Redis
- JWT
- bcrypt
- Docker
- Swagger
- GitHub Actions

## Run with Docker

```bash
docker compose up --build
```

The application will be available at:

```text
http://localhost:8080
```

Swagger:

```text
http://localhost:8080/swagger/index.html
```

Health check:

```text
GET /health
```

## Environment Variables

Example:

```env
APP_PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5433/auth?sslmode=disable
REDIS_ADDR=localhost:6379
JWT_SECRET=change-me
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=168h
```

## Tests

Run all tests:

```bash
go test ./...
```

## Build

```bash
go build ./...
```
