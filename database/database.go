package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/Viet-ph/xss-vulnerable/models"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Comments map[int]models.Comment `json:"comments"`
	Users  map[int]models.User  `json:"users"`
}

var (
	Db *DB
)

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		file, _ := os.Create(path)
		defer file.Close()
	}

	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}
	return db, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) EnsureDB() error {
	_, err := os.Stat(db.path)
	if errors.Is(err, os.ErrNotExist) {
		file, errCreate := os.Create(db.path)
		if errCreate != nil {
			return errCreate
		}

		defer file.Close()
	}

	return nil
}

// loadDB reads the database file into memory
func (db *DB) LoadDB() (DBStructure, error) {
	file, errReadFile := os.ReadFile(db.path)
	data := DBStructure{}

	if errReadFile != nil {
		return data, errReadFile
	}

	json.Unmarshal(file, &data)
	if data.Comments == nil {
		data.Comments = map[int]models.Comment{}
	}
	if data.Users == nil {
		data.Users = map[int]models.User{}
	}

	return data, nil
}

// writeDB writes the database file to disk
func (db *DB) WriteDB(dbStructure DBStructure) error {
	data, errMarshal := json.MarshalIndent(dbStructure, "", " ")
	if errMarshal != nil {
		return errMarshal
	}

	db.mux.Lock()
	if errWriteFile := os.WriteFile(db.path, data, 0666); errWriteFile != nil {
		return errWriteFile
	}
	defer db.mux.Unlock()

	return nil
}
