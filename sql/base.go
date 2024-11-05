package sql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db   *sql.DB
	once sync.Once
)

func getDB() *sql.DB {
	once.Do(func() {
		var err error
		username := os.Getenv("SQL_USER")
		password := os.Getenv("SQL_PASSWORD")
		host := os.Getenv("SQL_HOST")
		port := os.Getenv("SQL_PORT")
		database := os.Getenv("SQL_DB")
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, database)
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(time.Hour)

		if err := db.Ping(); err != nil {
			log.Fatalf("Failed to ping database: %v", err)
		}
		fmt.Println("Connected to the database with connection pooling!")
	})

	return db
}
