-- name: CreateUser :one
INSERT INTO users (
  fullname,
  username,
  email,
  hashed_password
) values (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE email = $1 limit 1;

-- name: UpdateUser :one
UPDATE users
SET 
  hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password),
  password_changed_at = COALESCE(sqlc.narg(password_changed_at), password_changed_at),
  fullname = COALESCE(sqlc.narg(fullname), fullname)
  username = COALESCE(sqlc.narg(username), username)
WHERE 
  id = sqlc.arg(id)
RETURNING *;