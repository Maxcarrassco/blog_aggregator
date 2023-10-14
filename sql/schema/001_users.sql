-- +goose Up

CREATE TABLE users (
id VARCHAR(40) NOT NULL,
name VARCHAR(256) NOT NULL,
created_at TIMESTAMP NOT NULL,
updated_at TIMESTAMP,
PRIMARY KEY(id)
);


-- +goose Down
DROP TABLE users;
