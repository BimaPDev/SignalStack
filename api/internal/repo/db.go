package repo

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// func Open(databaseURL string) (*sql.DB, error)
// - open postgres connection using lib/pq driver
// - ping to verify connectivity
// - configure pool: MaxOpenConns, MaxIdleConns, ConnMaxLifetime
// - return *sql.DB or error

func Open(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Error: Could not establish a connection with the database", err)
	}
	fmt.Println("Connected to the database!")
	return db, nil
}