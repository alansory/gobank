-- name: CreateUser :one
INSERT INTO users (
  username,
  email,
  hashed_password
) values (
  $1, $2, $3
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 limit 1;

-- name: UpdateUser :one
UPDATE users
SET 
  hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password),
  password_changed_at = COALESCE(sqlc.narg(password_changed_at), password_changed_at),
  username = COALESCE(sqlc.narg(username), username)
WHERE 
  id = sqlc.arg(id)
RETURNING *;