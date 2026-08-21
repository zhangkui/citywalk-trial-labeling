package service

import (
	"context"
	"encoding/csv"
	"os"
	"strconv"
)

type ImportService struct { c *Container }

func (s *ImportService) ImportLandmarksCSV(ctx context.Context, path string, createdBy int64) (int, error) {
	file, err := os.Open(path)
	if err != nil { return 0, err }
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil { return 0, err }
	count := 0
	for i, row := range records {
		if i == 0 || len(row) < 5 { continue }
		lat, _ := strconv.ParseFloat(row[2], 64)
		lng, _ := strconv.ParseFloat(row[3], 64)
		var categoryID *int64
		if v, err := strconv.ParseInt(row[4], 10, 64); err == nil { categoryID = &v }
		if _, err := s.c.Deps.Store.CreateLandmark(ctx, row[0], row[1], lat, lng, categoryID, "", "", createdBy, 1); err == nil { count++ }
	}
	return count, nil
}

