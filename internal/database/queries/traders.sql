-- name: GetTrader :one
SELECT * FROM traders
WHERE id = ? LIMIT 1;

-- name: ListTraders :many
SELECT * FROM traders
ORDER BY id;

-- name: CreateTrader :one
INSERT INTO traders (
  balance
) VALUES (
  ?
)
RETURNING *;

-- name: UpdateTrader :exec
UPDATE traders
set balance = ?
WHERE id = ?;

-- name: DeleteTrader :exec
DELETE FROM traders
WHERE id = ?;
