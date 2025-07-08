package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /v1/url", app.createShortUrl)
	mux.HandleFunc("GET /v1/url/{id}", app.getShortUrl)

	return mux
}
