package services

import (
	"encoding/json"
	"fmt"
	"log"
	"script-check-time/config"
	"strconv"
	"strings"
	"time"
)

type Queue struct {
	EventID int
	Date    time.Time
}

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

func AddNewQueueByEventIdAndJsonEventData(eventsReferences []EventReference) {

	// id - Event
	eventsSlice := make(map[int]Event)

	for _, event := range eventsReferences {

		var e Event
		err := json.Unmarshal([]byte(event.Period), &e)

		if err != nil {
			log.Printf("Problem with unmarshall: %s", err)
		}

		eventsSlice[event.ID] = e
	}

	for id, event := range eventsSlice {

		if event.LastDayOfMonth != "" {

			// Находим дату последнюю дату предыдущего месяца вставляем в нее время из event. LastDayOfMonth
			t := time.Now()

			firstDayCurrentMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)

			lastDayPreviousMonth := firstDayCurrentMonth.AddDate(0, 0, -1)

			parts := strings.Split(event.LastDayOfMonth, ":")

			hours, err := strconv.Atoi(parts[0])

			if err != nil {
				fmt.Println("Parsing error with hours", err.Error())
			}

			minutes, err := strconv.Atoi(parts[1])

			if err != nil {
				fmt.Println("Parsing error with minutes", err.Error())
			}

			seconds, err := strconv.Atoi(parts[2])

			if err != nil {
				fmt.Println("Parsing error with seconds", err.Error())
			}

			desiredData := time.Date(lastDayPreviousMonth.Year(), lastDayPreviousMonth.Month(), lastDayPreviousMonth.Day(),
				hours, minutes, seconds, 0, time.UTC)

			// Получили id последнего дня предыдущего месяца
			idOfLastDayOfPreviousMonth := findQueueIdByDateAndEventId(id, desiredData)

			// Удаляем очередь по id, если такая у нас была
			if idOfLastDayOfPreviousMonth != 0 {
				DeleteQueueById(id)
			}

			if t.Day() >= 20 {
				firstDayOfAfterNextMonth := time.Date(t.Year(), t.Month()+2, 1, hours, minutes, seconds, 0, time.UTC)

				lastDayOfNextMonth := firstDayOfAfterNextMonth.AddDate(0, 0, -1)

				isExists := findQueueIdByDateAndEventId(id, lastDayOfNextMonth)

				if isExists == 0 {
					AddNewQueue(id, lastDayOfNextMonth)
				}

			}

			firstDayOfNextMonth := time.Date(t.Year(), t.Month()+1, 1, hours, minutes, seconds, 0, time.UTC)

			lastDayOfCurrentMonth := firstDayOfNextMonth.AddDate(0, 0, -1)

			isExists := findQueueIdByDateAndEventId(id, lastDayOfCurrentMonth)

			if isExists == 0 {
				AddNewQueue(id, lastDayOfCurrentMonth)
			}
		}

		if len(event.CertainDate) > 0 {

			// Создаем новый слайс длинной в длину certainDate
			newCertainDates := make([]time.Time, len(event.CertainDate))

			// Получаем слайс последних истекших дат
			expiredDates := GetExpiredQueueByEventId(id)

			// Реализуем логику с февралем
			//eventStartDates := event.CertainDate

			for _, date := range expiredDates {

				isFebruary := (int(date.Month()) + event.PeriodInMonths) % 12

				term := (date.Day() == 29 || date.Day() == 30 || date.Day() == 31) && event.PeriodInMonths > 0 && isFebruary == 2

				if term {
					now := time.Now()

					var newDate time.Time

					str := "2024-02-25"

					example, err := time.Parse("2006-01-02", str)

					if err != nil {
						fmt.Println(err)
					}

					year := int(date.Month()) + event.PeriodInMonths
					if year > 12 {
						year = now.Year() + 1
					} else {
						year = now.Year()
					}

					if year%4 == 0 {
						newDate = time.Date(year, example.Month(), 29, date.Hour(),
							date.Minute(), date.Second(), 0, time.UTC)
					} else {
						newDate = time.Date(year, example.Month(), 28, date.Hour(),
							date.Minute(), date.Second(), 0, time.UTC)
					}

					newCertainDates = append(newCertainDates, newDate)
				}

				// Добавляем дни и месяцы к дате исходя из нашего ивента
				newDate := date.AddDate(0, event.PeriodInMonths, event.PeriodInDays)

				// Добавляем даты в слайс которые мы хотим добавить в queue
				newCertainDates = append(newCertainDates, newDate)
			}

			// У нас есть новые даты, теперь нам нужно проверить есть ли они уже в базе
			// Получаем даты, которые уже есть в нашей базе []time.Time
			existingDates := GetExistingDatesByEvent(id, newCertainDates)

			for n, newDate := range newCertainDates {
				for _, existingDate := range existingDates {
					if newDate == existingDate {
						// Если очередь которую мы собираемся создать уже есть в таблице очередей удаляем ее из таблицы
						// :n - создает слайс от первого элемента до того который мы удаляем
						// newCertainDates[n+1:]... - вставляет все элементы после удаляемого элемента
						newCertainDates = append(newCertainDates[:n], newCertainDates[n+1:]...)
					}
				}
			}

			for _, newDate := range newCertainDates {
				AddNewQueue(id, newDate)
			}
		} else if len(event.DaysOfWeek) != 0 {

			// Создаем новый слайс длинной в длину DaysOfWeek
			// newCertainDates := make([]time.Time, len(event.DaysOfWeek))

			// Получаем слайс последних истекших дат
			expiredDates := GetExpiredQueueByEventId(id)

			// Получаем дни недели по числам из истекших дат
			for _, weekDay := range expiredDates {

				// Узнаем название дня
				nameOfWeekDay := weekDay.Weekday() // type time.Weekday

				now := time.Now()
				daysAhead := (nameOfWeekDay - now.Weekday() + 7) % 7

				if daysAhead == 0 {
					daysAhead = 7
				}

				// Добавили количество дней до дня недели
				newDay := now.AddDate(0, 0, int(daysAhead))

				// Берем год месяц и число создаваемого дня
				currentYear, currentMonth, currentDay := newDay.Date()

				// Берем часы минуты и секунды создаваемого дня

				newFullDate := time.Date(currentYear, currentMonth, currentDay,
					weekDay.Hour(), weekDay.Minute(), weekDay.Second(), 0, time.UTC)

				AddNewQueue(id, newFullDate)
			}
		} else if event.Odd != "" {

			todayDay := time.Now().Day()

			executeLogicToAddNewQueueByOddOrEven("Odd", todayDay, id, event.Odd)

		} else if event.Even != "" {
			todayDay := time.Now().Day()

			executeLogicToAddNewQueueByOddOrEven("Even", todayDay, id, event.Even)
		}
	}

}

