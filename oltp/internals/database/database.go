package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type Database interface {
	Insert(int64, URL) error
	//Update()
	Get(int64) (URL, error)
	//Delete()
}

type FileDB struct {
	filePath string
	mu       sync.Mutex
	data     map[int64]URL
}

func NewFileDB(filePath string) (*FileDB, error) {
	db := &FileDB{
		filePath: filePath,
		data:     make(map[int64]URL),
	}

	//if file does not exists
	if _, err := os.Stat(db.filePath); os.IsNotExist(err) {
		err = db.saveToFile()
		if err != nil {
			return nil, err
		}
		return db, err
	}

	//if it does
	err := db.loadDataFromFile()
	if err != nil {
		return nil, err
	}
	return db, err
}

func (db *FileDB) loadDataFromFile() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	bytes, err := os.ReadFile(db.filePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &db.data)
	if err != nil {
		return err
	}

	return nil
}

func (db *FileDB) saveToFile() error {
	bytes, err := json.MarshalIndent(db.data, "", "")
	if err != nil {
		return err
	}
	return os.WriteFile(db.filePath, bytes, 0644)
}

func (db *FileDB) Insert(id int64, url URL) error {
	if _, ok := db.data[id]; ok {
		return errors.New("Duplicate URL")
	}
	db.data[id] = url
	err := db.saveToFile()
	if err != nil {
		return err
	}

	err = db.loadDataFromFile()
	if err != nil {
		return err
	}
	return nil
}

func (db *FileDB) Get(id int64) (URL, error) {
	urlStruct, ok := db.data[id]
	if !ok {
		return URL{}, errors.New("Not Found!")
	}

	return urlStruct, nil
}
