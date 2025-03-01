package database

import (
	models "go_server/Models"
	"time"
)

func CreateForgetPasswordTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS forget_password (
	email VARCHAR(255) NOT NULL UNIQUE PRIMARY KEY,
	token TEXT NOT NULL,
	expired_at TIMESTAMP NOT NULL
	)
	`
	_, err := instance.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func InsertForgetPassword(email string, token string, expired_at time.Time) error {
	query := `INSERT INTO forget_password (email, token, expired_at) VALUES ($1, $2, $3) ON CONFLICT (email) DO UPDATE SET token = $2, expired_at = $3`
	_, err := instance.db.Exec(query, email, token, expired_at)
	if err != nil {
		return err
	}
	return nil
}

func GetTokenByEmail(email string) (models.ForgetPassword, error) {
	query := `SELECT email, token, expired_at FROM forget_password WHERE email = $1`
	var token models.ForgetPassword
	err := instance.db.QueryRow(query, email).Scan(&token.Email, &token.Token, &token.ExpiredAt)
	if err!= nil {
		return models.ForgetPassword{}, err
	}
	return token, nil
}
