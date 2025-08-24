# Auto Light Pi Backend

This is the backend service for the Auto Light Pi project, built with Go.


## Project Structure
The backend follows a **clean layered architecture**:
```
backend/
  main.go                # Application entry point
  .env                   # Environment variables
  schema.sql             # Database schema
  internal/
   config/                # Configuration and initialization
   controllers/           # HTTP controllers (handlers)
   middleware/            # Middleware (e.g., auth)
   models/                # Data models (entities)
   repositories/          # Data access layer
   routes/                # API route definitions
   services/              # Business logic layer
   bootstrap/             # Application bootstrapping
```

In particular, the flow of a request through the backend is as follows:
```
Request
  │
  ▼
[Routes] → [Middleware] → [Controllers] → [Services] → [Repositories] → [Database]
```

- **Routes**: Define API endpoints and attach middleware/controllers.
- **Middleware**: Handle cross-cutting concerns (e.g., authentication).
- **Controllers**: Handle HTTP requests, parse input, and return responses.
- **Services**: Contain business logic and orchestrate repository calls.
- **Repositories**: Abstract database access and queries.
- **Models**: Define data structures (entities).
- **Config**: Manage configuration and environment setup.
- **Bootstrap**: Application bootstrapping and initialization.

## Environment Variables

Create a `.env` file in the backend root with the following variables:

```
PORT=8080
# Application
APPLICATION_NAME=auto-light-pi

# Database
DB_DSN="host=localhost user=yourusername dbname=autolightpi port=5432 sslmode=disable"

# JWT Secret Key
JWT_SECRET="yoursecretkey"
```

---

## Running the Server

1. **Install dependencies:**
   ```sh
   go mod tidy
   ```
2. **Run the server:**
   ```sh
   go run main.go
   ```

---

## Testing
Unit and integration tests will be added in the future.

## License

This project is licensed under the MIT License. See the [LICENSE](../LICENSE) file for details.
