-- name: CreateFeed :execresult
INSERT INTO feeds (id, name, url, user_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?);

-- name: GetAllFeeds :many
SELECT * FROM feeds;

-- name: GetFeedById :one
SELECT * FROM feeds WHERE id = ?;
