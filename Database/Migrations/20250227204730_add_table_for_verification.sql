-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS forget_password (
	email VARCHAR(255) NOT NULL UNIQUE PRIMARY KEY,
	token TEXT NOT NULL,
	expired_at TIMESTAMP NOT NULL
	);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS forget_password;
-- +goose StatementEnd
