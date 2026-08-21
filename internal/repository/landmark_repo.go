package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"citywalk/internal/platform/geocoding"
)

type LandmarkFilter struct {
	CategoryID int64
	Limit      int
	Offset     int
}

func (s *Store) CreateLandmark(ctx context.Context, name, address string, lat, lng float64, categoryID *int64, description, coverImage string, createdBy int64, status int8) (int64, error) {
	res, err := s.Exec(ctx, `INSERT INTO landmarks(name, address, lat, lng, category_id, description, cover_image, created_by, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`, name, nullString(address), lat, lng, categoryID, nullString(description), nullString(coverImage), createdBy, status)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) ListLandmarks(ctx context.Context, filter LandmarkFilter) ([]map[string]any, int64, error) {
	clauses := []string{"deleted_at IS NULL"}
	args := make([]any, 0)
	if filter.CategoryID > 0 {
		clauses = append(clauses, "category_id = ?")
		args = append(args, filter.CategoryID)
	}
	where := strings.Join(clauses, " AND ")
	var total int64
	if err := s.QueryRow(ctx, fmt.Sprintf(`SELECT COUNT(1) FROM landmarks WHERE %s`, where), args...).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := s.Query(ctx, fmt.Sprintf(`SELECT id, name, address, lat, lng, category_id, description, cover_image, created_by, status, rating, rating_count, view_count, created_at, updated_at FROM landmarks WHERE %s ORDER BY created_at DESC LIMIT ? OFFSET ?`, where), append(args, filter.Limit, filter.Offset)...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var result []map[string]any
	for rows.Next() {
		var r LandmarkRecord
		if err := rows.Scan(&r.ID, &r.Name, &r.Address, &r.Lat, &r.Lng, &r.CategoryID, &r.Description, &r.CoverImage, &r.CreatedBy, &r.Status, &r.Rating, &r.RatingCount, &r.ViewCount, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, 0, err
		}
		result = append(result, landmarkRecordToMap(r))
	}
	return nonNilSlice(result), total, rows.Err()
}

func (s *Store) FindLandmark(ctx context.Context, id int64) (map[string]any, error) {
	row := s.QueryRow(ctx, `SELECT id, name, address, lat, lng, category_id, description, cover_image, created_by, status, rating, rating_count, view_count, created_at, updated_at FROM landmarks WHERE id = ? AND deleted_at IS NULL`, id)
	var r LandmarkRecord
	if err := row.Scan(&r.ID, &r.Name, &r.Address, &r.Lat, &r.Lng, &r.CategoryID, &r.Description, &r.CoverImage, &r.CreatedBy, &r.Status, &r.Rating, &r.RatingCount, &r.ViewCount, &r.CreatedAt, &r.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return landmarkRecordToMap(r), nil
}

func (s *Store) UpdateLandmark(ctx context.Context, id int64, name, address string, lat, lng float64, categoryID *int64, description, coverImage string, status int8) error {
	_, err := s.Exec(ctx, `UPDATE landmarks SET name = ?, address = ?, lat = ?, lng = ?, category_id = ?, description = ?, cover_image = ?, status = ? WHERE id = ?`, name, nullString(address), lat, lng, categoryID, nullString(description), nullString(coverImage), status, id)
	return err
}

func (s *Store) SoftDeleteLandmark(ctx context.Context, id int64) error {
	_, err := s.Exec(ctx, `UPDATE landmarks SET deleted_at = NOW() WHERE id = ?`, id)
	return err
}

func (s *Store) RateLandmark(ctx context.Context, id int64, score float64) error {
	tx, err := s.Begin(ctx)
	if err != nil {
		return err
	}
	var currentRating float64
	var currentCount int64
	if err := tx.QueryRowContext(ctx, `SELECT rating, rating_count FROM landmarks WHERE id = ? FOR UPDATE`, id).Scan(&currentRating, &currentCount); err != nil {
		_ = tx.Rollback()
		return err
	}
	newCount := currentCount + 1
	newRating := ((currentRating * float64(currentCount)) + score) / float64(newCount)
	if _, err := tx.ExecContext(ctx, `UPDATE landmarks SET rating = ?, rating_count = ? WHERE id = ?`, newRating, newCount, id); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s *Store) TouchLandmarkView(ctx context.Context, id int64) error {
	_, err := s.Exec(ctx, `UPDATE landmarks SET view_count = view_count + 1 WHERE id = ?`, id)
	return err
}

func (s *Store) NearbyLandmarks(ctx context.Context, lat, lng, radius float64, limit int) ([]map[string]any, error) {
	rows, err := s.Query(ctx, `SELECT id, name, address, lat, lng, category_id, description, cover_image, created_by, status, rating, rating_count, view_count, created_at, updated_at FROM landmarks WHERE deleted_at IS NULL AND status = 1 ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []map[string]any
	for rows.Next() {
		var r LandmarkRecord
		if err := rows.Scan(&r.ID, &r.Name, &r.Address, &r.Lat, &r.Lng, &r.CategoryID, &r.Description, &r.CoverImage, &r.CreatedBy, &r.Status, &r.Rating, &r.RatingCount, &r.ViewCount, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}
		dist := geocoding.DistanceMeters(lat, lng, r.Lat, r.Lng)
		if dist <= radius {
			m := landmarkRecordToMap(r)
			m["distance"] = dist
			list = append(list, m)
		}
		if len(list) >= limit {
			break
		}
	}
	return nonNilSlice(list), rows.Err()
}

func landmarkRecordToMap(r LandmarkRecord) map[string]any {
	return map[string]any{"id": r.ID, "name": r.Name, "address": r.Address.String, "lat": r.Lat, "lng": r.Lng, "categoryId": r.CategoryID.Int64, "description": r.Description.String, "coverImage": r.CoverImage.String, "createdBy": r.CreatedBy, "status": r.Status, "rating": r.Rating, "ratingCount": r.RatingCount, "viewCount": r.ViewCount, "createdAt": r.CreatedAt, "updatedAt": r.UpdatedAt}
}
