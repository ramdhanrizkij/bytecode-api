# Bytecode API

Bytecode API is a robust, production-ready, and modular Go web API boilerplate. It follows clean architecture principles and comes pre-configured with essential features such as Role-Based Access Control (RBAC), JWT Authentication, database migrations, and background worker processing.

## 🚀 Features

- **Clean & Modular Architecture**: Structured by features (auth, user, role, permission).
- **Authentication & Authorization**: Secure JWT-based authentication combined with an RBAC middleware to handle user permissions precisely.
- **Robust Routing capability**: Utilizing [Fiber](https://gofiber.io/) as our performant HTTP engine.
- **Database Operations**: Integrated with PostgreSQL using [GORM](https://gorm.io/) and [Gormigrate](https://github.com/go-gormigrate/gormigrate) for versioned schema migrations.
- **Background Jobs**: Built-in generic worker pool & job scheduling mechanism.
- **Optional Redis Caching**: Redis can be enabled via environment variables for caching frequently accessed endpoints.
- **Pluggable File Storage**: Choose between local filesystem storage and MinIO (S3-compatible) through environment variables.
- **Structured Logging**: Leveraging Uber's [Zap logger](https://github.com/uber-go/zap) for fast, structured, and leveled logging.
- **Interactive API Documentation**: Swagger UI is generated from Swaggo annotations and exposed at `/swagger/index.html`.

## 🛠️ Tech Stack

- **Language**: [Go (Golang)](https://go.dev/)
- **Framework**: [Fiber v3](https://github.com/gofiber/fiber)
- **Database**: PostgreSQL
- **ORM**: [GORM](https://gorm.io/)
- **Migrations**: [Gormigrate](https://github.com/go-gormigrate/gormigrate)
- **Logging**: [Uber Zap](https://github.com/uber-go/zap)
- **Cache**: Redis (optional)
- **Object Storage**: Local filesystem or MinIO (optional)
- **API Docs**: [Swaggo](https://github.com/swaggo/swag) + [Fiber Contrib Swaggo](https://github.com/gofiber/contrib/tree/main/v3/swaggo)

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

4. **Run the Migrations:**
   Create your tables to prepare the database schema.
   ```bash
   make migrate-up
   ```

5. **Optional Redis cache:**
   Redis is disabled by default. To enable it in the app, set `REDIS_ENABLED=true` in `.env`.
   If you also want Docker Compose to start Redis, set `COMPOSE_PROFILES=redis`.

6. **Optional file storage provider:**
   Local storage is the default with `STORAGE_PROVIDER=local`, storing files under `STORAGE_LOCAL_PATH` and serving them from `STORAGE_BASE_URL`.
   To use MinIO instead, set `STORAGE_PROVIDER=minio` and configure the `MINIO_*` variables in `.env`.
   If you want Docker Compose to start MinIO too, set `COMPOSE_PROFILES=minio` or combine profiles such as `COMPOSE_PROFILES=redis,minio`.

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
make migrate-create name=create_users_table # Generate a registered Go migration stub
```

**Swagger documentation:**
```bash
make swagger-setup # Install the swag CLI
make swagger       # Re-generate docs/, swagger.json, and swagger.yaml
```
Swagger UI is available at `http://127.0.0.1:8080/swagger/index.html` when enabled. It is exposed by default outside production. In production, set `SWAGGER_ENABLED=true` together with `SWAGGER_USERNAME` and `SWAGGER_PASSWORD` to protect it with basic auth; if credentials are missing, the endpoint is not registered.

## 📁 Project Structure

```text
├── cmd/
│   ├── api/          # Main entry point for the REST API
│   └── worker/       # Main entry point for the background worker
├── internal/
│   ├── core/         # Core system utilities (config, middleware, server, job pool)
│   ├── features/     # Feature modules (Auth, User, Role, Permission)
│   └── shared/       # Shared generic utilities (Response format, constants, Validation)
├── migrations/       # Versioned Gormigrate definitions
├── test/             # Integration tests
├── .env.example      # Environment variables template
├── Makefile          # Useful execution commands
└── go.mod            # Go module dependencies
```

## 📄 License
This project is licensed under the MIT License.
