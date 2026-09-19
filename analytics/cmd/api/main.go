package main

import (
	"fmt"

	messagequeue "github.com/devendrasankhla/short-my-url/analytics/internals/messageQueue"
)

type application struct {
	messaageQueue messagequeue.MessageQueue
}

func main() {
	rabbitMQ, err := messagequeue.NewRabbmitMQ()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rabbitMQ.Close()

	app := &application{
		messaageQueue: rabbitMQ,
	}

	app.messaageQueue.Pull()

	err = app.serve()
	if err != nil {
		fmt.Println(err)
	}
}
