package types

import "time"

type SQLKeyword struct {
	ID       string
	Keyword  string
	DocCount int
	IDF      float64
}

type SQLPage struct {
	ID          string
	URL         string
	Title       string
	MetaImage   string
	Description string
	Favicon     string
	TotalWords  int
}

type SQLPageKeyword struct {
	PageID       string
	KeywordID    string
	WeightedFreq float64
}

type SQLIndexList struct {
	ID      string
	URL     string
	AddedOn time.Time
}
