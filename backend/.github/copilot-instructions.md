# Copilot Instructions for Auto-Light-Pi Backend

## Architecture Overview
This project follows a **Clean Layered Architecture** in Go.
- **Project Structure**:
  - `internal/`: Contains all application code, organized by domain/feature (e.g., `auth`, `user`, `refresh_token`).
  - `bootstrap/`: Application initialization and dependency injection (`server.go`, `database.go`).
  - `routes/`: API route definitions (`routes.go`).
  - `main.go`: Entry point, loads configuration and starts the server.

- **Layer Responsibilities**:
  1. **Routes** (`internal/routes`): Defines HTTP endpoints and maps them to controllers.
  2. **Controllers** (`internal/*/controller.go`): Handles HTTP requests/responses (Gin), input validation, and calls Services.
  3. **Services** (`internal/*/service.go`): Contains business logic. Defines interfaces for its dependencies (repositories).
  4. **Repositories** (`internal/*/repository.go`): Handles raw SQL database access. `database/sql` is used directly with `lib/pq`.

## Key Conventions & Patterns

### 1. Dependency Injection
- Use **manual dependency injection** in `internal/bootstrap/server.go`.
- Avoid global state for dependencies.
- Services should define interfaces for the repositories they need.
- **Example**: `auth.NewAuthService(userRepo, refreshTokenRepo)`

### 2. Error Handling
- **Services**: Return domain-specific errors (defined as `var ErrX = errors.New(...)` at the top of `service.go`).
- **Repositories**: Handle `sql.ErrNoRows` appropriately (e.g., return `nil, nil` or a specific error).
- **Controllers**: Check for specific errors using `errors.Is()` and map them to appropriate HTTP status codes.
- **Context**: Always pass `ctx context.Context` through all layers (Controller -> Service -> Repository).

### 3. Database Access
- Use **Raw SQL** via `database/sql`. Do not use ORMs like GORM.
- Use `ExecContext`, `QueryRowContext`, etc., to respect context cancellation/timeouts.
- Use explicit Entity structs (e.g., `UserEntity` in `repository.go`) to map DB columns, then convert to Domain Models (e.g., `User` struct) before returning.
- **Pattern**: `toEntity()` and `toDomain()` helper methods in repository files.

### 4. Logging
- Use `log/slog` for structured logging.
- Configure logging in `main.go` with `lumberjack` for rotation.
- Log significant events and errors in Controllers.

### 5. API & Validation
- Framework: **Gin Web Framework**.
- Use struct tags for validation (e.g., `binding:"required,email"`).
- Define request/response structs inside the controller file (e.g., `registerRequest`, `loginUserResponse`).

## Developer Workflow

### Environment Setup
- Configuration is loaded from **environment variables** (see `main.go`).
- Essential vars: `BACKEND_PORT`, `APPLICATION_NAME`, database credentials.

### File Locations
- **Routes**: `internal/routes/routes.go`
- **Wiring**: `internal/bootstrap/server.go`
- **Feature Code**: `internal/<feature_name>/`

### Testing
- Use standard `testing` package.
- Mocks are located in `mocks/` subdirectories (e.g., `internal/auth/mocks/`).
- Use `mockgen` or manual mocks for service/repository interfaces.

## External Dependencies
- **PostgreSQL**: Primary data store (`lib/pq`).
- **Redis**: Used for token storage (`go-redis/v9`).
- **JWT**: Token generation/validation (`golang-jwt/jwt/v5`).
