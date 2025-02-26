package database

import models "go_server/Models"

func CreateTokenTable() error {
	query := `CREATE TABLE IF NOT EXISTS tokens (
		id SERIAL PRIMARY KEY,
		app_id INT NOT NULL,
		token TEXT NOT NULL,
		refresh_token TEXT ,
		FOREIGN KEY (app_id) REFERENCES apps(id) ON DELETE CASCADE
	)`
	_, err := instance.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func InsertToken(appId int, token , refreshToken string) (int, error) {
	query := `INSERT INTO tokens (app_id, token,refresh_token) VALUES ($1, $2 , $3) RETURNING id`
	var pk int
	err := instance.db.QueryRow(query, appId, token , refreshToken).Scan(&pk)
	if err != nil {
		return 0, err
	}
	return pk, nil
}

func GetTokenById(tokenId int) (models.Token, error) {
	query := `SELECT id, app_id, token , refresh_token FROM tokens WHERE id = $1`
	var token models.Token
	err := instance.db.QueryRow(query, tokenId).Scan(&token.ID, &token.AppId, &token.Token, &token.RefreshToken)
	if err != nil {
		return models.Token{}, err
	}
	return token, nil
}

func DeleteToken(tokenId int) error {
	query := `DELETE FROM tokens WHERE id = $1`
	_, err := instance.db.Exec(query, tokenId)
	if err != nil {
		return err
	}
	return nil
}
