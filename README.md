# Go-Loyalty - Random Go/Gin-Boilerplate Playground for Loyalty Points Management

[![Go Build](https://github.com/dydanz/go-playground/actions/workflows/go.yml/badge.svg)](https://github.com/dydanz/go-playground/actions/workflows/go.yml) [![Docker](https://github.com/dydanz/go-playground/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/dydanz/go-playground/actions/workflows/docker-publish.yml) [![Codacy Badge](https://app.codacy.com/project/badge/Grade/00bb0e4faf7c4cd493b14ff5d587ea68)](https://app.codacy.com/gh/dydanz/go-playground/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![Jekyll / GitHub Pages](https://github.com/dydanz/go-playground/actions/workflows/jekyll-gh-pages.yml/badge.svg)](https://github.com/dydanz/go-playground/actions/workflows/jekyll-gh-pages.yml)

A (playground) RESTful API service, built with Go (Gin framework) that handles Loyalty Points management with PostgreSQL for data persistence and Redis for caching.

It receives inbound transactions from merchant clients, classifies transactions, and generates points for each transaction based on predefined program rules and constraints.

#### Disclaimer
As designated for my personal research AI-generated code, most of the codes are less-caffeinated-machine-generated

But feel free to fork, clone or whatever you want at your own risk. 
For questions or professional inquiries: [Linkedin](https://www.linkedin.com/in/dandi-diputra/)

Tech stack:
- Go 1.21+ with Gin framework
- PostgreSQL & Redis
- Docker & Docker Compose, comply with GKE.
- Generated using Cursor 0.8 + Claude 3 Sonnet in VSCode

## Features

- RESTful API endpoints for user management (CRUD operations)
- PostgreSQL database with Replication for CQRS Approach, Redis for caching Session Management.
--- future plans, table archival will be implemented to introduce advanced hot-cold data separation.
- Password hashing using bcrypt
--- not sure SSO will be implemented, but it's a good idea to implement it.
- Swagger, is easier to check and verify your work.
- Locust for Load Testing (Python knowledge required to implement the test scenario.
- Integrated with static file config to start and build your JS-based web-app. Hello World page provided!
- Docker support for local development and deployment to GKE
- Environment-based configuration

## Prerequisites

Before you begin, ensure you have installed:
- Go 1.16 or later (1.21 is recommended)
- Docker and Docker Compose
- Git
- Python 3.10 or later
- A good PC/laptop is needed because Docker will be hungry!

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/dydanz/go-playground.git
```

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

Connect to PostgreSQL and create the user's table:

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

You can find the API endpoints in the Swagger UI at [Swagger URL](http://localhost:8080/swagger/index.html)

or run it with `$ python3 test.py` to see the end-to-end tests result.

## Project Structure
```
go-playground/
├── cmd/
│   └── api/
│       └── main.go
├── server/
│   ├── bootstrap/ # Initializes the application
│   │   ├── database.go
│   │   ├── repository.go
│   │   ├── router.go
│   │   └── service.go
│   ├── config/
│   │   └── config.go
│   ├── docs/
│   │   ├── docs.go
│   │   ├── swagger.json
│   │   └── swagger.yaml
│   ├── domain/
│   │   ├── auth.go
│   │   ├── event_log.go
│   │   ├── event_log_repository.go
│   │   ├── interfaces.go
│   │   ├──  ... (more DTO/interface files)
│   ├── handler/
│   │   ├── auth_handler.go
│   │   ├── merchant_handler.go
│   │   ├── ping.go
│   │   ├──  ... (more router handler files)
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── csrf.go
│   │   └── session.go
│   ├── migrations/
│   │   ├── 000001_create_users_table.up.sql
│   │   ├── 000002_add_auth_tables.up.sql
│   │   ├── 000003_add_auth_token_unique_constraint.up.sql
│   │   └── ... (more migration files)
│   ├── mocks/ # Mock files for testing
│   │   ├── repository/
│   │   └── service/
│   ├── repository/
│   │   ├── postgres/
│   │   └── redis/
│   ├── service/
│   │   ├── auth_service.go
│   │   ├── merchant_service.go
│   │   ├── points_service.go
│   │   ├── redemption_service.go
│   │   └── ... (more service files)
│   └── util/
│       └── entitlement.go
├── pkg/
│   ├── channel/
│   │   └── pubsub_channel.go
│   ├── database/
│   │   ├── migration.go
│   │   ├── postgres.go
│   │   └── redis.go
│   └── kafka/
│       └── kafka.go
├── web/ # Static files for simple web pages, fool proof you can add reactjs etc within go/gin project.
│   ├── assets/
│   │   ├── css/
│   │   ├── js/
│   └── pages/
│       ├── components/ # reusable components
│       ├── dashboard.html
│       ├── sign-in.html
│       ├── sign-up.html
│       └── transactions.html
├── locust-test/ # Locust load testing files, run separately under python-venv
│   ├── common/
│   ├── locustfile.py # locust load-test file
│   ├── requirements.txt
│   └── itgtest.py # run integration test
├── docker-compose.yml
├── Dockerfile
├── .env
└── go.mod
```

## Development

### Running Tests

First, you'll need to add the testify package to your dependencies. Run:
``` $ go get github.com/stretchr/testify ```
To run the test, you can use these commands:
Run all tests in the project:
``` $ go test ./... -v ```
Run specific test file:
``` $ go test ./server/service/user_service_test.go ```
Run with verbose output:
``` $ go test -v ./server/service ```


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
swag init -g cmd/api/main.go -o server/docs
```

3. Access the documentation at [Swagger URL](http://localhost:8080/swagger/index.html) after starting the server

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

Migrations are stored in `server/migrations/` directory:
```
server/migrations/
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
migrate -path server/migrations -database "postgres://postgres:your_password_here@localhost:5432/go_cursor?sslmode=disable" up

# Rollback migrations
migrate -path server/migrations -database "postgres://postgres:your_password_here@localhost:5432/go_cursor?sslmode=disable" down
```

### Creating New Migrations

To create a new migration:
```bash
migrate create -ext sql -dir server/migrations -seq add_new_feature
```

This creates two files:
- `XXXXXX_add_new_feature.up.sql`: Forward migration
- `XXXXXX_add_new_feature.down.sql`: Rollback migration

### Migration Commands

```bash
# Apply all migrations
migrate -path server/migrations -database ${DATABASE_URL} up

# Rollback last migration
migrate -path server/migrations -database ${DATABASE_URL} down 1

# Rollback all migrations
migrate -path server/migrations -database ${DATABASE_URL} down

# Force a specific version
migrate -path server/migrations -database ${DATABASE_URL} force VERSION
```

### Common Migration Issues

1. "Dirty" Database State
```bash
migrate -path server/migrations -database ${DATABASE_URL} force VERSION
```

2. Version Mismatch
```bash
# Check current version
migrate -path server/migrations -database ${DATABASE_URL} version
```

3. Connection Issues
- Verify database credentials
- Check if database is running
- Ensure correct permissions

### Locust Load Testing
This project uses Locust for API load testing and performance analysis.

1. Create a Python virtual environment:
```bash
cd locust-test
python3 -m venv venv
```

2. Activate virtual environment and install Locust:
```bash
# Activate virtual environment
source venv/bin/activate

# Install Locust package
pip install locust
```
3. Start Locust server:
```bash
# Make sure you're in the locust-test directory with activated virtual environment
locust -f locustfile.py <test-class-name>
```
4. Access Locust Web Interface:
- Open your browser and navigate to http://localhost:8089
- Set number of users, spawn rate, and target host
- Click "Start swarming" to begin the load test
Note: Ensure your API server is running before starting the load test.
