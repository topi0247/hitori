-- name: InsertCard :one
INSERT INTO cards (theme_id, profile_id, guest_name, card_number, word, expires_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, uuid;

-- name: GetCardByID :one
SELECT * FROM cards WHERE id = $1;

-- name: CountCardsByThemeID :one
SELECT COUNT(*) FROM cards WHERE theme_id = $1;

-- name: GetAvailableCardNumber :one
SELECT n::int FROM generate_series($1::int, $2::int) AS n
WHERE n NOT IN (SELECT card_number FROM cards WHERE theme_id = $3)
ORDER BY random()
LIMIT 1;

-- name: GetGameCards :many
SELECT * FROM cards
WHERE theme_id = $1 AND is_confirmed = true
ORDER BY random()
LIMIT $2::int;

-- name: ConfirmCard :exec
UPDATE cards SET word = $1, is_confirmed = true, expires_at = NULL, updated_at = now()
WHERE id = $2;

-- name: DeleteCard :exec
DELETE FROM cards WHERE id = $1;

-- name: GetCardsByUUIDs :many
SELECT * FROM cards WHERE uuid::text = ANY($1::text[]);

-- name: AddMatchPoints :exec
UPDATE cards SET match_points = match_points + $1::int, updated_at = now()
WHERE uuid::text = $2;
