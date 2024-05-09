package parser

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"xkcd/internal/models"
)

var ErrNoRecordFound = fmt.Errorf("no record found")

type xkcdResponse struct {
	Link       string `json:"link"`
	News       string `json:"news"`
	SafeTitle  string `json:"safe_title"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
	Title      string `json:"title"`
	Year       int    `json:"year,string"`
	Month      int    `json:"month,string"`
	Day        int    `json:"day,string"`
	Num        int    `json:"num"`
}

func Parse(repo models.ComicRepo) (isEnd bool, err error) {
	comic, err := (models.Comic{}).ReadLast()
	if !errors.Is(err, sql.ErrNoRows) && err != nil {
		err = fmt.Errorf("on get last: %w", err)
		return false, err
	}

	for i := comic.Num + 1; i <= comic.Num+100; i++ {
		res, reqErr := sendRequest(i)
		if reqErr != nil && !errors.Is(reqErr, ErrNoRecordFound) {
			err = fmt.Errorf("on sending request: %w", reqErr)
			return false, err
		} else if errors.Is(reqErr, ErrNoRecordFound) {
			return true, nil
		}

		err = saveComic(repo, res)
		if err != nil {
			err = fmt.Errorf("on saving comic: %w", err)
			return false, err
		}

	}

	return false, nil
}

func saveComic(repo models.ComicRepo, r xkcdResponse) (err error) {
	comic := models.Comic{
		Link:       r.Link,
		News:       r.News,
		SafeTitle:  r.SafeTitle,
		Transcript: r.Transcript,
		Alt:        r.Alt,
		Img:        r.Img,
		Title:      r.Title,
		Year:       r.Year,
		Month:      r.Month,
		Day:        r.Day,
		Num:        r.Num,
	}

	err = repo.Save(&comic)
	if err != nil {
		err = fmt.Errorf("on saving comic to db: %w", err)
		return err
	}

	return nil
}

func sendRequest(page int) (result xkcdResponse, err error) {
	const xkcdUrl = "https://xkcd.com/%d/info.0.json"

	resp, err := http.Get(fmt.Sprintf(xkcdUrl, page))
	defer func() {
		_ = resp.Body.Close()
	}()
	if err != nil {
		err = fmt.Errorf("on send request: %w", err)
		return result, err
	}

	if resp.StatusCode == 404 {
		return result, ErrNoRecordFound
	} else if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("request failed with http status: %s", resp.Status)
		return result, err
	}

	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		err = fmt.Errorf("on parse response: %w", err)
		return result, err
	}

	return result, nil
}
