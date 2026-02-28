# Setup Instructions

## Prerequisites

- Go 1.25.5 or later
- PostgreSQL database
- SQLC (for SQLC endpoints) - install with: `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`

## Quick Start

1. **Install dependencies:**
   ```bash
   go mod tidy
   ```

2. **Create `.env` file** in the project root:
   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=your_db_user
   DB_PASSWORD=your_db_password
   DB_NAME=your_db_name
   DB_SSLMODE=disable
   SERVER_PORT=8080
   ```

3. **Set up your database:**
   - Create a PostgreSQL database
   - Run the schema: `psql -d your_db_name -f schema.sql`

4. **Generate SQLC code** (required for SQLC endpoints):
   ```bash
   sqlc generate
   ```
   This will create the `sqlc/` directory with generated code.

5. **Enable SQLC handlers:**
   - After running `sqlc generate`, edit `api/main.go`:
     - Add import: `"gorm-sqlc/sqlc/db"`
     - Uncomment the SQLC handler initialization lines
   - Edit `handlers/sqlc/user_handler.go` and `handlers/sqlc/transaction_handler.go`:
     - Uncomment the `queries` field in structs
     - Uncomment all the TODO sections with actual implementations

6. **Run the API server:**
   ```bash
   go run api/main.go
   ```

7. **Test the API:**
   ```bash
   # Health check
   curl http://localhost:8080/health

   # Create a user (GORM)
   curl -X POST http://localhost:8080/api/gorm/users \
     -H "Content-Type: application/json" \
     -d '{
       "email": "john@example.com",
       "username": "johndoe",
       "full_name": "John Doe",
       "age": 30,
       "is_active": true
     }'
   ```

## GORM Query Logging

GORM queries are automatically logged to the console. You'll see SQL queries when making requests to GORM endpoints.

## API Documentation

See `API_README.md` for complete API documentation.
