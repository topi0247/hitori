-- name: GetAllThemes :many
SELECT id, title FROM themes ORDER BY id;

-- name: GetThemeByID :one
SELECT id, title FROM themes WHERE id = $1;
