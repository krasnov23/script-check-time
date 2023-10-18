package services

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"script-check-time/config"
	"strconv"
)

type EventReference struct {
	ID     int
	Period string
}

type Event struct {
	PeriodInDays   int    `json:"periodInDays"`
	PeriodInMonths int    `json:"periodInMonths"`
	StartDate      string `json:"startDate"`
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

func AddNewEventWithGetIdRecordBack(name string, url string, insertedPeriod string) (int, Event) {

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
		log.Printf("Error with scanning and returning id,period: %s", err)
	}

	fmt.Println(period)

	var e Event
	err = json.Unmarshal([]byte(period), &e)

	if err != nil {
		log.Printf("Problem with unmarshall: %s", err)
	}

	return id, e

}

func FindEventsByIdsAndGetIdAndPeriod(ids []int) []EventReference {

	db := config.DB

	strIDs := make([]interface{}, len(ids))
	for i, id := range ids {
		strIDs[i] = id
	}

	fmt.Println(strIDs)

	var placeholder string
	//placeholders := strings.Repeat(",?", len(ids)-1)
	for i, _ := range ids {
		if i != 0 {
			a := i + 1
			num := strconv.Itoa(a)
			placeholder += "," + "$" + num
		}

	}

	fmt.Println(reflect.TypeOf(placeholder))

	query := fmt.Sprintf(`SELECT id, period FROM event_reference WHERE id IN ($1%s)`, placeholder)
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

	fmt.Println(records)
	return records

}

func EditEventsById(id int, event Event) {

	db := config.DB

	jsonData, err := json.Marshal(&event)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = db.Exec(`UPDATE event_reference SET period = $1 WHERE id = $2`, jsonData, id)
	if err != nil {
		panic(err)
	}

	fmt.Println("Record updated successfully!")

}

// Удалить все из event_reference
func DeleteAllEventReferences() {
	db := config.GetDB()

	if db == nil {
		log.Printf("DB connection is not established")
	}

	db.Query("DELETE FROM event_reference")
}
