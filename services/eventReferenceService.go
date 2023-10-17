package services

import (
	"fmt"
	"log"
	"reflect"
	"script-check-time/config"
	"strings"
)

type EventReference struct {
	ID     int
	Period string
}

func AddNewEvent(name string, url string, period string) {

	db := config.DB

	sqlStatement := `INSERT INTO event_reference (name, url, period) VALUES ($1,$2,$3)`

	if db == nil {
		log.Printf("DB connection is not established")
		return
	}

	_, err := db.Exec(sqlStatement, name, url, period)
	if err != nil {
		log.Printf("Query Error: %s", err)
		return
	}

}

func AddNewEventWithGetIdRecordBack(name string, url string, insertedPeriod string) (int, string) {

	db := config.DB

	// Here you prepare your SQL query
	stmt, err := db.Prepare("INSERT INTO event_reference (name, url, period) VALUES ($1,$2,$3) RETURNING id, period")
	if err != nil {
		log.Fatal(err)
	}

	var id int

	var period string

	// Execute the query. Scan for the returned 'id'
	err = stmt.QueryRow(name, url, insertedPeriod).Scan(&id, &period)
	if err != nil {
		log.Fatal(err)
	}

	return id, period
}

func FindEventsByIdsAndGetIdAndPeriod(ids []int) []EventReference {

	db := config.DB

	strIDs := make([]interface{}, len(ids))
	for i, id := range ids {
		strIDs[i] = id
	}

	placeholders := strings.Repeat(",?", len(ids)-1)
	fmt.Println(reflect.TypeOf(placeholders))

	query := fmt.Sprintf(`SELECT id, name, url, period FROM event_reference WHERE id IN (?%s)`, placeholders)
	fmt.Println(query)

	rows, err := db.Query(query, strIDs...)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var records []EventReference
	for rows.Next() {
		var r EventReference
		err = rows.Scan(&r.ID, &r.Period) // added missing fields
		if err != nil {
			log.Fatal(err)
		}
		records = append(records, r)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return records

}

// Удалить все из event_reference
func DeleteAllEventReferences() {
	db := config.GetDB()

	if db == nil {
		log.Printf("DB connection is not established")
	}

	db.Query("DELETE FROM event_reference")
}
