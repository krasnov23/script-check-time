package main

import (
	"fmt"
	"script-check-time/config"
	"script-check-time/services"
	"time"
)

func main() {

	config.Connect()

	//addQueue()
	for {
		services, _ := services.GetAllQueues(time.Now())

		for i := range services {
			fmt.Println(services[i])
		}

		time.Sleep(60 * time.Second)
	}
}

func addQueue() {

	/*layout := time.RFC3339
	str := "2006-01-02T15:04:05Z"

	t, err := time.Parse(layout, str)

	if err != nil {
		fmt.Println(err)
	}*/

	now := time.Now()
	times := []time.Time{
		now, now.Add(time.Second * 40),
		now.Add(time.Second * 70),
		now.Add(time.Second * 80),
		now.Add(time.Second * 90),
		now.Add(time.Second * 125),
	}

	for i := range times {
		services.AddNewQueue(fmt.Sprintf("Event number is %d", i), times[i])
	}

}
