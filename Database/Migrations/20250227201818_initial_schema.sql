-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL
	);

CREATE TABLE IF NOT EXISTS apps (
		id SERIAL PRIMARY KEY,
		app_name VARCHAR(255) NOT NULL,
		callback_url TEXT NOT NULL,
		user_id INT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

CREATE TABLE IF NOT EXISTS sessions (
		user_id INT NOT NULL,
		app_id INT NOT NULL,
		refresh_token TEXT,
		PRIMARY KEY (user_id, app_id),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (app_id) REFERENCES apps(id) ON DELETE CASCADE
	);

CREATE TABLE IF NOT EXISTS tokens (
		id SERIAL PRIMARY KEY,
		app_id INT NOT NULL,
		token TEXT NOT NULL,
		refresh_token TEXT,
		FOREIGN KEY (app_id) REFERENCES apps(id) ON DELETE CASCADE
	);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS tokens;
DROP TABLE IF EXISTS apps;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
