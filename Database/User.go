package database

import models "go_server/Models"

func CreateUserTable() error {
	query := `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL
	)`
	_, err := instance.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func InsertUser(email, name,password string) (int, error) {
	query := `INSERT INTO users (email, password , name) VALUES ($1, $2, $3) RETURNING id`
	var pk int
	err := instance.db.QueryRow(query, email, password,name).Scan(&pk)
	if err != nil {
		return 0, err
	}
	return pk, nil
}

func GetAllUsers() ([]models.User, error) {
	query := `SELECT id, email, password FROM users`
	rows, err := instance.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func GetUserByEmail(email string) (models.User, error) {
	query := `SELECT id, email, password FROM users WHERE email = $1`
	var user models.User
	err := instance.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func GetUserById(id string) (models.User , error) {
	query := `SELECT id, email, password FROM users WHERE id = $1`
	var user models.User
	err := instance.db.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.Password)
	if err!= nil {
		return models.User{}, err
	}
	return user, nil
}

func CheckIfUserExists(email string) bool {
	query := `SELECT id FROM users WHERE email = $1`
	var id int
	err := instance.db.QueryRow(query, email).Scan(&id)
	return err == nil
}


