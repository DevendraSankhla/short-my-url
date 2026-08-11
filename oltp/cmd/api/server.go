package main

import (
	"fmt"
	"net/http"
)

func (app *application) serve() error {
	server := &http.Server{
		Addr:    ":4000",
		Handler: app.routes(),
	}

	fmt.Println("Starting server at 4000")

	err := server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
