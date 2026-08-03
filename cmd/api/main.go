package main

import (
	"fmt"

	"github.com/devendrasankhla/short-my-url/internals/database"
)

type application struct {
	models  database.Models
	counter uint64
}

func main() {
	db, err := database.NewFileDB("D:\\Downloads\\urlShortner.text")
	if err != nil {
		fmt.Println("Error initializing database : ", err)
	}

	app := &application{
		models:  database.NewModels(db),
		counter: 0,
	}

	err = app.serve()
	if err != nil {
		fmt.Println(err)
	}
}
