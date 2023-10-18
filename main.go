package main

import (
	"script-check-time/config"
	"script-check-time/services"
)

func main() {

	config.Connect()

	services.CheckQueueTimeAndMakeNewEvent()
	/*_, eventObject := services.AddNewEventWithGetIdRecordBack("Проснуться", "son.com",
	`{"periodInDays": 1,"startDate":"2023-10-18 15:04:05","periodInMonths" : 0}`)*/

	//services.CheckQueueTimeAndMakeNewEvent()

	// Добавление событий и очередей с айдишниками событий
	/*services.AddAnEventAndAddQueueWithIdEvent("Пообедать", "lunch.com",
	`{"periodInDays": 1,"startDate":"2023-10-18 16:04:05","periodInMonths": 0}`)*/

	// Добавление очередей
	/*services.AddQueue()

	for {
		queues, _ := services.GetAllQueuesByDate(time.Now())

		for i := range queues {
			fmt.Println(queues[i])
		}

		time.Sleep(60 * time.Second)
	}*/

}
