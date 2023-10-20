-- name: CreateFeedFollow :execresult
INSERT INTO feed_follows (id, user_id, feed_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?);

-- name: GetUserFeedFollow :many
SELECT `feed_follows`.* FROM feeds
JOIN feed_follows ON feeds.id = feed_follows.feed_id 
WHERE feed_follows.user_id = ?;

-- name: DeleteUserFeedFollow :execresult
DELETE FROM feed_follows WHERE id = ?;


-- name: GetFeedFollowById :one
SELECT * FROM feed_follows WHERE id = ?;
