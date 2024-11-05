package sql

import (
	"database/sql"

	"github.com/Arinji2/search-backend/types"
	"github.com/google/uuid"
)

func GetPage(uid string) (types.SQLPage, error) {
	db := getDB()
	query := "SELECT id, url, title, total_words, meta_image, description FROM pages WHERE id = ?"

	rows, err := db.Query(query, uid)
	if err != nil {
		return types.SQLPage{}, err
	}
	defer rows.Close()

	var page types.SQLPage
	var id []byte
	for rows.Next() {
		err := rows.Scan(&id, &page.URL, &page.Title, &page.MetaImage, &page.Description, &page.TotalWords)
		if err != nil {
			return types.SQLPage{}, err
		}
		page.ID = string(id)
	}

	if err = rows.Err(); err != nil {
		return types.SQLPage{}, err
	}

	return page, nil
}

func UpdatePage(uid string, page types.SQLPage) error {
	db := getDB()
	query := "UPDATE pages SET url = ?, title = ?, total_words = ? WHERE id = ?"

	_, err := db.Exec(query, page.URL, page.Title, page.TotalWords, uid)
	if err != nil {
		return err
	}

	return nil
}

func CreatePage(page types.SQLPage) (string, error) {
	db := getDB()
	id := uuid.NewString()
	query := "INSERT INTO pages (id, url, title, meta_image, description, favicon, total_words) VALUES (?, ?, ?, ?, ?, ?, ?)"

	_, err := db.Exec(query, id, page.URL, page.Title, page.MetaImage, page.Description, page.Favicon, page.TotalWords)
	if err != nil {
		return "", err
	}

	return id, nil
}

func GetPagesCount() (int, error) {
	db := getDB()
	query := "SELECT COUNT(*) FROM pages"

	var count int
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func PageExists(url string) (bool, string, error) {
	db := getDB()
	query := "SELECT id FROM pages WHERE url = ?"
	row := db.QueryRow(query, url)

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
