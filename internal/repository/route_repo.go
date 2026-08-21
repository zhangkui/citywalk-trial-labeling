package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

type RouteFilter struct {
	City         string
	ThemeID      int64
	MinRate      float64
	CountMinRate float64
	SortBy       string
	Limit        int
	Offset       int
}

type RouteWithWaypoints struct {
	Route     RouteRecord
	Waypoints []map[string]any
}

func (s *Store) CreateRoute(ctx context.Context, title, description string, themeID *int64, city string, startLat, startLng, endLat, endLng *float64, totalDistance, duration *int, difficulty int8, coverImage string, createdBy int64, status int8, waypoints []map[string]any) (int64, error) {
	tx, err := s.Begin(ctx)
	if err != nil {
		return 0, err
	}
	res, err := tx.ExecContext(ctx, `INSERT INTO routes(title, description, theme_id, city, start_lat, start_lng, end_lat, end_lng, total_distance, duration, difficulty, cover_image, created_by, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		title, nullString(description), themeID, nullString(city), startLat, startLng, endLat, endLng, totalDistance, duration, difficulty, nullString(coverImage), createdBy, status)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}
	routeID, err := res.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}
	for idx, wp := range waypoints {
		_, err := tx.ExecContext(ctx, `INSERT INTO waypoints(route_id, name, lat, lng, stay_duration, `+"`order`"+`) VALUES (?, ?, ?, ?, ?, ?)`, routeID, wp["name"], wp["lat"], wp["lng"], wp["stayDuration"], idx)
		if err != nil {
			_ = tx.Rollback()
			return 0, err
		}
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return routeID, nil
}

func (s *Store) ListRoutes(ctx context.Context, filter RouteFilter) ([]map[string]any, int64, error) {
	clauses := []string{"r.deleted_at IS NULL"}
	args := make([]any, 0)
	if filter.City != "" {
		clauses = append(clauses, "r.city = ?")
		args = append(args, filter.City)
	}
	if filter.ThemeID > 0 {
		clauses = append(clauses, "r.theme_id = ?")
		args = append(args, filter.ThemeID)
	}
	countClauses := append([]string(nil), clauses...)
	countArgs := append([]any(nil), args...)
	if filter.CountMinRate > 0 {
		countClauses = append(countClauses, "r.rating >= ?")
		countArgs = append(countArgs, filter.CountMinRate)
	}
	if filter.MinRate > 0 {
		clauses = append(clauses, "r.rating >= ?")
		args = append(args, filter.MinRate)
	}
	where := strings.Join(clauses, " AND ")
	var total int64
	if err := s.QueryRow(ctx, fmt.Sprintf(`SELECT COUNT(1) FROM routes r WHERE %s`, strings.Join(countClauses, " AND ")), countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}
	sortBy := "r.created_at DESC"
	switch filter.SortBy {
	case "hot":
		sortBy = "r.favorite_count DESC, r.view_count DESC"
	case "rating":
		sortBy = "r.rating DESC, r.rating_count DESC"
	case "latest":
		sortBy = "r.created_at DESC"
	}
	query := fmt.Sprintf(`SELECT r.id, r.title, r.description, r.theme_id, r.city, r.start_lat, r.start_lng, r.end_lat, r.end_lng, r.total_distance, r.duration, r.difficulty, r.cover_image, r.created_by, r.status, r.favorite_count, r.rating, r.rating_count, r.view_count, r.created_at, r.updated_at FROM routes r WHERE %s ORDER BY %s LIMIT ? OFFSET ?`, where, sortBy)
	args = append(args, filter.Limit, filter.Offset)
	rows, err := s.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var list []map[string]any
	for rows.Next() {
		var r RouteRecord
		if err := rows.Scan(&r.ID, &r.Title, &r.Description, &r.ThemeID, &r.City, &r.StartLat, &r.StartLng, &r.EndLat, &r.EndLng, &r.TotalDistance, &r.Duration, &r.Difficulty, &r.CoverImage, &r.CreatedBy, &r.Status, &r.FavoriteCount, &r.Rating, &r.RatingCount, &r.ViewCount, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, 0, err
		}
		list = append(list, routeRecordToMap(r))
	}
	return nonNilSlice(list), total, rows.Err()
}

func (s *Store) FindRoute(ctx context.Context, id int64) (map[string]any, error) {
	row := s.QueryRow(ctx, `SELECT id, title, description, theme_id, city, start_lat, start_lng, end_lat, end_lng, total_distance, duration, difficulty, cover_image, created_by, status, favorite_count, rating, rating_count, view_count, created_at, updated_at FROM routes WHERE id = ? AND deleted_at IS NULL`, id)
	var r RouteRecord
	if err := row.Scan(&r.ID, &r.Title, &r.Description, &r.ThemeID, &r.City, &r.StartLat, &r.StartLng, &r.EndLat, &r.EndLng, &r.TotalDistance, &r.Duration, &r.Difficulty, &r.CoverImage, &r.CreatedBy, &r.Status, &r.FavoriteCount, &r.Rating, &r.RatingCount, &r.ViewCount, &r.CreatedAt, &r.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	waypointRows, err := s.Query(ctx, `SELECT id, route_id, name, lat, lng, stay_duration, `+"`order`"+` FROM waypoints WHERE route_id = ? ORDER BY `+"`order`"+` ASC`, id)
	if err != nil {
		return nil, err
	}
	defer waypointRows.Close()
	var waypoints []map[string]any
	for waypointRows.Next() {
		var wp map[string]any
		var idv, routeID int64
		var name string
		var lat, lng float64
		var stay, ord int
		if err := waypointRows.Scan(&idv, &routeID, &name, &lat, &lng, &stay, &ord); err != nil {
			return nil, err
		}
		wp = map[string]any{"id": idv, "routeId": routeID, "name": name, "lat": lat, "lng": lng, "stayDuration": stay, "order": ord}
		waypoints = append(waypoints, wp)
	}
	result := routeRecordToMap(r)
	result["waypoints"] = nonNilSlice(waypoints)
	return result, nil
}

func (s *Store) UpdateRoute(ctx context.Context, id int64, title, description string, themeID *int64, city string, startLat, startLng, endLat, endLng *float64, totalDistance, duration *int, difficulty int8, coverImage string, status int8, waypoints []map[string]any) error {
	if _, err := s.Exec(ctx, `UPDATE routes SET title = ?, description = ?, theme_id = ?, city = ?, start_lat = ?, start_lng = ?, end_lat = ?, end_lng = ?, total_distance = ?, duration = ?, difficulty = ?, cover_image = ?, status = ? WHERE id = ?`,
		title, nullString(description), themeID, nullString(city), startLat, startLng, endLat, endLng, totalDistance, duration, difficulty, nullString(coverImage), status, id); err != nil {
		return err
	}
	if _, err := s.Exec(ctx, `DELETE FROM waypoints WHERE route_id = ?`, id); err != nil {
		return err
	}
	for idx, wp := range waypoints {
		if _, err := s.Exec(ctx, `INSERT INTO waypoints(route_id, name, lat, lng, stay_duration, `+"`order`"+`) VALUES (?, ?, ?, ?, ?, ?)`, id, wp["name"], wp["lat"], wp["lng"], wp["stayDuration"], idx); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) SoftDeleteRoute(ctx context.Context, id int64) error {
	tx, err := s.Begin(ctx)
	if err != nil {
		return err
	}
	res, err := tx.ExecContext(ctx, `UPDATE routes SET deleted_at = NOW(), status = 3 WHERE id = ? AND deleted_at IS NULL`, id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		_ = tx.Rollback()
		return ErrNotFound
	}
	if _, err := tx.ExecContext(ctx, `UPDATE stories SET route_id = NULL WHERE route_id IN (SELECT id FROM routes WHERE id = ? AND deleted_at IS NULL)`, id); err != nil {
		_ = tx.Rollback()
		return err
	}
	if _, err := tx.ExecContext(ctx, `UPDATE events SET route_id = NULL WHERE route_id IN (SELECT id FROM routes WHERE id = ? AND deleted_at IS NULL)`, id); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s *Store) UpdateRouteStatus(ctx context.Context, id int64, status int8) error {
	tx, err := s.Begin(ctx)
	if err != nil {
		return err
	}
	var ownerID int64
	if err := tx.QueryRowContext(ctx, `SELECT created_by FROM routes WHERE id = ? AND deleted_at IS NULL`, id).Scan(&ownerID); err != nil {
		_ = tx.Rollback()
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return err
	}
	rows, err := tx.QueryContext(ctx, `SELECT id FROM routes WHERE created_by = ? AND deleted_at IS NULL ORDER BY id`, ownerID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	var routeIDs []int64
	for rows.Next() {
		var routeID int64
		if err := rows.Scan(&routeID); err != nil {
			_ = rows.Close()
			_ = tx.Rollback()
			return err
		}
		routeIDs = append(routeIDs, routeID)
	}
	if err := rows.Close(); err != nil {
		_ = tx.Rollback()
		return err
	}
	for _, routeID := range routeIDs {
		if _, err := tx.ExecContext(ctx, `UPDATE routes SET status = ? WHERE id = ?`, status, routeID); err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) FavoriteRoute(ctx context.Context, userID, routeID int64, add bool) error {
	tx, err := s.Begin(ctx)
	if err != nil {
		return err
	}
	if add {
		if _, err := tx.ExecContext(ctx, `INSERT IGNORE INTO favorites(user_id, target_type, target_id) VALUES (?, 'route', ?)`, userID, routeID); err != nil {
			_ = tx.Rollback()
			return err
		}
	} else {
		if _, err := tx.ExecContext(ctx, `DELETE FROM favorites WHERE user_id = ? AND target_type = 'route' AND target_id = ?`, userID, routeID); err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	// Recompute favorite_count from the actual favorite rows so the operation is
	// idempotent and self-healing: a repeated add/remove (or a click on an
	// already-correct state) leaves the count unchanged, while a stale count is
	// resynced on the next favorite touch.
	if _, err := tx.ExecContext(ctx, `UPDATE routes SET favorite_count = (SELECT COUNT(*) FROM favorites WHERE target_type = 'route' AND target_id = ?) WHERE id = ?`, routeID, routeID); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s *Store) RateRoute(ctx context.Context, routeID int64, score float64) error {
	tx, err := s.Begin(ctx)
	if err != nil {
		return err
	}
	var currentRating float64
	var currentCount int64
	if err := tx.QueryRowContext(ctx, `SELECT rating, rating_count FROM routes WHERE id = ? FOR UPDATE`, routeID).Scan(&currentRating, &currentCount); err != nil {
		_ = tx.Rollback()
		return err
	}
	newCount := currentCount + 1
	newRating := ((currentRating * float64(currentCount)) + score) / float64(newCount)
	_, err = tx.ExecContext(ctx, `UPDATE routes SET rating = ?, rating_count = ? WHERE id = ?`, newRating, newCount, routeID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func routeRecordToMap(r RouteRecord) map[string]any {
	return map[string]any{
		"id": r.ID, "title": r.Title, "description": r.Description.String, "themeId": r.ThemeID.Int64, "city": r.City.String,
		"startLat": r.StartLat.Float64, "startLng": r.StartLng.Float64, "endLat": r.EndLat.Float64, "endLng": r.EndLng.Float64,
		"totalDistance": r.TotalDistance.Int64, "duration": r.Duration.Int64, "difficulty": r.Difficulty.Int64, "coverImage": r.CoverImage.String,
		"createdBy": r.CreatedBy, "status": r.Status, "favoriteCount": r.FavoriteCount, "rating": r.Rating,
		"ratingCount": r.RatingCount, "viewCount": r.ViewCount, "createdAt": r.CreatedAt, "updatedAt": r.UpdatedAt,
	}
}

func encodeWaypoints(waypoints []map[string]any) string {
	if len(waypoints) == 0 {
		return "[]"
	}
	data, _ := json.Marshal(waypoints)
	return string(data)
}
