package main

import (
	"embed"
	"fmt"

	"github.com/devendrasankhla/short-my-url/oltp/internals/database"
	messagequeue "github.com/devendrasankhla/short-my-url/oltp/internals/messageQueue"
)

type application struct {
	models        database.Models
	counter       uint64
	messaageQueue messagequeue.MessageQueue
}

//go:embed web/*
var embedFrontend embed.FS

func main() {
	db, err := database.NewFileDB("D:\\Downloads\\urlShortner")
	if err != nil {
		fmt.Println("Error initializing database : ", err)
		return
	}

	rabbitMQ, err := messagequeue.NewRabbmitMQ()
	if err != nil {
		fmt.Println(err)
	}
	defer rabbitMQ.Close()

	app := &application{
		models:        database.NewModels(db),
		counter:       0,
		messaageQueue: rabbitMQ,
	}

	err = app.serve()
	if err != nil {
		fmt.Println(err)
	}
}
