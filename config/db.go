package config

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

var DB *sql.DB

func Connect() {
	connStr := "postgres://postgres:postgres@localhost:5434/postgres?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	query := "CREATE TABLE IF NOT EXISTS queues (id SERIAL PRIMARY KEY,query VARCHAR,date TIMESTAMP)"
	_, err = db.Exec(query)

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Table created successfully!")
	}

	DB = db
}

func GetDB() *sql.DB {
	return DB
}
