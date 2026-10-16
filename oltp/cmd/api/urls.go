package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/devendrasankhla/short-my-url/oltp/internals/database"
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

	//NOTE: if alias is there, it will use it, else it will add encoded snowflake id as alias
	shortKey := ""
	uniqueId := app.snowflakeNode.Generate()
	fmt.Println("New ID generated for ", req, " ", uniqueId)

	if req.Alias != "" {
		//TODO: check if alias already exists
		shortKey = req.Alias
	} else {
		s, _ := sqids.New()
		shortKey, err = s.Encode([]uint64{uint64(uniqueId.Int64())})
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	url := database.URL{
		ID:          int64(uniqueId),
		OriginalUrl: req.OriginalUrl,
		ShortKey:    shortKey,
		CreatedAt:   time.Now(),
		ModifiedAt:  time.Now(),
		TotalHits:   1,
	}

	err = app.models.Urls.AddUrl(uniqueId.Int64(), url)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Url Shorted %s\n ", (r.Host + "/" + shortKey))))
}

func (app *application) getShortUrl(w http.ResponseWriter, r *http.Request) {
	defer func(start time.Time) {
		duration := time.Since(start)
		slog.Info("execution_time",
			"duration_ms", duration.Milliseconds(), // Will show 0 if under 1ms
			"duration_str", duration.String(), // Will show exact time like "341µs" or "820ns"
		)
	}(time.Now())

	shortKey := r.PathValue("id")
	sqd, err := sqids.New()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//TODO: need to check here if we got system shortKey or user shortKey, use sqd to decode and check in DB with ID, if returns nothings, search with short key
	idSlice := sqd.Decode(shortKey)
	url := database.URL{}
	//if len is 0, need to search with shortKey
	if len(idSlice) > 0 {
		url, err = app.models.Urls.GetUrl(int64(idSlice[0]))
		if err != nil {
			//if err != nil, two cases, actual error, or not found, if not found, need to check with shortKey
			if err.Error() == "Not Found!" {
				fmt.Println("Url not found when searching with ID, need to check with shortKey")
			}
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	err = app.messaageQueue.Push(url.ShortKey)
	if err != nil {
		fmt.Println(err)
	}

	http.Redirect(w, r, url.OriginalUrl, http.StatusFound)
}

func (app *application) getUrlDetails(w http.ResponseWriter, r *http.Request) {
	defer func(start time.Time) {
		duration := time.Since(start)
		slog.Info("execution_time",
			"duration_ms", duration.Milliseconds(), // Will show 0 if under 1ms
			"duration_str", duration.String(), // Will show exact time like "341µs" or "820ns"
		)
	}(time.Now())

	//TODO: need to fix this after adding snowflake ids
	//urlID := r.PathValue("id")
	url, err := app.models.Urls.GetUrl(1)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-type", "application/json")
	jsonUrl, err := json.MarshalIndent(url, "", "\t")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(jsonUrl)
}
