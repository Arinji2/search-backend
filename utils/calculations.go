package utils

import (
	"math"

	"github.com/Arinji2/search-backend/sql"
)

func CalculateIDF(docCount int) (float64, error) {
	totalPages, err := sql.GetPagesCount()
	if err != nil {
		return 0, err
	}
	if totalPages == 0 {
		return 0, nil
	}
	return math.Log(float64(docCount) / float64(totalPages)), nil
}
