package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type Database interface {
	Insert(string, URL) error
	//Update()
	Get(string) (URL, error)
	//Delete()
}

type FileDB struct {
	filePath string
	mu       sync.Mutex
	data     map[string]URL
}

func NewFileDB(filePath string) (*FileDB, error) {
	db := &FileDB{
		filePath: filePath,
		data:     make(map[string]URL),
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

func (db *FileDB) Insert(id string, url URL) error {
	if _, ok := db.data[id]; ok {
		return errors.New("Duplicate URL")
	}
	db.data[id] = url
	err := db.saveToFile()
	if err != nil {
		return err
	}
	return nil
}

func (db *FileDB) Get(id string) (URL, error) {
	err := db.loadDataFromFile()
	if err != nil {
		return URL{}, err
	}

	urlStruct, ok := db.data[id]
	if !ok {
		return URL{}, errors.New("Not Found!")
	}

	return urlStruct, nil
}
