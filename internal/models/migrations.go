package models

import (
	"fmt"
	"xkcd/internal/db"
)

func CreateTables() error {
	conn, err := db.GetInstance()
	if err != nil {
		err = fmt.Errorf("on GetInstance: %w", err)
		return err
	}

	sqlString := `CREATE TABLE IF NOT EXISTS comic (
		id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		link string,
		news string,
		safetitle string,
		transcript string,
		alt string,
		img string,
		title string,
		year int,
		month int,
		day int,
		num string
	);`

	_, err = conn.Exec(sqlString)

	if err != nil {
		err = fmt.Errorf("create table comic: %w", err)
		return err
	}

	return nil
}
