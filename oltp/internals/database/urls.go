package database

import (
	"fmt"
	"time"
)

type URL struct {
	ID          int64     `json:"id"`
	OriginalUrl string    `json:"originalUrl"`
	ShortKey    string    `json:"alias"`
	CreatedAt   time.Time `json:"createdAt"`
	ModifiedAt  time.Time `json:"modifiedAt"`
	TotalHits   int       `json:"totalHits"`
}

type URlModel struct {
	db Database
}

func (m *URlModel) AddUrl(id int64, url URL) error {
	err := m.db.Insert(id, url)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (m *URlModel) GetUrl(id int64) (URL, error) {
	url, err := m.db.Get(id)
	if err != nil {
		return URL{}, err
	}
	return url, nil
}
