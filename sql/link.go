package sql

func LinkPageKeyword(keywordID string, pageID string, frequency int) error {
	db := getDB()
	query := "INSERT INTO page_keywords (keyword_id, page_id, frequency) VALUES (?, ?, ?)"
	_, err := db.Exec(query, keywordID, pageID, frequency)
	if err != nil {
		return err
	}
	return nil
}
