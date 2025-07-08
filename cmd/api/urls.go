package main

import "net/http"

func (app *application) createShortUrl(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Created short url"))
}

func (app *application) getShortUrl(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("redirect to ./shorurl"))
}
