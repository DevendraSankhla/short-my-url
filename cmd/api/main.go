package main

import (
	"fmt"
)

type application struct {
	counter    uint64
	shornedUrl map[string]string
}

func main() {
	app := &application{
		counter:    0,
		shornedUrl: make(map[string]string),
	}

	err := app.serve()
	if err != nil {
		fmt.Println(err)
	}
}
