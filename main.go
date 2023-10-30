package main

import (
	"script-check-time/config"
)

func main() {

	config.Connect()

	/*services.AddAnEventAndAddQueueWithIdEvent("Сходить в зал", "gym.com",
	`{ "odd" : "15:05:00" }`)*/

	//services.CheckQueueTimeAndMakeNewEvent()

	//config.DeleteEventTable()
	//config.DeleteQueueTable()

}

// Тестовые данные для проверки

/*services.AddAnEventAndAddQueueWithIdEvent("Пробник", "test.com",
`{"periodInMonths": 1,"CertainDate" : ["25 15:35:00","25 15:40:00","25 15:43:00"]}`)*/

/*services.AddAnEventAndAddQueueWithIdEvent("Заплатит зарплату", "salary.com",
`{ "lastDayOfMonth" : "20:00:00" }`)*/

/*services.AddAnEventAndAddQueueWithIdEvent("Проснуться", "son.com",
`{"periodInDays": 1,"CertainDate" : ["23 16:25:00","24 19:00:00"]}`)*/

/*services.AddAnEventAndAddQueueWithIdEvent("Потренероваться", "son.com",
`{"DaysOfWeek" : {"Sunday" : ["18:00:00","22:00:00"],
	"Friday" : ["18:00:00","22:00:00"],"Monday": ["18:00:00","22:00:00"]}}`)*/

/*services.AddAnEventAndAddQueueWithIdEvent("Сходить в зал", "gym.com",
`{ "even" : "15:05:00" }`)*/

/*services.AddAnEventAndAddQueueWithIdEvent("Сходить в зал", "gym.com",
`{ "odd" : "15:05:00" }`)*/
