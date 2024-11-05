package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"sync"

	"github.com/Arinji2/search-backend/types"
	"github.com/google/uuid"
)

func GetKeyword(uid string) (types.SQLKeyword, error) {
	db := getDB()
	query := "SELECT id, keyword, doc_count, idf FROM keywords WHERE id = ?"

	rows, err := db.Query(query, uid)
	if err != nil {
		return types.SQLKeyword{}, err
	}
	defer rows.Close()

	var keyword types.SQLKeyword
	var id []byte
	for rows.Next() {
		err := rows.Scan(&id, &keyword.Keyword, &keyword.DocCount, &keyword.IDF)
		if err != nil {
			return types.SQLKeyword{}, err
		}
		keyword.ID = string(id)
	}

	if err = rows.Err(); err != nil {
		return types.SQLKeyword{}, err
	}

	return keyword, nil
}

func UpdateKeyword(uid string, keyword types.SQLKeyword) error {
	db := getDB()
	query := "UPDATE keywords SET keyword = ?, doc_count = ?, idf = ? WHERE id = ?"

	_, err := db.Exec(query, keyword.Keyword, keyword.DocCount, keyword.IDF, uid)
	if err != nil {
		return err
	}

	return nil
}

func CreateKeyword(keyword types.SQLKeyword) (string, error) {
	db := getDB()
	id := uuid.NewString()
	query := "INSERT INTO keywords (id, keyword, doc_count, idf) VALUES (?, ?, ?, ?)"

	_, err := db.Exec(query, id, keyword.Keyword, keyword.DocCount, keyword.IDF)
	if err != nil {
		return "", err
	}

	return id, nil
}

func KeywordExists(keyword string) (bool, string, error) {
	db := getDB()
	query := "SELECT id FROM keywords WHERE keyword = ?"
	row := db.QueryRow(query, keyword)

	var id string
	err := row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, "", nil
		}
		return false, "", err
	}

	return true, id, nil
}

func GetKeywordsCount() (int, error) {
	db := getDB()
	query := "SELECT COUNT(*) FROM keywords"

	var count int
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func GetAllKeywords() ([]types.SQLKeyword, error) {
	db := getDB()
	query := "SELECT id, keyword, doc_count, idf FROM keywords"

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keywords []types.SQLKeyword
	for rows.Next() {
		var keyword types.SQLKeyword
		var id []byte
		err := rows.Scan(&id, &keyword.Keyword, &keyword.DocCount, &keyword.IDF)
		if err != nil {
			return nil, err
		}
		keyword.ID = string(id)
		keywords = append(keywords, keyword)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return keywords, nil
}
func UpdateIDFScores() error {
	keywords, err := GetAllKeywords()
	if err != nil {
		return err
	}
	pagesCount, err := GetPagesCount()

	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	workersCount := 10
	keywordsChan := make(chan types.SQLKeyword, 10)
	errorChan := make(chan error, len(keywords)*2)

	for i := 0; i < workersCount; i++ {
		go func() {
			for keyword := range keywordsChan {
				existingIDF := keyword.IDF
				keyword.IDF = math.Log(float64(pagesCount) / float64(keyword.DocCount))
				err := UpdateKeyword(keyword.ID, keyword)
				if err != nil {
					errorChan <- err
				}
				fmt.Printf("\n Updated IDF For Keyword %s from %f to %f\n", keyword.Keyword, existingIDF, keyword.IDF)
				wg.Done()
			}
		}()
	}

	for _, keyword := range keywords {
		wg.Add(1)
		keywordsChan <- keyword
	}

	go func() {
		wg.Wait()
		close(errorChan)
	}()

	close(keywordsChan)
	var keywordsError error

	for err := range errorChan {
		keywordsError = errors.Join(keywordsError, err)
	}

	if keywordsError != nil {
		return keywordsError
	}

	return nil
}
