package services

import (
	"fmt"
	"time"
)

// Добавляет Ивент и Создает очередь
func AddAnEventAndAddQueueWithIdEvent(name string, url string, insertedPeriod string) {

	id, period := AddNewEventWithGetIdRecordBack(name, url, insertedPeriod) // (int, Event)

	//now := time.Now()
	//year, month, day := now.Date()

	layout := "2006-01-02 15:04:05"

	// Concatenate current date and incoming time for parsing
	//dateTimeStr := fmt.Sprintf("%d-%02d-%02d %s", period)

	// Parsing the string to time
	myTime, err := time.Parse(layout, period.StartDate)

	// Handle error
	if err != nil {
		fmt.Println("Parsing error", err.Error())
	} else {
		fmt.Println("Parsed time in current date: ", myTime)
	}

	AddNewQueue(id, myTime)
}

func CheckQueueTimeAndMakeNewEvent() {

	// Получаем очереди, которые уже прошли (время которых истекло) и возвращаем слайс объектов их event_id и дат
	pastQueues := GetExpiredQueues()

	// Создаем массив из айди очередей

	if len(pastQueues) != 0 {
		// По полученным event_id находим все
		//event и возвращаем слайс объектов EventReference(ID,Period)
		eventReferences := FindEventsByIdsAndGetIdAndPeriod(pastQueues)

		// Удаляем прошедшие очереди
		DeleteExpireQueues()

		// Создаем эти же очереди на завтра
		for _, event := range eventReferences {
			AddNewQueueByEventDataAndEditEventDate(event.ID, event.Period)
		}
	}

}
