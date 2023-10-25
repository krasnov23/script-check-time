package services

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Добавляет Ивент и Создает очередь
func AddAnEventAndAddQueueWithIdEvent(name string, url string, insertedPeriod string) {

	id, period := AddNewEventWithGetIdRecordBack(name, url, insertedPeriod) // (int, Event)

	fmt.Println(id, period)

	for nameDayOfWeek, arrayOfTimes := range period.DaysOfWeek {

		fmt.Println(nameDayOfWeek, arrayOfTimes)
		datesSlice := setNextWeekday(nameDayOfWeek, arrayOfTimes) // Вернет массив с датами дату формата 2023-10-27 09:58:57.83082 где 9:58 время которое сейчас

		for index := range datesSlice {
			AddNewQueue(id, datesSlice[index])
		}

	}

	for _, date := range period.CertainDate {

		parts := strings.Split(date, " ")
		fmt.Println(parts)

		day, _ := strconv.Atoi(parts[0])

		hourMinuteAndSeconds := strings.Split(parts[1], ":")

		hours, _ := strconv.Atoi(hourMinuteAndSeconds[0])
		minutes, _ := strconv.Atoi(hourMinuteAndSeconds[1])
		seconds, _ := strconv.Atoi(hourMinuteAndSeconds[2])

		todayDay := time.Now().Day()

		month := time.Now().Month()

		if day >= todayDay {

			t := time.Date(time.Now().Year(), month, day, hours, minutes, seconds, 0, time.UTC)

			AddNewQueue(id, t)
		} else {

			t := time.Date(time.Now().Year(), month, day, hours, minutes, seconds, 0, time.UTC).AddDate(0, 1, 0)

			AddNewQueue(id, t)
		}

	}

	if period.Odd != "" {

		// Получаем сегодняшний день
		todayDay := time.Now().Day()

		year, month, _ := time.Now().Date()
		firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
		firstOfNextMonth := firstOfMonth.AddDate(0, 1, 0)
		numberOfDays := firstOfNextMonth.Add(-time.Nanosecond).Day()

		periodOdd := strings.Split(period.Odd, ":")

		hours, _ := strconv.Atoi(periodOdd[0])
		minutes, _ := strconv.Atoi(periodOdd[1])
		seconds, _ := strconv.Atoi(periodOdd[2])

		numbersDays := make([]int, numberOfDays)

		for num := range numbersDays {
			numbersDays[num] = num + 1
		}

		//fmt.Println(numbersDays)

		for _, day := range numbersDays {
			fmt.Println(day)
			if day%2 != 0 && day >= todayDay {
				t := time.Date(time.Now().Year(), month, day, hours, minutes, seconds, 0, time.UTC) //.AddDate(0, 1, 0)
				AddNewQueue(id, t)
			}
		}

	}

	if period.Even != "" {

		year, month, _ := time.Now().Date()
		firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
		firstOfNextMonth := firstOfMonth.AddDate(0, 1, 0)
		numberOfDays := firstOfNextMonth.Add(-time.Nanosecond).Day()

		periodEven := strings.Split(period.Even, ":")

		hours, _ := strconv.Atoi(periodEven[0])
		minutes, _ := strconv.Atoi(periodEven[1])
		seconds, _ := strconv.Atoi(periodEven[2])

		numbersDays := make([]int, numberOfDays)

		todayDay := time.Now().Day()

		for num := range numbersDays {
			numbersDays[num] = num + 1
		}

		for _, day := range numbersDays {
			if day%2 == 0 && day >= todayDay {
				t := time.Date(time.Now().Year(), month, day, hours, minutes, seconds, 0, time.UTC) //.AddDate(0, 1, 0)
				AddNewQueue(id, t)
			}
		}
	}

	if period.LastDayOfMonth != "" {

		now := time.Now()

		currentYear, currentMonth, _ := now.Date()

		firstDayOfNextMonth := time.Date(currentYear, currentMonth+1, 1, 0, 0, 0, 0, now.Location())

		lastDayOfCurrentMonth := firstDayOfNextMonth.AddDate(0, 0, -1)

		periodLastDay := strings.Split(period.LastDayOfMonth, ":")

		hours, _ := strconv.Atoi(periodLastDay[0])
		minutes, _ := strconv.Atoi(periodLastDay[1])
		seconds, _ := strconv.Atoi(periodLastDay[2])

		timeWithHours := lastDayOfCurrentMonth.Add(time.Hour*time.Duration(hours) +
			time.Minute*time.Duration(minutes) + time.Second*time.Duration(seconds))
		AddNewQueue(id, timeWithHours)
	}

}

//now := time.Now()
//year, month, day := now.Date()

//layout := "2006-01-02 15:04:05"

// Concatenate current date and incoming time for parsing
//dateTimeStr := fmt.Sprintf("%d-%02d-%02d %s", period)

//myTime, err := time.Parse(layout, period.DaysOfWeek)

// Handle error
/*if err != nil {
	fmt.Println("Parsing error", err.Error())
} else {
	fmt.Println("Parsed time in current date: ", myTime)
}*/

func CheckQueueTimeAndMakeNewEvent() {

	// Получаем очереди, которые уже прошли (время которых истекло) и возвращаем слайс объектов их event_id и дат
	pastQueues := GetExpiredQueues()

	// Создаем массив из айди ивентов
	queuesIds := make([]int, len(pastQueues))

	for _, id := range pastQueues {
		queuesIds = append(queuesIds, id.EventID)
	}

	if len(pastQueues) != 0 {
		// По полученным event_id находим все
		//event и возвращаем слайс объектов EventReference(ID,Period)
		eventReferences := FindEventsByIdsAndGetIdAndPeriod(queuesIds)

		AddNewQueueByEventIdAndJsonEventData(eventReferences)

	}

	// Удаляем прошедшие очереди
	DeleteExpireQueues()
}

func setNextWeekday(name string, arrayOfTimes []string) []time.Time {

	var weekday time.Weekday

	switch name {
	case "Sunday":
		weekday = time.Sunday
	case "Monday":
		weekday = time.Monday
	case "Tuesday":
		weekday = time.Tuesday
	case "Wednesday":
		weekday = time.Wednesday
	case "Thursday":
		weekday = time.Thursday
	case "Friday":
		weekday = time.Friday
	case "Saturday":
		weekday = time.Saturday
	default:
		fmt.Println("Not Valid day of Week")
		return nil
	}

	now := time.Now()
	daysAhead := (weekday - now.Weekday() + 7) % 7

	if daysAhead == 0 {
		daysAhead = 7
	}

	rightTime := now.AddDate(0, 0, int(daysAhead))

	var dates []time.Time

	for _, individualTime := range arrayOfTimes {

		parts := strings.Split(individualTime, ":")
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

		currentYear, currentMonth, currentDay := rightTime.Date()

		t := time.Date(currentYear, currentMonth, currentDay, hours, minutes, seconds, 0, time.UTC)

		dates = append(dates, t)
	}

	return dates

}

/*`{
	"periodInDays": 0,
	"periodInMonths": 1,
	"daysOfWeek": { "Monday": ["9:00:00", "17:00:00"], "Tuesday": ["10:00:00", "18:00:00"] },
	"certainDate": null,
	"lastDay" : "12-00-00"
}`*/

/*s := `{
	"periodInDays": 7,
	"periodInMonths": 0,
	"daysOfWeek": null
	"certainDate": ["2023-10-01 15:00:00", "2023-11-01 15:00:00"]
	"lastDay": "12-00-00"
}`*/
