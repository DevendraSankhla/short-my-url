package main

import "fmt"

type application struct {
}

func main() {
	app := &application{}

	err := app.serve()
	if err != nil {
		fmt.Println(err)
	}
}
