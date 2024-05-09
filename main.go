package main

import (
	"fmt"
	"xkcd/internal/db"
	"xkcd/internal/models"
	"xkcd/internal/parser"
)

func main() {

	conn, err := db.GetInstance()
	defer func() {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}()

	if err != nil {
		panic(err)
	}

	err = models.CreateTables()
	if err != nil {
		err = fmt.Errorf("on create table: %w", err)
		panic(err)
	}

	isEnd := false

	for !isEnd {
		isEnd, err = parser.Parse(models.Comic{})

		if err != nil {
			err = fmt.Errorf("on parse: %w", err)
			panic(err)
		}
	}
}
