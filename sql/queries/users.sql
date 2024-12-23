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

-- name: CreateChirp :one
INSERT INTO chirps(id, created_at, updated_at, body, user_id)
VALUES(
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetChirps :many
SELECT * FROM chirps
ORDER BY created_at;

-- name: GetChirpsDesc :many
SELECT * FROM chirps
ORDER BY created_at DESC;

-- name: GetChirp :one
SELECT * FROM chirps
WHERE chirps.id = $1;

-- name: GetChirpsByUserID :many
SELECT * FROM chirps
WHERE chirps.user_id = $1
ORDER BY created_at;

-- name: GetChirpsByUserIDDesc :many
SELECT * FROM chirps
WHERE chirps.user_id = $1
ORDER BY created_at DESC;


-- name: GetUser :one
SELECT * FROM users
WHERE users.email = $1;

-- name: CreateToken :one
INSERT INTO refresh_tokens(token, created_at,
updated_at, user_id, expires_at, revoked_at)
VALUES(
    $1,
    NOW(),
    NOW(),
    $2,
    NOW() + INTERVAL '60day',
    NULL
)
RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token = $1;

-- name: RevokeToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1;

-- name: UpdateUserEmailPw :exec
UPDATE users
SET hashed_password = $1, email = $2
WHERE id = $3;

-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id = $1;

-- name: UpgradeChripy :exec
UPDATE users
SET is_chirpy_red = TRUE
WHERE id = $1;