func executeLogicToAddNewQueueByOddOrEven(oddOrEven string, todayDay int, eventId int, hoursMinutesSeconds string) {

	if todayDay >= 27 {

		year, month, _ := time.Now().Date()
		firstOfMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
		// Получаем первое число того месяца который будет через один (например сейчас 27 октября получает 1 декабря)
		firstOfNextMonth := firstOfMonth.AddDate(0, 1, 0)
		// Получаем количество дней в следующем месяце
		numberOfDays := firstOfNextMonth.Add(-time.Nanosecond).Day()

		periodOdd := strings.Split(hoursMinutesSeconds, ":")

		hours, _ := strconv.Atoi(periodOdd[0])
		minutes, _ := strconv.Atoi(periodOdd[1])
		seconds, _ := strconv.Atoi(periodOdd[2])

		numbersDays := make([]int, numberOfDays)

		// Получаем слайс дней следующего месяца
		for num := range numbersDays {
			numbersDays[num] = num + 1
		}

		var checkNewQueuesDates []time.Time

		if oddOrEven == "Odd" {
			for _, day := range numbersDays {
				fmt.Println(day)
				if day%2 != 0 {

					newDate := time.Date(time.Now().Year(), month+1, day, hours, minutes, seconds, 0, time.UTC)

					checkNewQueuesDates = append(checkNewQueuesDates, newDate)
				}
			}
		} else {
			for _, day := range numbersDays {
				fmt.Println(day)
				if day%2 == 0 {

					newDate := time.Date(time.Now().Year(), month+1, day, hours, minutes, seconds, 0, time.UTC)

					checkNewQueuesDates = append(checkNewQueuesDates, newDate)
				}
			}
		}

		// Проверяем есть ли уже данные числа в очереди
		existingQueuesDates := GetExistingDatesByEvent(eventId, checkNewQueuesDates)

		// Находим те очереди, что не были созданы
		for n, newDate := range checkNewQueuesDates {
			for _, queueDate := range existingQueuesDates {
				if queueDate == newDate {
					checkNewQueuesDates = append(checkNewQueuesDates[:n], checkNewQueuesDates[n+1:]...)
				}
			}
		}

		// Создаем очереди, которые еще не были созданы
		for n := range checkNewQueuesDates {
			AddNewQueue(eventId, checkNewQueuesDates[n])
		}

	} else {

		year, month, _ := time.Now().Date()
		firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
		// Получаем первое число cледующего месяца
		firstOfNextMonth := firstOfMonth.AddDate(0, 1, 0)
		// Получаем количество дней в нашем месяце
		amountOfDays := firstOfNextMonth.Add(-time.Nanosecond).Day()

		periodOdd := strings.Split(hoursMinutesSeconds, ":")

		hours, _ := strconv.Atoi(periodOdd[0])
		minutes, _ := strconv.Atoi(periodOdd[1])
		seconds, _ := strconv.Atoi(periodOdd[2])

		numbersDays := make([]int, amountOfDays)

		// Получаем слайс дней действующего месяца
		for num := range numbersDays {
			numbersDays[num] = num + 1
		}

		var checkNewQueuesDates []time.Time

		if oddOrEven == "Odd" {
			for _, day := range numbersDays {
				fmt.Println(day)
				if day%2 != 0 {

					newDate := time.Date(time.Now().Year(), month, day, hours, minutes, seconds, 0, time.UTC)

					checkNewQueuesDates = append(checkNewQueuesDates, newDate)
				}
			}
		} else {
			for _, day := range numbersDays {
				fmt.Println(day)
				if day%2 == 0 {

					newDate := time.Date(time.Now().Year(), month+1, day, hours, minutes, seconds, 0, time.UTC)

					checkNewQueuesDates = append(checkNewQueuesDates, newDate)
				}
			}
		}

		// Проверяем есть ли уже данные числа в очереди
		existingQueuesDates := GetExistingDatesByEvent(eventId, checkNewQueuesDates)

		// Находим те очереди, что не были созданы
		for n, newDate := range checkNewQueuesDates {
			for _, queueDate := range existingQueuesDates {
				if queueDate == newDate {
					checkNewQueuesDates = append(checkNewQueuesDates[:n], checkNewQueuesDates[n+1:]...)
				}
			}
		}

		// Создаем очереди, которые еще не были созданы
		for n := range checkNewQueuesDates {
			AddNewQueue(eventId, checkNewQueuesDates[n])
		}
	}

}

