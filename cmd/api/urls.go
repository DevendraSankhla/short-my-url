package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/devendrasankhla/short-my-url/internals/database"
	"github.com/sqids/sqids-go"
)

type CreateShortUrlRequest struct {
	OriginalUrl string `json:"originalUrl"`
	Alias       string `json:"alias"`
}

func (app *application) createShortUrl(w http.ResponseWriter, r *http.Request) {
	var req CreateShortUrlRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println(err)
	}

	if req.OriginalUrl == "" {
		fmt.Println(errors.New("Need URL"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	urlId := ""
	if req.Alias != "" {
		urlId = req.Alias
	} else {
		s, _ := sqids.New()
		urlId, err = s.Encode([]uint64{app.counter})
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	url := database.URL{
		ID:          int64(app.counter),
		OriginalUrl: req.OriginalUrl,
		CreatedAt:   time.Now(),
		ModifiedAt:  time.Now(),
		TotalHits:   1,
	}

	err = app.models.Urls.AddUrl(urlId, url)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	app.counter++

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Url Shorted %s\n ", (r.Host + "/" + urlId))))
}

func (app *application) getShortUrl(w http.ResponseWriter, r *http.Request) {
	urlID := r.PathValue("id")
	originalUrl, err := app.models.Urls.GetUrl(urlID)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	http.Redirect(w, r, originalUrl, http.StatusFound)
}
