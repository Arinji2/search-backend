package sql

import (
	"database/sql"

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
