# Quizer - Quiz Application Backend

A comprehensive quiz application backend built with Go, Gin framework, and PostgreSQL. The application provides a REST API for managing quiz topics, conducting quiz sessions, and tracking user progress.

## 📖 Table of Contents

- [Features](#features)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Installation & Setup](#installation--setup)
- [Configuration](#configuration)
- [Database Migrations](#database-migrations)
- [Running the Application](#running-the-application)
- [API Documentation](#api-documentation)
- [Usage Examples](#usage-examples)
- [Development](#development)
- [Contributing](#contributing)

## ✨ Features

- **User Authentication**: Registration, login, and JWT-based authentication
- **Quiz Management**: Topic-based quiz sessions with multiple difficulty levels
- **Question Management**: Support for multiple-choice questions with options
- **Session Tracking**: Complete quiz session management from start to completion
- **Progress Tracking**: User quiz history and performance statistics
- **RESTful API**: Well-structured REST API with comprehensive endpoints
- **Database Migrations**: Structured database schema with migration support
- **Middleware Support**: CORS, authentication, logging, and recovery middleware

## 📁 Project Structure

```
quizer/
├── cmd/                          # Application entry points
├── datasource/                   # Database connection and configuration
│   └── database.go
├── internal/                     # Private application code
│   ├── config/                   # Configuration management
│   │   └── config.go
│   ├── handlers/                 # HTTP request handlers
│   │   ├── auth_handler.go       # Authentication endpoints
│   │   ├── quiz_handler.go       # Quiz management endpoints
│   │   └── user_handler.go       # User management endpoints
│   ├── middleware/               # HTTP middleware
│   │   ├── auth.go              # Authentication middleware
│   │   └── middleware.go        # Common middleware
│   ├── models/                   # Data models and structures
│   │   ├── auth.go              # Authentication models
│   │   ├── quiz_*.go            # Quiz-related models
│   │   ├── user.go              # User models
│   │   └── *.go                 # Other domain models
│   ├── routes/                   # Route definitions
│   │   └── routes.go
│   └── services/                 # Business logic layer
│       ├── auth_service.go       # Authentication business logic
│       ├── quiz_service.go       # Quiz business logic
│       ├── token_service.go      # Token management
│       └── user_service.go       # User business logic
├── migrations/                   # Database migration files
│   ├── *_create_users_table.*
│   ├── *_create_topics_table.*
│   ├── *_create_questions_table.*
│   ├── *_create_quiz_sessions_table.*
│   └── *.sql                    # Other migration files
├── pkg/                         # Public shared packages
│   └── utils/                   # Utility functions
│       ├── env.go               # Environment utilities
│       ├── logger.go            # Logging utilities
│       ├── password.go          # Password hashing
│       └── response.go          # HTTP response utilities
├── scripts/                     # Deployment and setup scripts
├── go.mod                       # Go module definition
├── go.sum                       # Go module checksums
├── main.go                      # Application entry point
├── Makefile                     # Build and development commands
└── README.md                    # This file
```

## 🔧 Prerequisites

Before running this application, make sure you have the following installed:

- **Go**: Version 1.25.0 or higher
- **PostgreSQL**: Version 12 or higher
- **golang-migrate**: For database migrations ([Installation Guide](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate))
- **Git**: For cloning the repository
- **Make**: For using the Makefile commands (optional)

## 🚀 Installation & Setup

### 1. Clone the Repository

```bash
git clone https://github.com/rifqicode/quizer.git
cd quizer
```

### 2. Install Dependencies

```bash
# Install Go dependencies
go mod download
go mod tidy

# Or using Makefile
make install
```

### 3. Environment Configuration

Create a `.env` file in the root directory:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=quizer_db
DB_SSLMODE=disable

# Database URL for migrations (required for golang-migrate)
DATABASE_URL=postgres://your_db_user:your_db_password@localhost:5432/quizer_db?sslmode=disable

# Server Configuration
SERVER_PORT=8080
SERVER_MODE=debug
SERVER_READ_TIMEOUT=30s
SERVER_WRITE_TIMEOUT=30s

# JWT Configuration
JWT_SECRET=your_jwt_secret_key_here
JWT_EXPIRES_IN=24h

# Application Configuration
APP_NAME=Quizer
APP_VERSION=1.0.0
```

### 4. Database Setup

```bash
# Create PostgreSQL database
createdb quizer_db

# Or connect to PostgreSQL and create manually
psql -U postgres
CREATE DATABASE quizer_db;
```

## 🗃️ Database Migrations

This project uses [golang-migrate](https://github.com/golang-migrate/migrate) for database schema management. Migration files are located in the `migrations/` directory.

### Install golang-migrate

```bash
# On macOS using Homebrew
brew install golang-migrate

# On Ubuntu/Debian
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/

# On Windows using Chocolatey
choco install migrate

# Or download directly from GitHub releases
# https://github.com/golang-migrate/migrate/releases
```

### Running Migrations

```bash
# Make sure DATABASE_URL is set in your .env file
# DATABASE_URL=postgres://username:password@localhost:5432/database_name?sslmode=disable

# Apply all pending migrations
make migrate-up

# Or run directly with migrate command
migrate -database "$DATABASE_URL" -path migrations up

# Rollback the last migration
make migrate-down

# Or rollback directly
migrate -database "$DATABASE_URL" -path migrations down 1

# Check migration status
migrate -database "$DATABASE_URL" -path migrations version
```

### Creating New Migrations

```bash
# Create a new migration using Makefile (interactive)
make migrate-create

# Or create directly with migrate command
migrate create -ext sql -dir migrations -seq your_migration_name
```

### Migration Files Structure

The project includes the following migration files:
- `20250821141439_create_users_table.*` - User accounts and authentication
- `20250821141440_create_access_tokens_table.*` - JWT token management
- `20250822100001_create_topics_table.*` - Quiz topics/categories
- `20250822100002_create_questions_table.*` - Quiz questions
- `20250822100003_create_question_options_table.*` - Multiple choice options
- `20250822100004_create_quiz_sessions_table.*` - Quiz session tracking
- `20250822100005_create_quiz_answers_table.*` - User answers
- `20250822100006_create_user_question_history_table.*` - Answer history
- `20250822100007_seed_topics_data.*` - Default topic data
- `20250822100008_add_deleted_at_columns.*` - Soft delete support
- `20250822100009_seed_database_questions.*` - Sample quiz questions

## 🏃‍♂️ Running the Application

### Development Mode

```bash
# Run directly
go run main.go

# Or using Makefile
make dev

# Or using the development command
make run
```

### Production Build

```bash
# Build the application
make build

# Run the built binary
make run

# Or run directly
./bin/quizer
```

The application will start on `http://localhost:8080` (or the port specified in your `.env` file).

## 📚 API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Authentication Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/auth/register` | Register a new user |
| POST | `/auth/login` | User authentication |
| GET | `/auth/me` | Get current user profile |
| POST | `/auth/logout` | Logout user |

### Quiz Management Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/quiz/topics` | Get all active topics |
| GET | `/quiz/topics/{id}` | Get topic by ID |
| GET | `/quiz/topics/{id}/stats` | Get topic question statistics |

### Quiz Session Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/quiz/sessions` | Start a new quiz session |
| GET | `/quiz/sessions/{id}` | Get quiz session details |
| PUT | `/quiz/sessions/{id}/abandon` | Abandon a quiz session |
| GET | `/quiz/sessions/{id}/next-question` | Get next question |
| POST | `/quiz/sessions/{id}/answer` | Submit an answer |
| POST | `/quiz/sessions/{id}/complete` | Complete quiz session |
| GET | `/quiz/sessions/{id}/results` | Get quiz results |

### User Statistics Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/quiz/my-sessions` | Get user quiz history |
| GET | `/quiz/my-stats` | Get user quiz statistics |

### Health Check Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Basic health check |
| GET | `/health/db` | Database connectivity check |

## 🔧 Usage Examples

### 1. User Registration

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

### 2. User Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

### 3. Start Quiz Session

```bash
curl -X POST http://localhost:8080/api/v1/quiz/sessions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "topic_id": 1,
    "total_questions": 10,
    "difficulty": "easy"
  }'
```

### 4. Submit Answer

```bash
curl -X POST http://localhost:8080/api/v1/quiz/sessions/1/answer \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "question_id": 1,
    "selected_option_id": 2
  }'
```

## 🛠️ Development

### Available Make Commands

```bash
make install         # Install dependencies
make dev            # Run in development mode
make build          # Build the application
make run            # Run the built binary
make test           # Run tests
make test-coverage  # Run tests with coverage
make clean          # Clean build artifacts
make fmt            # Format code
make lint           # Lint code

# Migration commands
make migrate-create # Create a new migration (interactive)
make migrate-up     # Apply all pending migrations
make migrate-down   # Rollback the last migration

# Docker commands
make docker-build   # Build Docker image
make docker-run     # Run Docker container
```

### Code Structure Guidelines

- **handlers/**: HTTP request handlers, responsible for request/response processing
- **services/**: Business logic layer, contains the core application logic
- **models/**: Data structures and database models
- **middleware/**: HTTP middleware for cross-cutting concerns
- **utils/**: Shared utility functions

### Environment Variables

The application supports the following environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | Database host | localhost |
| `DB_PORT` | Database port | 5432 |
| `DB_USER` | Database username | - |
| `DB_PASSWORD` | Database password | - |
| `DB_NAME` | Database name | - |
| `DATABASE_URL` | Complete database URL for migrations | - |
| `SERVER_PORT` | Server port | 8080 |
| `SERVER_MODE` | Gin mode (debug/release) | debug |
| `JWT_SECRET` | JWT signing secret | - |

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📞 Support

For support, email your-email@example.com or create an issue in the GitHub repository.

---

**Happy Quizzing! 🎯**