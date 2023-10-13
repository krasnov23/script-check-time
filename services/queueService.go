package services

import (
	"log"
	"script-check-time/config"
	"time"
)

func AddNewQueue(name string, queryCreated time.Time) {
	db := config.GetDB()

	sqlStatement := `INSERT INTO queues (query, date) VALUES ($1, $2)`

	if db == nil {
		log.Printf("DB connection is not established")
		return
	}

	result, err := db.Exec(sqlStatement, name, queryCreated)
	if err != nil {
		log.Printf("Query Error: %s", err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting affected rows: %s", err)
		return
	}

	log.Printf("Affected rows: %d", rowsAffected)

}

func GetAllQueues(timeStamp time.Time) ([]string, error) {
	db := config.GetDB()

	if db == nil {
		log.Printf("DB connection is not established")
		return nil, nil
	}

	rows, err := db.Query("SELECT query FROM queues WHERE date <= $1", timeStamp)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var queries []string
	for rows.Next() {
		var query string
		if err := rows.Scan(&query); err != nil {
			log.Fatal(err)
		}
		queries = append(queries, query)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return queries, nil
}
