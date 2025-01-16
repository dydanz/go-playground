# Go-Playground - Random Go/Gin-Boilerplate Playground 

[![Go](https://github.com/dydanz/go-playground/actions/workflows/go.yml/badge.svg)](https://github.com/dydanz/go-playground/actions/workflows/go.yml)

A (playground) RESTful API service built with Go (Gin framework) that handles user management with PostgreSQL for data persistence and Redis for caching.

#### Disclaimer
As designated for my personal research AI-generated code, ALL OF THESE CODE ARE GENERATED AUTOMATICALLY! 

But feel free to fork, clone or whatever you want at your own risk. 
For questions or professional inquiries: [Linkedin](https://www.linkedin.com/in/dandi-diputra/)

Tech stack:
- Go 1.21+ with Gin framework
- PostgreSQL & Redis
- Docker & Docker Compose
- Generated using Cursor 0.8 + Claude 3 Sonnet in VSCode

## Features

- RESTful API endpoints for user management (CRUD operations)
- PostgreSQL database with UUID as primary key
- Password hashing using bcrypt
- Redis caching
- Docker support for local development
- Environment-based configuration

## Prerequisites

Before you begin, ensure you have installed:
- Go 1.16 or later
- Docker and Docker Compose
- Git

## Getting Started

### 1. Clone the Repository

### 2. Set Up Environment Variables

Create a `.env` file in the project root:
```env
# PostgreSQL Settings
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres123
DB_NAME=go_cursor

# Redis Settings
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=redis123
```

### 3. Start Dependencies (PostgreSQL and Redis)

```bash
docker-compose up -d
```

### 4. Install Go Dependencies

```bash
go mod tidy
```

### 5. Initialize Database

Connect to PostgreSQL and create the users table:

```bash
docker exec -it $(docker ps -qf "name=postgres") psql -U postgres -d go_cursor
```

Then run the following SQL:

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
```

### 6. Run the Application

```bash
go run cmd/api/main.go
```

The API will be available at `http://localhost:8080`

## API Endpoints

### Create User
```bash
curl -X POST http://localhost:8080/api/users \
-H "Content-Type: application/json" \
-d '{
    "email": "john@example.com",
    "password": "password123",
    "name": "John Doe",
    "phone": "1234567890"
}'
```

### Get All Users
```bash
curl http://localhost:8080/api/users
```

### Get User by ID
```bash
curl http://localhost:8080/api/users/{user_id}
```

### Update User
```bash
curl -X PUT http://localhost:8080/api/users/{user_id} \
-H "Content-Type: application/json" \
-d '{
    "name": "John Updated",
    "phone": "0987654321"
}'
```

### Delete User
```bash
curl -X DELETE http://localhost:8080/api/users/{user_id}
```

## Project Structure
```
go-cursor/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── domain/
│   │   └── user.go
│   ├── repository/
│   │   ├── postgres/
│   │   │   └── user_repository.go
│   │   └── redis/
│   │       └── cache_repository.go
│   ├── handler/
│   │   └── user_handler.go
│   └── service/
│       └── user_service.go
├── pkg/
│   └── database/
│       ├── postgres.go
│       └── redis.go
├── docker-compose.yml
├── .env
└── go.mod
```

## Development

### Running Tests

```bash
go test ./... -v
```

### Common Issues

1. Database Connection Issues
   - Check if PostgreSQL container is running: `docker ps`
   - Verify .env credentials match docker-compose.yml
   - Wait a few seconds after starting containers

2. Redis Connection Issues
   - Check if Redis container is running: `docker ps`
   - Verify Redis password in .env matches docker-compose.yml

3. "Module Not Found" Errors
   - Run `go mod tidy` to fix dependencies
   - Check if module name in imports matches go.mod

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| DB_HOST | PostgreSQL host | localhost |
| DB_PORT | PostgreSQL port | 5432 |
| DB_USER | PostgreSQL username | postgres |
| DB_PASSWORD | PostgreSQL password | postgres123 |
| DB_NAME | PostgreSQL database name | go_cursor |
| REDIS_HOST | Redis host | localhost |
| REDIS_PORT | Redis port | 6379 |
| REDIS_PASSWORD | Redis password | redis123 |

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details

## Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Go-Redis](https://github.com/go-redis/redis)
- [Lib/pq](https://github.com/lib/pq)

## API Documentation (Swagger)

The API documentation is available via Swagger UI at:
```
http://localhost:8080/swagger/index.html
```

### Generating Swagger Documentation

1. Install swag:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

2. Generate documentation:
```bash
swag init -g cmd/api/main.go -o internal/docs
```

3. Access the documentation at http://localhost:8080/swagger/index.html after starting the server

### Swagger Annotations

The API endpoints are documented using Swagger annotations in the handler files. Example:
```go
// @Summary Create a new user
// @Description Create a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Param user body domain.CreateUserRequest true "User details"
// @Success 201 {object} domain.User
// @Router /users [post]
```

## Database Migrations

This project uses `golang-migrate` for database schema management.

### Migration Files

Migrations are stored in `internal/migrations/` directory:
```
internal/migrations/
├── 000001_create_users_table.up.sql   # Create tables
└── 000001_create_users_table.down.sql  # Rollback changes
```

### Running Migrations

Migrations run automatically when the application starts. To run them manually:

1. Install golang-migrate:
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

2. Run migrations:
```bash
# Apply migrations
migrate -path internal/migrations -database "postgres://postgres:your_password_here@localhost:5432/go_cursor?sslmode=disable" up

# Rollback migrations
migrate -path internal/migrations -database "postgres://postgres:your_password_here@localhost:5432/go_cursor?sslmode=disable" down
```

### Creating New Migrations

To create a new migration:
```bash
migrate create -ext sql -dir internal/migrations -seq add_new_feature
```

This creates two files:
- `XXXXXX_add_new_feature.up.sql`: Forward migration
- `XXXXXX_add_new_feature.down.sql`: Rollback migration

### Migration Commands

```bash
# Apply all migrations
migrate -path internal/migrations -database ${DATABASE_URL} up

# Rollback last migration
migrate -path internal/migrations -database ${DATABASE_URL} down 1

# Rollback all migrations
migrate -path internal/migrations -database ${DATABASE_URL} down

# Force a specific version
migrate -path internal/migrations -database ${DATABASE_URL} force VERSION
```

### Common Migration Issues

1. "Dirty" Database State
```bash
migrate -path internal/migrations -database ${DATABASE_URL} force VERSION
```

2. Version Mismatch
```bash
# Check current version
migrate -path internal/migrations -database ${DATABASE_URL} version
```

3. Connection Issues
- Verify database credentials
- Check if database is running
- Ensure correct permissions
