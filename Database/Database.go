package database

import (
	"database/sql"
	"fmt"
	"sync"
	_ "github.com/lib/pq"
)

var (
	instance *Database
	once     sync.Once
)

const connectionString = "postgres://postgres:PASS@localhost:5432/gopgtest?sslmode=disable"

type Database struct {
	db *sql.DB
}

func init() {
	GetInstance()
	err := CreateUserTable()
	if err != nil {
		panic(err)
	}
	err = CreateAppTable()
	if err != nil {
		panic(err)
	}
	err = CreateTokenTable()
	if err != nil {
		panic(err)
	}
	err = CreateSessionTable()
	if err!= nil {
		panic(err)
	}
}

func GetInstance() *Database {
	once.Do(func() {
		instance = &Database{
			db: connectDB(),
		}
	})
	return instance
}

func connectDB() *sql.DB {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to the database")
	return db
}

