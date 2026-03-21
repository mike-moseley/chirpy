-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
	gen_random_uuid(),
	NOW(),
	NOW(),
	$1,
	$2
	)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetHashedPass :one
SELECT hashed_password FROM users
WHERE users.email = $1;

-- name: GetUser :one
SELECT id, created_at, updated_at, email FROM users
WHERE users.email = $1;

-- name: GetUserByID :one
SELECT id, created_at, updated_at, email FROM users
WHERE users.id = $1;

-- name: GetUserFromRefreshToken :one
SELECT id, users.created_at, users.updated_at, email
FROM users
INNER JOIN refresh_tokens ON users.id = refresh_tokens.user_id
WHERE refresh_tokens.token = $1;

-- name: UpdateUserEmailPassword :exec
UPDATE users
SET email = $2, hashed_password = $3, updated_at = NOW()
WHERE id = $1;
