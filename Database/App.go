package database

import models "go_server/Models"

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
