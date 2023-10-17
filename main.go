package main

import (
	"script-check-time/config"
	"script-check-time/services"
)

func main() {

	config.Connect()

	services.CheckQueueTimeAndMakeNewEvent()

	//services.AddAnEventAndAddQueueWithIdEvent("Dream", "lunch.com", "17:05:00")

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
