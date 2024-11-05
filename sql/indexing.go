package sql

import "database/sql"

func DeletePageIndex(url string) error {
	db := getDB()
	query := "DELETE FROM index_list WHERE url = ?"

	_, err := db.Exec(query, url)
	if err != nil {
		return err
	}

	return nil
}

func CheckIndexURLExists(url string) (bool, error) {
	db := getDB()
	query := "SELECT url FROM index_list WHERE url = ?"
	row := db.QueryRow(query, url)
	var scannedURL string

	err := row.Scan(&scannedURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func AddIndexList(url string) error {
	db := getDB()

	exists, err := CheckIndexURLExists(url)

	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	query := "INSERT INTO index_list (url) VALUES (?)"

	_, err = db.Exec(query, url)
	if err != nil {

		return err
	}

	return nil
}
