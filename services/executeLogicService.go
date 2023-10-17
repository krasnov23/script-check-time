package services

import (
	"fmt"
	"time"
)

// Добавляет Ивент и Создает очередь
func AddAnEventAndAddQueueWithIdEvent(name string, url string, insertedPeriod string) {

	id, period := AddNewEventWithGetIdRecordBack(name, url, insertedPeriod)

	now := time.Now()
	year, month, day := now.Date()

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

func CheckQueueTimeAndMakeNewEvent() {

	// Получаем очереди, которые уже прошли (время которых истекло) и возвращаем слайс их id
	pastQueues := GetExpiredQueues()

	// По полученным event_id находим все event и возвращаем слайс объектов EventReference(ID,Period)
	eventReferences := FindEventsByIdsAndGetIdAndPeriod(pastQueues)

	// Удаляем прошедшие очереди
	DeleteExpireQueues()

	// Создаем эти же очереди на завтра
	for _, event := range eventReferences {
		AddNewQueueByEventData(event.ID, event.Period)
	}
}
