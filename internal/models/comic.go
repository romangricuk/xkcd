package models

import (
	"database/sql"
	"fmt"
	"xkcd/internal/db"
)

type ComicRepo interface {
	Read(id int) (comic Comic, err error)
	Save(comic *Comic) (err error)
	Delete(comic Comic) error
	ReadAll(limit int, isDesc bool) ([]Comic, error)
	ReadLast() (Comic, error)
}

type Comic struct {
	ID         int    `json:"id"`
	Link       string `json:"link"`
	News       string `json:"news"`
	SafeTitle  string `json:"safe_title"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
	Title      string `json:"title"`
	Year       int    `json:"year"`
	Month      int    `json:"month"`
	Day        int    `json:"day"`
	Num        int    `json:"num"`
}

func (c Comic) Read(id int) (comic Comic, err error) {
	conn, err := db.GetInstance()
	if err != nil {
		err = fmt.Errorf("on GetInstance: %v", err)
		return comic, err
	}

	row, err := conn.Query(`select * from comic where id = ?`, id)

	if err != nil {
		err = fmt.Errorf("on query: %w", err)
		return comic, err
	}

	if err = scanRowToComic(row, &c); err != nil {
		return comic, err
	}

	return comic, nil
}

func scanRowToComic(row *sql.Rows, c *Comic) error {
	if err := row.Scan(
		&(c.ID),
		&(c.Link),
		&(c.News),
		&(c.SafeTitle),
		&(c.Transcript),
		&(c.Alt),
		&(c.Img),
		&(c.Title),
		&(c.Year),
		&(c.Month),
		&(c.Day),
		&(c.Num),
	); err != nil {
		err = fmt.Errorf("on row.Scan: %w", err)
		return err
	}

	return nil
}

func (Comic) Save(c *Comic) error {
	conn, err := db.GetInstance()
	if err != nil {
		err = fmt.Errorf("on GetInstance: %v", err)
		return err
	}

	var query string
	var result sql.Result

	if c.ID == 0 {
		query = `INSERT INTO comic
		(Link, News, SafeTitle, Transcript, Alt, Img, Title, Year, Month, Day, Num)
		values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		result, err = conn.Exec(query, c.Link, c.News, c.SafeTitle, c.Transcript, c.Alt, c.Img, c.Title, c.Year, c.Month, c.Day, c.Num)
	} else {
		query = `UPDATE comic
		SET 
		    Link = ?,
			News = ?,
			SafeTitle = ?,
			Transcript = ?,
			Alt = ?,
			Img = ?,
			Title = ?,
			Year = ?,
			Month = ?,
			Day = ?,
			Num = ?
		WHERE id = ?`
		result, err = conn.Exec(query, c.Link, c.News, c.SafeTitle, c.Transcript, c.Alt, c.Img, c.Title, c.Year, c.Month, c.Day, c.Num, c.ID)
	}

	if err != nil {
		err = fmt.Errorf("on exec save query: %w", err)
		return err
	}

	if c.ID == 0 {
		id, err := result.LastInsertId()
		if err != nil {
			err = fmt.Errorf("on get last insert id: %w", err)
			return err
		}

		c.ID = int(id)
	}

	return nil
}

func (Comic) Delete(c Comic) error {
	conn, err := db.GetInstance()
	if err != nil {
		err = fmt.Errorf("on GetInstance: %v", err)
		return err
	}

	_, err = conn.Exec(`DELETE FROM comic WHERE id = ?`, c.ID)
	if err != nil {
		err = fmt.Errorf("on deleting comic: %w", err)
		return err
	}

	return nil
}

func (Comic) ReadAll(limit int, isDesc bool) ([]Comic, error) {
	comics := make([]Comic, 0)
	conn, err := db.GetInstance()
	if err != nil {
		err = fmt.Errorf("on GetInstance: %v", err)
		return nil, err
	}

	query := `SELECT * FROM comic ORDER BY num LIMIT ?`
	if isDesc {
		query = `SELECT * FROM comic ORDER BY num DESC LIMIT ?`
	}

	rows, err := conn.Query(query, limit)

	if err != nil {
		err = fmt.Errorf("on rows query: %w", err)
		return nil, err
	}

	count := 0
	for rows.Next() {
		count++
		comic := new(Comic)
		if err = scanRowToComic(rows, comic); err != nil {
			return nil, err
		}
		comics = append(comics, *comic)
	}

	if count == 0 {
		return nil, sql.ErrNoRows
	}

	return comics, nil
}

func (c Comic) ReadLast() (comic Comic, err error) {
	var comics []Comic
	comics, err = (c).ReadAll(1, true)

	if err != nil {
		err = fmt.Errorf("on ReadAll: %w", err)
		return comic, err
	}

	return comics[0], nil
}
