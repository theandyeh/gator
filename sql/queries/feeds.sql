-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, url, name, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: DeleteFeeds :exec
DELETE FROM feeds;

-- name: GetFeeds :many
SELECT feeds.id, feeds.url, feeds.name, feeds.user_id, users.name AS username
FROM feeds
LEFT JOIN users ON feeds.user_id = users.id;
