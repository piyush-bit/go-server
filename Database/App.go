package database

import (
	"errors"
	"fmt"
	models "go_server/Models"
	"strings"
)

func CreateAppTable() error {
	query := `CREATE TABLE IF NOT EXISTS apps (
		id SERIAL PRIMARY KEY,
		app_name VARCHAR(255) NOT NULL,
		callback_url TEXT NOT NULL,
		user_id INT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	)`
	_, err := instance.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func InsertApp(name, callbackUrl string, userId int) (int, error) {
	querry := `INSERT INTO apps (app_name , callback_url , user_id) VALUES ($1 , $2 , $3) RETURNING id`
	var pk int
	err := instance.db.QueryRow(querry, name, callbackUrl).Scan(&pk)
	if err != nil {
		return 0, err
	}
	return pk, nil
}

func GetAllAppsOfUser(userId int) ([]models.App, error) {
	query := `SELECT id, app_name, callback_url FROM apps WHERE user_id = $1`

	rows, err := instance.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var apps []models.App
	for rows.Next() {
		var app models.App
		err := rows.Scan(&app.ID, &app.Name, &app.CallbackUrl)
		if err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}
	return apps, nil
}

func GetAppById(appId int) (models.App, error) {
	query := `SELECT id, app_name, callback_url FROM apps WHERE id = $1`
	var app models.App
	err := instance.db.QueryRow(query, appId).Scan(&app.ID, &app.Name, &app.CallbackUrl)
	if err != nil {
		return models.App{}, err
	}
	return app, nil
}
func UpdateApp(appId, userId int, name, callbackUrl string) error {
	if name == "" && callbackUrl == "" {
		return errors.New("name or callback URL are required")
	}

	var query string
	var args []interface{}

	// Build dynamic query based on provided fields
	setParts := make([]string, 0, 2)
	if name != "" {
		setParts = append(setParts, "app_name = $1")
		args = append(args, name)
	}
	if callbackUrl != "" {
		setParts = append(setParts, fmt.Sprintf("callback_url = $%d", len(args)+1))
		args = append(args, callbackUrl)
	}

	// Build the complete query
	query = fmt.Sprintf(
		"UPDATE apps SET %s WHERE id = $%d AND user_id = $%d",
		strings.Join(setParts, ", "),
		len(args)+1,
		len(args)+2,
	)

	// Add WHERE clause parameters
	args = append(args, appId, userId)

	_, err := instance.db.Exec(query, args...)
	return err
}

func DeleteApp(appId, userId int) error {
	query := `DELETE FROM apps WHERE id = $1 AND user_id = $2`
	_, err := instance.db.Exec(query, appId, userId)
	return err
}