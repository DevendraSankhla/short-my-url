package main

import (
	"fmt"
	"net/http"

	messagequeue "github.com/devendrasankhla/short-my-url/analytics/internals/messageQueue"
)

func healthTest(w http.ResponseWriter, r *http.Request) {
	mq, _ := messagequeue.NewRabbmitMQ()
	mq.Pull()

	w.Write([]byte(fmt.Sprint("Helllo, From analyics microservice")))
}

func main() {
	fmt.Println("Welcome to analytics microservice, server started on 3000")
	http.HandleFunc("/", healthTest)
	http.ListenAndServe(":3000", nil)
}
