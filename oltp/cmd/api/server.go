package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) serve() error {
	server := &http.Server{
		Addr:    ":" + strconv.Itoa(app.cfg.port),
		Handler: app.routes(),
	}

	fmt.Println("Starting server at ", string(app.cfg.port))

	err := server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
