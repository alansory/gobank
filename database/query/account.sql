-- name: createAccount :one
INSERT INTO accounts (
  user_id,
  balance,
  currency
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 limit 1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1 limit 1
FOR NO KEY UPDATE;

-- name: listAccounts :many
SELECT * FROM accounts
WHERE user_id = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdateAccount :one
UPDATE accounts 
SET balance = $2, currency = $3
WHERE id = $1
RETURNING *;