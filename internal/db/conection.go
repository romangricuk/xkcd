package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"sync"
)

var db *sql.DB

func GetInstance() (*sql.DB, error) {
	var err error

	once := new(sync.Once)
	once.Do(func() {
		db, err = sql.Open("sqlite3", "file:db.sqlite")
	})

	if err != nil {
		err = fmt.Errorf("on connect to database: %w", err)
		return nil, err
	}

	return db, nil
}
