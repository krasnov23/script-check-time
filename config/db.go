package config

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

var DB *sql.DB

func Connect() {
	connStr := "postgres://postgres:postgres@localhost:5435/postgres?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	DB = db
	CreateQueueTable()
	CreateEventTable()

}

func CreateQueueTable() {

	db := GetDB()

	query := "CREATE TABLE IF NOT EXISTS queues (id SERIAL PRIMARY KEY,event_id BIGINT,date TIMESTAMP)"

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Table Queue created successfully!")
	}
}

func CreateEventTable() {
	db := GetDB()
	query := "CREATE TABLE IF NOT EXISTS event_reference (id SERIAL PRIMARY KEY,name VARCHAR(255),url VARCHAR(255),period VARCHAR(255))"
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Table Event created successfully!")
	}
}

func GetDB() *sql.DB {
	return DB
}
