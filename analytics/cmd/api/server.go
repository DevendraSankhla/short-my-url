package main

import (
	"fmt"
	"net/http"
)

func (app *application) serve() error {
	server := &http.Server{
		Addr: ":3000",
	}

	fmt.Println("Starting server at 3000")

	err := server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