func GetExistingDatesByEvent(id int, dates []time.Time) []time.Time {

	db := config.DB

	args := make([]interface{}, len(dates)+1)
	args[0] = id
	for i, date := range dates {
		args[i+1] = date
	}

	var placeholders string
	for i := range dates {
		if i != 0 {
			a := i + 2
			num := strconv.Itoa(a)
			placeholders += "," + "$" + num
		}
	}

	query := fmt.Sprintf(`SELECT date FROM queues WHERE event_id = $1 AND date IN ($2%s)`, placeholders)

	fmt.Println(query)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil
	}

	defer rows.Close()

	var resultDates []time.Time
	for rows.Next() {
		var date time.Time
		if err := rows.Scan(&date); err != nil {
			log.Fatal(err)
		}
		resultDates = append(resultDates, date)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return resultDates

}

// period - нужно распарсить из строки в объект Event
func AddNewQueueByEventDataAndEditEventDate(id int, period string) {

	var e Event
	err := json.Unmarshal([]byte(period), &e)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(e)

	//now := time.Now()
	//tomorrow := now.AddDate(0, 0, e.periodInDays)
	//year, month, day := tomorrow.Date()

	//layout := "2006-01-02 15:04:05"

	// Concatenate current date and incoming time for parsing
	//dateTimeStr := fmt.Sprintf("%d-%02d-%02d %s", year, month, day, period)

	// Parsing the string to time
	//myTime, err := time.Parse(layout, e.StartDate)

	// Handle error
	/*if err != nil {
		fmt.Println("Parsing error", err.Error())
	} else {
		fmt.Println("Parsed time in current date: ", myTime)
	}*/

	// Новая дата
	//newDate := myTime.AddDate(0, e.PeriodInMonths, e.PeriodInDays)
	//fmt.Println(newDate)

	//e.StartDate = newDate.Format(layout)
	//fmt.Println(e.StartDate)

	//EditEventsById(id, e)

	//AddNewQueue(id, newDate)
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

func GetExpiredQueueByEventId(id int) []time.Time {

	db := config.GetDB()

	if db == nil {
		log.Printf("DB connection is not established")
	}

	rows, err := db.Query(`
		SELECT date
		FROM queues
		WHERE event_id = $1 AND date < $2
		ORDER BY date DESC
	`, id, time.Now())

	if err != nil {
		panic(err)
	}

	var dates []time.Time
	for rows.Next() {
		var date time.Time
		err = rows.Scan(&date)
		if err != nil {
			panic(err)
		}
		dates = append(dates, date)
	}

	return dates

	/*for _, date := range dates {
		fmt.Println(date)
	}*/

}

func findQueueIdByDateAndEventId(id int, date time.Time) int {
	db := config.GetDB()

	if db == nil {
		log.Printf("DB connection is not established")
	}

	row, err := db.Query(`
		SELECT id
		FROM queues
		WHERE event_id = $1 AND date = $2
		ORDER BY date DESC
		LIMIT 1
	`, id, date)

	if err != nil {
		log.Printf("Error while executing the query: %v", err)
		return 0
	}

	defer row.Close()

	var queueId int
	if row.Next() {
		err := row.Scan(&queueId)
		if err != nil {
			log.Printf("Error while scanning the row: %v", err)
			return 0
		}
	}

	if err := row.Err(); err != nil {
		log.Printf("Row error: %v", err)
		return 0
	}

	return queueId

}

func GetExpiredQueues() []Queue {

	db := config.GetDB()

	if db == nil {
		log.Printf("DB connection is not established")
	}

	timeStamp := time.Now()

	rows, err := db.Query("SELECT event_id,date FROM queues WHERE date <= $1", timeStamp)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var queues []Queue

	for rows.Next() {
		var r Queue
		err = rows.Scan(&r.EventID, &r.Date) // added missing fields
		if err != nil {
			log.Fatal(err)
		}
		queues = append(queues, r)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(queues)
	return queues
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

func DeleteQueueById(id int) {
	db := config.GetDB()

	if db == nil {
		log.Printf("DB connection is not established")
	}

	rows, err := db.Query("DELETE FROM queues WHERE id = $1", id)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
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
