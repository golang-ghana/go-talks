# API Documentation

This project provides REST APIs for both SQLC and GORM examples, organized by functionality.

## Setup

1. **Install dependencies:**
   ```bash
   go mod tidy
   ```

2. **Create `.env` file** (copy from `.env.example`):
   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=user
   DB_PASSWORD=pass
   DB_NAME=dbname
   DB_SSLMODE=disable
   SERVER_PORT=8080
   ```

3. **Generate SQLC code** (required for SQLC endpoints):
   ```bash
   sqlc generate
   ```
   Then uncomment the SQLC handler implementations in the code.

4. **Run the server:**
   ```bash
   go run api/main.go
   ```

## API Endpoints

### Health Check
- `GET /health` - Server health check

### GORM API (`/api/gorm`)

#### User Operations
- `POST /api/gorm/users` - Create a new user
- `GET /api/gorm/users/:id` - Get user by ID
- `GET /api/gorm/users/email/:email` - Get user by email
- `GET /api/gorm/users` - List users (paginated, query params: `limit`, `offset`)
- `GET /api/gorm/users/active` - List active users
- `PUT /api/gorm/users/:id` - Update user
- `PATCH /api/gorm/users/:id/status` - Update user active status
- `DELETE /api/gorm/users/:id` - Delete user
- `GET /api/gorm/users/search?q=name` - Search users by name
- `GET /api/gorm/users/age-range?min_age=25&max_age=35` - Get users by age range
- `GET /api/gorm/users/count` - Count total users

#### Advanced GORM Operations
- `GET /api/gorm/users/complex` - Complex query (active users with age >= 25 or username like 'admin%')
- `GET /api/gorm/users/partial` - Select specific fields only
- `GET /api/gorm/users/raw?min_age=20&limit=5` - Raw SQL query

#### Transactions
- `POST /api/gorm/transactions/users` - Create multiple users in a transaction

### SQLC API (`/api/sqlc`)

**Note:** SQLC endpoints require running `sqlc generate` first and uncommenting the handler implementations.

#### User Operations
- `POST /api/sqlc/users` - Create a new user
- `GET /api/sqlc/users/:id` - Get user by ID
- `GET /api/sqlc/users/email/:email` - Get user by email
- `GET /api/sqlc/users` - List users (paginated, query params: `limit`, `offset`)
- `GET /api/sqlc/users/active` - List active users
- `PUT /api/sqlc/users/:id` - Update user
- `PATCH /api/sqlc/users/:id/status` - Update user active status
- `DELETE /api/sqlc/users/:id` - Delete user
- `GET /api/sqlc/users/search?q=name` - Search users by name
- `GET /api/sqlc/users/age-range?min_age=25&max_age=35` - Get users by age range
- `GET /api/sqlc/users/count` - Count total users

#### Transactions
- `POST /api/sqlc/transactions/users` - Create multiple users in a transaction

## Request/Response Examples

### Create User
```bash
POST /api/gorm/users
Content-Type: application/json

{
  "email": "john@example.com",
  "username": "johndoe",
  "full_name": "John Doe",
  "age": 30,
  "is_active": true
}
```

### Update User
```bash
PUT /api/gorm/users/1
Content-Type: application/json

{
  "email": "john.doe@example.com",
  "username": "johndoe",
  "full_name": "John Michael Doe",
  "age": 31
}
```

### Update User Status
```bash
PATCH /api/gorm/users/1/status
Content-Type: application/json

{
  "is_active": false
}
```

### List Users (Paginated)
```bash
GET /api/gorm/users?limit=10&offset=0
```

### Search Users
```bash
GET /api/gorm/users/search?q=John
```

### Age Range Filter
```bash
GET /api/gorm/users/age-range?min_age=25&max_age=35
```

### Transaction (Multiple Users)
```bash
POST /api/gorm/transactions/users
Content-Type: application/json

{
  "users": [
    {
      "email": "alice@example.com",
      "username": "alice",
      "full_name": "Alice Smith",
      "age": 25,
      "is_active": true
    },
    {
      "email": "bob@example.com",
      "username": "bob",
      "full_name": "Bob Johnson",
      "age": 28,
      "is_active": true
    }
  ]
}
```

## GORM Query Logging

GORM queries are automatically logged to the console. You'll see SQL queries being executed when you make API requests to GORM endpoints.

## Notes

- GORM endpoints are fully functional and include query logging
- SQLC endpoints require code generation first (`sqlc generate`)
- After generating SQLC code, uncomment the implementations in the handler files
- All endpoints return JSON responses
- Error responses follow the format: `{"error": "error message"}`
