-- name: InsertProfile :one
INSERT INTO profiles (auth_user_id, user_name) VALUES ($1, $2)
RETURNING id, auth_user_id, user_name;

-- name: GetProfileByAuthUserID :one
SELECT id, auth_user_id, user_name FROM profiles WHERE auth_user_id = $1;

-- name: UpdateProfileUserName :exec
UPDATE profiles SET user_name = $1, updated_at = now() WHERE auth_user_id = $2;

-- name: DeleteProfile :exec
DELETE FROM profiles WHERE auth_user_id = $1;
