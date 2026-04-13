# Bytecode API

Bytecode API is a robust, production-ready, and modular Go web API boilerplate. It follows clean architecture principles and comes pre-configured with essential features such as Role-Based Access Control (RBAC), JWT Authentication, database migrations, and background worker processing.

## 🚀 Features

- **Clean & Modular Architecture**: Structured by features (auth, user, role, permission).
- **Authentication & Authorization**: Secure JWT-based authentication combined with an RBAC middleware to handle user permissions precisely.
- **Robust Routing capability**: Utilizing [Fiber](https://gofiber.io/) as our performant HTTP engine.
- **Database Operations**: Integrated with PostgreSQL using [GORM](https://gorm.io/) and `golang-migrate` for version control of schema migrations.
- **Background Jobs**: Built-in generic worker pool & job scheduling mechanism.
- **Structured Logging**: Leveraging Uber's [Zap logger](https://github.com/uber-go/zap) for fast, structured, and leveled logging.

## 🛠️ Tech Stack

- **Language**: [Go (Golang)](https://go.dev/)
- **Framework**: [Fiber v2](https://github.com/gofiber/fiber)
- **Database**: PostgreSQL
- **ORM**: [GORM](https://gorm.io/)
- **Migrations**: [golang-migrate](https://github.com/golang-migrate/migrate)
- **Logging**: [Uber Zap](https://github.com/uber-go/zap)

## 📋 Prerequisites

Before you begin, ensure you have met the following requirements:
- **Go** >= 1.21 installed
- **PostgreSQL** installed and running
- **Make** installed (for executing Makefile commands)

## ⚙️ Installation & Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/ramdhanrizkij/bytecode-api.git
   cd bytecode-api
   ```

2. **Setup environment variables:**
   Copy the example environment file and adjust the parameters to match your local setup.
   ```bash
   cp .env.example .env
   ```
   *(Make sure to update the database connection strings, JWT secret, etc., inside `.env`)*

3. **Install dependencies:**
   Download all required Go modules.
   ```bash
   go mod tidy
   ```

4. **Setup Database Migrations:**
   Install the migration CLI tool. A convenience command is provided in the Makefile:
   ```bash
   make migrate-setup
   ```
   *Note: If the `migrate` command is not recognized after installation, ensure that `~/go/bin` is added to your system's `PATH`.*

5. **Run the Migrations:**
   Create your tables to prepare the database schema.
   ```bash
   make migrate-up
   ```

## 💻 Usage

The API includes several pre-defined `make` scripts to streamline the development workflow:

**Run the API Server:**
```bash
make run
```
By default, the server will start on `http://127.0.0.1:8080`.

**Run the Worker (Background jobs):**
```bash
make run-worker
```

**Testing commands:**
```bash
make test             # Run all tests
make test-unit        # Run unit tests only
make test-integration # Run integration tests
```

**Database commands:**
```bash
make migrate-up      # Apply pending database migrations
make migrate-down    # Revert the latest database migration
make migrate-refresh # Drop all and re-apply all migrations completely
make migrate-create name=create_users_table # Generate a new blank migration file
```

## 📁 Project Structure

```text
├── cmd/
│   ├── api/          # Main entry point for the REST API
│   └── worker/       # Main entry point for the background worker
├── internal/
│   ├── core/         # Core system utilities (config, middleware, server, job pool)
│   ├── features/     # Feature modules (Auth, User, Role, Permission)
│   └── shared/       # Shared generic utilities (Response format, constants, Validation)
├── migrations/       # SQL migration scripts
├── test/             # Integration tests
├── .env.example      # Environment variables template
├── Makefile          # Useful execution commands
└── go.mod            # Go module dependencies
```

## 📄 License
This project is licensed under the MIT License.
