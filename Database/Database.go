package database

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

var (
	instance *Database
	once     sync.Once
)

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
	if err != nil {
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
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		panic("DATABASE_URL is not set")
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to the database")
	return db
}
