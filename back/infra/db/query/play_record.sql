-- name: InsertPlayRecord :exec
INSERT INTO play_records (theme_id, profile_id, card_amount, correct_rate)
VALUES ($1, $2, $3, $4);
