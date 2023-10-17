package services

import (
	"fmt"
	"log"
	"script-check-time/config"
	"time"
)

func AddNewQueue(eventId int, queryCreated time.Time) {

	db := config.GetDB()

	sqlStatement := `INSERT INTO queues (event_id, date) VALUES ($1, $2)`

	if db == nil {
		log.Printf("DB connection is not established")
		return
	}

	result, err := db.Exec(sqlStatement, eventId, queryCreated)
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

func AddNewQueueByEventData(id int, period string) {

	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1)
	year, month, day := tomorrow.Date()

	layout := "2006-01-02 15:04:05"

	// Concatenate current date and incoming time for parsing
	dateTimeStr := fmt.Sprintf("%d-%02d-%02d %s", year, month, day, period)

	// Parsing the string to time
	myTime, err := time.Parse(layout, dateTimeStr)

	// Handle error
	if err != nil {
		fmt.Println("Parsing error", err.Error())
	} else {
		fmt.Println("Parsed time in current date: ", myTime)
	}

	AddNewQueue(id, myTime)
}

func GetAllQueuesByDate(timeStamp time.Time) ([]string, error) {

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

func GetExpiredQueues() []int {

	db := config.GetDB()

	if db == nil {
		log.Printf("DB connection is not established")
	}

	timeStamp := time.Now()

	rows, err := db.Query("SELECT event_id FROM queues WHERE date <= $1", timeStamp)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var queries []int
	for rows.Next() {
		var query int
		if err := rows.Scan(&query); err != nil {
			log.Fatal(err)
		}
		queries = append(queries, query)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(queries)
	return queries

}

func DeleteExpireQueues() {

	db := config.GetDB()

	if db == nil {
		log.Printf("DB connection is not established")
	}

	timeStamp := time.Now()

	rows, err := db.Query("DELETE FROM queues WHERE date <= $1", timeStamp)
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

}

func AddQueue() {

	now := time.Now()
	times := []time.Time{
		now, now.Add(time.Second * 40),
		now.Add(time.Second * 70),
		now.Add(time.Second * 80),
		now.Add(time.Second * 90),
		now.Add(time.Second * 125),
	}

	for i := range times {
		AddNewQueue(i, times[i])
	}

}
