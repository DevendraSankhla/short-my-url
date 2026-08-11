package database

import (
	"fmt"
	"time"
)

type URL struct {
	ID          int64     `json:"id"`
	OriginalUrl string    `json:"originalUrl"`
	Alias       string    `json:"alias"`
	CreatedAt   time.Time `json:"createdAt"`
	ModifiedAt  time.Time `json:"modifiedAt"`
	TotalHits   int       `json:"totalHits"`
}

type URlModel struct {
	db Database
}

func (m *URlModel) AddUrl(id string, url URL) error {
	err := m.db.Insert(id, url)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (m *URlModel) GetUrl(id string) (string, error) {
	url, err := m.db.Get(id)
	if err != nil {
		return "", err
	}
	return url.OriginalUrl, nil
}
