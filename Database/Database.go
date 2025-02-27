package database

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	"github.com/pressly/goose/v3"
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
}

func GetInstance() *Database {
	once.Do(func() {
		db := connectDB()

		// Run migrations
		err := runMigrations(db)
		if err!= nil {
			panic(fmt.Sprintf("Error running migrations: %v", err))
		}
		instance = &Database{
			db: db,
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

func runMigrations(db *sql.DB) error {
    // Set dialect to postgres
    if err := goose.SetDialect("postgres"); err != nil {
        return err
    }
    
    // Run migrations from the migrations directory
    if err := goose.Up(db, "Database/Migrations"); err != nil {
        return err
    }
    
    fmt.Println("Database migrations completed successfully")
    return nil
}