package main

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	homeFS, err := fs.Sub(embedFrontend, "web")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	mux.HandleFunc("POST /v1/url", app.createShortUrl)
	mux.HandleFunc("GET /{id}", app.getShortUrl)
	mux.Handle("GET /", http.FileServer(http.FS(homeFS)))

	return mux
}
