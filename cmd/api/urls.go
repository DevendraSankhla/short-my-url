package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sqids/sqids-go"
)

func (app *application) createShortUrl(w http.ResponseWriter, r *http.Request) {
	var input struct {
		OrignalUrl string `json:"orignalUrl"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		fmt.Println(err)
	}

	s, _ := sqids.New()
	id, _ := s.Encode([]uint64{app.counter})

	app.shornedUrl[id] = input.OrignalUrl

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Url Shorted %s\n ", (r.Host + "/" + id))))
}

func (app *application) getShortUrl(w http.ResponseWriter, r *http.Request) {
	urlID := r.PathValue("id")
	http.Redirect(w, r, app.shornedUrl[urlID], http.StatusFound)
}
