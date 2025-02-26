package database

import models "go_server/Models"


func CreateSessionTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS sessions (
		user_id INT NOT NULL,
		app_id INT NOT NULL,
		refresh_token TEXT,
		PRIMARY KEY (user_id, app_id),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (app_id) REFERENCES apps(id) ON DELETE CASCADE
	);
	`
	_, err := instance.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func InsertSession(userId, appId int, refreshToken string) (int, error) {
	query := `INSERT INTO sessions (user_id, app_id, refresh_token) VALUES ($1, $2, $3) RETURNING id`
	var pk int
	err := instance.db.QueryRow(query, userId, appId, refreshToken).Scan(&pk)
	if err!= nil {
		return 0, err
	}
	return pk, nil
}

func UpdateRefreshToken(userId , appId int, refreshToken string) error {
	query := `UPDATE sessions SET refresh_token = $1 WHERE user_id = $2 AND app_id = $3`
	_, err := instance.db.Exec(query, refreshToken, userId, appId)
	if err!= nil {
		return err
	}
	return nil
}

func DeleteSession(userId, appId int) error {
	query := `DELETE FROM sessions WHERE user_id = $1 AND app_id = $2`
	_, err := instance.db.Exec(query, userId, appId)
	if err!= nil {
		return err
	}
	return nil
}

func InsertOrUpdateSession(userId, appId int, refreshToken string) error {
	query := `
		INSERT INTO sessions (user_id, app_id, refresh_token)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, app_id) DO UPDATE
		SET refresh_token = $3
	`
	_, err := instance.db.Exec(query, userId, appId, refreshToken)
	if err!= nil {
		return err
	}
	return nil
}

func GetRefreshToken(userId, appId int) (string, models.User, error) {
	query := `
		SELECT users.name , users.email , sessions.refresh_token 
		FROM sessions 
		INNER JOIN users ON sessions.user_id = users.id 
		WHERE sessions.user_id = $1 AND sessions.app_id = $2
	`
	var name , email , refreshToken string
	err := instance.db.QueryRow(query, userId, appId).Scan(&name,&email,&refreshToken)
	if err != nil {
		return "", models.User{}, err
	}
	return refreshToken, models.User{Name: name,Email: email}, nil
}
