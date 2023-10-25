package main

import (
	"script-check-time/config"
)

func main() {

	config.Connect()

}

// Тестовые данные для проверки

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
