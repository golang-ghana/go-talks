-- queries.sql (for sqlc)

-- name: CreateUser :one
INSERT INTO users (email, username, full_name, age, is_active)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListActiveUsers :many
SELECT * FROM users
WHERE is_active = true
ORDER BY created_at DESC;

-- name: UpdateUser :one
UPDATE users
SET 
    email = $2,
    username = $3,
    full_name = $4,
    age = $5,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateUserActiveStatus :exec
UPDATE users
SET is_active = $2, updated_at = NOW()
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: SearchUsersByName :many
SELECT * FROM users
WHERE full_name ILIKE '%' || $1 || '%'
ORDER BY full_name;

-- name: GetUsersByAgeRange :many
SELECT * FROM users
WHERE age BETWEEN $1 AND $2
ORDER BY age;
