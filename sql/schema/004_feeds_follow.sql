-- +goose Up

CREATE TABLE IF NOT EXISTS feed_follows (
	id VARCHAR(40) NOT NULL,
	user_id VARCHAR(40) NOT NULL,
	feed_id VARCHAR(40) NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP,
	PRIMARY KEY (user_id, feed_id),
	FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
	FOREIGN KEY (feed_id) REFERENCES feeds (id) ON DELETE CASCADE
);



-- +goose Down
DROP TABLE IF EXISTS feed_follows;
