package main

import (
	"script-check-time/config"
	"script-check-time/services"
)

func main() {

	config.Connect()

	services.CheckQueueTimeAndMakeNewEvent()

	// Добавление событий и очередей с айдишниками событий
	//services.AddAnEventAndAddQueueWithIdEvent("Проснуться", "up.com", "11:01:00")

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
