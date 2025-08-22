# Quizer - Quiz Application Backend

A comprehensive quiz application backend built with Go, Gin framework, and PostgreSQL. The application provides a REST API for managing quiz topics, conducting quiz sessions, and tracking user progress.

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

### Code Structure Guidelines

- **handlers/**: HTTP request handlers, responsible for request/response processing
- **services/**: Business logic layer, contains the core application logic
- **models/**: Data structures and database models
- **middleware/**: HTTP middleware for cross-cutting concerns
- **utils/**: Shared utility functions

---

**Happy Quizzing! 🎯**
